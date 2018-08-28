package workflows

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/supergiant/pkg/runner"
	"github.com/supergiant/supergiant/pkg/runner/ssh"
	"github.com/supergiant/supergiant/pkg/storage"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

type TaskHandler struct {
	runnerFactory  func(config ssh.Config) (runner.Runner, error)
	cloudAccGetter cloudAccountGetter
	repository     storage.Interface
}

type RunTaskRequest struct {
	WorkflowName string       `json:"workflowName"`
	Cfg          steps.Config `json:"config"`
}

type BuildTaskRequest struct {
	StepNames []string     `json:"stepNames"`
	Cfg       steps.Config `json:"config"`
	SshConfig ssh.Config   `json:"sshConfig"`
}

type TaskResponse struct {
	ID string `json:"id"`
}

func NewTaskHandler(repository storage.Interface, runnerFactory func(config ssh.Config) (runner.Runner, error), getter cloudAccountGetter) *TaskHandler {
	return &TaskHandler{
		runnerFactory:  runnerFactory,
		repository:     repository,
		cloudAccGetter: getter,
	}
}

func (h *TaskHandler) Register(m *mux.Router) {
	m.HandleFunc("/tasks", h.RunTask).Methods(http.MethodPost)
	m.HandleFunc("/tasks/{id}", h.GetTask).Methods(http.MethodGet)
	m.HandleFunc("/tasks/{id}/restart", h.RestartTask).Methods(http.MethodPost)
	m.HandleFunc("/tasks/{id}/logs", h.StreamLogs).Methods(http.MethodGet)
	m.HandleFunc("/tasks/{id}/logs/ws", h.GetLogs).Methods(http.MethodGet)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		http.Error(w, "need id of task", http.StatusBadRequest)
		return
	}

	data, err := h.repository.Get(r.Context(), prefix, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func (h *TaskHandler) RunTask(w http.ResponseWriter, r *http.Request) {
	req := &RunTaskRequest{}
	err := json.NewDecoder(r.Body).Decode(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get cloud account fill appropriate config structure with cloud account credentials
	err = FillCloudAccountCredentials(r.Context(), h.cloudAccGetter, &req.Cfg)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	task, err := NewTask(req.WorkflowName, h.repository)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	task.Run(context.Background(), req.Cfg, os.Stdout)

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(&TaskResponse{
		task.ID,
	})
}

func (h *TaskHandler) RestartTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		http.Error(w, "need id of task", http.StatusBadRequest)
		return
	}

	data, err := h.repository.Get(r.Context(), prefix, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	task, err := deserializeTask(data, h.repository)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	task.Restart(context.Background(), id, os.Stdout)
	w.WriteHeader(http.StatusAccepted)
}

func (h *TaskHandler) BuildAndRunTask(w http.ResponseWriter, r *http.Request) {
	req := &BuildTaskRequest{}
	err := json.NewDecoder(r.Body).Decode(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create newTask sshRunner with config provided
	sshRunner, err := h.runnerFactory(req.SshConfig)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Cfg.Runner = sshRunner
	s := make([]steps.Step, 0, len(req.StepNames))
	// Get steps for task
	for _, stepName := range req.StepNames {
		s = append(s, steps.GetStep(stepName))
	}

	// TODO(stgleb): pass here workflow type DOMaster or DONode
	task := newTask("", s, h.repository)
	// We ignore cancel function since we cannot get it back
	ctx, _ := context.WithTimeout(context.Background(), req.Cfg.Timeout*time.Second)
	task.Run(ctx, req.Cfg, os.Stdout)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&TaskResponse{
		task.ID,
	})
}

func (h *TaskHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	page := `<!DOCTYPE html>
<html lang="en">
    <head>
        <title>WebSocket Example</title>
    </head>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/3.3.1/jquery.js"></script>
    
    <body>
        <div id="logs"></div>
        <script type="text/javascript">
            (function() {
                var conn = new WebSocket("ws://{{ .Host }}/tasks/{{ .TaskId }}/logs");
                conn.onmessage = function(evt) {
                    console.log('file updated');
 					$('#logs').append("<p>" + evt.data + "</p>");
                }
            })();
        </script>
    </body>
</html>`
	tpl := template.Must(template.New("").Parse(page))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		http.Error(w, "need id of task", http.StatusBadRequest)
		return
	}

	var v = struct {
		Host   string
		TaskId string
	}{
		r.Host,
		id,
	}
	tpl.Execute(w, &v)
}

func (h *TaskHandler) StreamLogs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		http.Error(w, "need id of task", http.StatusBadRequest)
		return
	}

	var upgrader = websocket.Upgrader{
		HandshakeTimeout: time.Second * 10,
		WriteBufferSize:  1024,
		ReadBufferSize:   0,
	}

	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fd, err := os.OpenFile(path.Join("/tmp", fmt.Sprintf("%s.log", id)), os.O_RDONLY, 0666)

	if os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := ioutil.ReadAll(fd)
	c.SetWriteDeadline(time.Now().Add(time.Second * 10))

	err = c.WriteMessage(websocket.TextMessage, data)

	if err != nil {
		logrus.Errorf("Error while write to websocket %v", err)
		return
	}

	go func(){
		pingTicker := time.NewTicker(time.Second * 60)

		t, err := tail.TailFile(path.Join("/tmp", fmt.Sprintf("%s.log", id)), tail.Config{Follow: true, MaxLineSize: 160})

		for {
			select {
			case line := <-t.Lines:
				c.SetWriteDeadline(time.Now().Add(time.Second * 10))
				err = c.WriteMessage(websocket.TextMessage, []byte(line.Text))

				if err != nil {
					logrus.Errorf("Error while write to websocket %v", err)
					return
				}
			case <-pingTicker.C:
				err = c.SetWriteDeadline(time.Now().Add(time.Second * 10))

				if err := c.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					return
				}
			}
		}
	}()
}
