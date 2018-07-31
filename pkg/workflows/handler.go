package workflows

import (
	"net/http"

	"github.com/gorilla/mux"

	"context"
	"encoding/json"
	"os"

	"github.com/supergiant/supergiant/pkg/runner/ssh"
	"github.com/supergiant/supergiant/pkg/storage"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

type WorkflowHandler struct {
	repository storage.Interface
}

type BuildWorkFlowRequest struct {
	stepNames []string     `json:"step_names"`
	config    steps.Config `json:"config"`
	sshConfig ssh.Config   `json:"ssh_config"`
}

func (h *WorkflowHandler) GetWorkflow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		http.Error(w, "need id of workflow", http.StatusBadRequest)
		return
	}

	data, err := h.repository.Get(r.Context(), "workflows", id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func (h *WorkflowHandler) BuildWorkflow(w http.ResponseWriter, r *http.Request) {
	var req BuildWorkFlowRequest

	err := json.NewDecoder(r.Body).Decode(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	runner, err := ssh.NewRunner(req.sshConfig)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.config.Runner = runner
	s := make([]steps.Step, 0, len(req.stepNames))

	for _, stepName := range req.stepNames {
		s = append(s, steps.GetStep(stepName))
	}

	workflow := BuildCustomWorkflow(s, req.config, h.repository)
	workflow.Run(context.Background(), os.Stdout)
}
