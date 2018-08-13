package workflows

import (
	"context"
	"encoding/json"
	"net/http"

	"os"

	"time"

	"github.com/gorilla/mux"

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

type RestartTaskRequest struct {
	Id string `json:"id"`
}

type TaskResponse struct {
	Id string `json:"id"`
}

func NewTaskHandler(repository storage.Interface, runnerFactory func(config ssh.Config) (runner.Runner, error), getter cloudAccountGetter) *TaskHandler {
	return &TaskHandler{
		runnerFactory:  runnerFactory,
		repository:     repository,
		cloudAccGetter: getter,
	}
}

func (h *TaskHandler) Register(m *mux.Router) {
	m.HandleFunc("/tasks", h.GetTask).Methods(http.MethodGet)
	m.HandleFunc("/tasks", h.RunTask).Methods(http.MethodPost)
	m.HandleFunc("/tasks/restart", h.RestartTask).Methods(http.MethodPost)
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
	//fillCloudAccountCredentials(r.Context(), h.cloudAccGetter, &req.Cfg)

	task, err := NewTask(req.WorkflowName, req.Cfg, h.repository)
	task.Run(context.Background(), os.Stdout)

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(&TaskResponse{
		task.Id,
	})
}

func (h *TaskHandler) RestartTask(w http.ResponseWriter, r *http.Request) {
	req := &RestartTaskRequest{}
	err := json.NewDecoder(r.Body).Decode(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := h.repository.Get(r.Context(), prefix, req.Id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	task, err := deserializeTask(data, h.repository)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	task.Restart(context.Background(), req.Id, os.Stdout)
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
	task := newTask("", s, req.Cfg, h.repository)
	// We ignore cancel function since we cannot get it back
	ctx, _ := context.WithTimeout(context.Background(), req.Cfg.Timeout*time.Second)
	task.Run(ctx, os.Stdout)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&TaskResponse{
		task.Id,
	})
}
