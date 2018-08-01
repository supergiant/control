package workflows

import (
	"net/http"

	"github.com/gorilla/mux"

	"context"
	"encoding/json"
	"os"

	"github.com/supergiant/supergiant/pkg/runner"
	"github.com/supergiant/supergiant/pkg/runner/ssh"
	"github.com/supergiant/supergiant/pkg/storage"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

type TaskHandler struct {
	runnerFactory func(config ssh.Config) (runner.Runner, error)
	repository    storage.Interface
}

type BuildTaskRequest struct {
	StepNames []string     `json:"step_names"`
	Cfg       steps.Config `json:"Cfg"`
	SshConfig ssh.Config   `json:"ssh_config"`
}

type BuildTaskResponse struct {
	Id string `json:"id"`
}

func NewTaskHandler(repository storage.Interface, runnerFactory func(config ssh.Config) (runner.Runner, error)) *TaskHandler {
	return &TaskHandler{
		runnerFactory: runnerFactory,
		repository:    repository,
	}
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

func (h *TaskHandler) BuildAndRunTask(w http.ResponseWriter, r *http.Request) {
	req := &BuildTaskRequest{}
	err := json.NewDecoder(r.Body).Decode(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create new runner with config provided
	runner, err := h.runnerFactory(req.SshConfig)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Cfg.Runner = runner
	s := make([]steps.Step, 0, len(req.StepNames))
	// Get steps for task
	for _, stepName := range req.StepNames {
		s = append(s, steps.GetStep(stepName))
	}

	task := New(s, req.Cfg, h.repository)
	// TODO(stgleb): We should provide custom timeout for task execution
	task.Run(context.Background(), os.Stdout)

	respData, _ := json.Marshal(BuildTaskResponse{task.Id})

	w.WriteHeader(http.StatusCreated)
	w.Write(respData)
}
