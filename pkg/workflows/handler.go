package workflows

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/message"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/runner"
	"github.com/supergiant/control/pkg/runner/ssh"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/storage"
	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows/steps"
)

type cloudAccountGetter interface {
	Get(context.Context, string) (*model.CloudAccount, error)
}

type TaskHandler struct {
	runnerFactory  func(config ssh.Config) (runner.Runner, error)
	cloudAccGetter cloudAccountGetter
	repository     storage.Interface
	getWriter      func(string) (io.WriteCloser, error)
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
		getWriter: func(name string) (io.WriteCloser, error) {
			// TODO(stgleb): Add log directory to params of supergiant
			return os.OpenFile(path.Join("/tmp", name), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		},
	}
}

func (h *TaskHandler) Register(m *mux.Router) {
	m.HandleFunc("/tasks", h.RunTask).Methods(http.MethodPost)
	m.HandleFunc("/tasks/{id}", h.GetTask).Methods(http.MethodGet)
	m.HandleFunc("/tasks/{id}/restart",
		h.RestartTask).Methods(http.MethodPost)
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

	data, err := h.repository.Get(r.Context(), Prefix, id)

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

	acc, err := h.cloudAccGetter.Get(r.Context(), req.Cfg.CloudAccountName)

	if err != nil {
		if sgerrors.IsNotFound(err) {
			http.NotFound(w, r)
			return
		}

		message.SendUnknownError(w, err)
		return
	}

	// Get cloud account fill appropriate config structure with cloud account credentials
	err = util.FillCloudAccountCredentials(r.Context(), acc, &req.Cfg)

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

	logrus.Debugf("get task %s", id)
	data, err := h.repository.Get(r.Context(), Prefix, id)

	if err != nil {
		logrus.Debugf("task %s not found", id)
		http.NotFound(w, r)
		return
	}

	task, err := DeserializeTask(data, h.repository)

	if err != nil {
		logrus.Debugf("error deserializing task %s %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fileName := util.MakeFileName(id)
	writer, err := h.getWriter(fileName)

	if err != nil {
		http.Error(w, fmt.Sprintf("get writer %v", err), http.StatusInternalServerError)
		logrus.Errorf("Get writer %v", err)
		return
	}

	task.Restart(context.Background(), id, writer)
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

// NOTE(stgleb): This is made for testing purposes and example, remove when UI is done.
func (h *TaskHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	tokenString := ""
	ts := strings.Split(authHeader, " ")

	if len(ts) > 1 {
		tokenString = ts[1]
	}

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
                var conn = new WebSocket("ws://{{ .Host }}/v1/api/tasks/{{ .TaskID }}/logs?token={{ .Token }}");
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
		TaskID string
		Token  string
	}{
		r.Host,
		id,
		tokenString,
	}
	err := tpl.Execute(w, &v)
	if err != nil {
		logrus.Error(err)
	}
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
		// TODO(stgleb): Do something more safe in future
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	t, err := tail.TailFile(path.Join("/tmp", util.MakeFileName(id)),
		tail.Config{
			Follow:    true,
			MustExist: true,
			Location: &tail.SeekInfo{
				Offset: 0,
				Whence: io.SeekStart,
			},
			Logger:      tail.DiscardingLogger,
			MaxLineSize: 160,
		})

	if os.IsNotExist(err) {
		http.NotFound(w, r)
		logrus.Errorf("Not found %s", util.MakeFileName(id))
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Errorf("Open file %s for tail %v", util.MakeFileName(id), err)
		return
	}

	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Errorf("Upgrade connection %v", err)
		return
	}

	go func() {
		pingTicker := time.NewTicker(time.Second * 60)

		for {
			select {
			case line := <-t.Lines:
				c.SetWriteDeadline(time.Now().Add(time.Second * 10))
				err = c.WriteMessage(websocket.TextMessage, []byte(line.Text))

				// Do not log this error, since client can simply disconnect
				if err != nil {
					return
				}
			case <-pingTicker.C:
				c.SetWriteDeadline(time.Now().Add(time.Second * 10))
				if err := c.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					return
				}
			}
		}
	}()
}
