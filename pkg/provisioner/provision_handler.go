package provisioner

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/supergiant/supergiant/pkg/account"
	"github.com/supergiant/supergiant/pkg/profile"
)

type ProvisionHandler struct {
	kubeService         profile.KubeProfileService
	cloudAccountService account.Service
	provisioner         Provisioner
}

type ProvisionRequest struct {
	ProfileId        string `json:"profileId"`
	CloudAccountName string `json:"cloudAccountName"`
}

type ProvisionResponse struct {
	TaskIds []string `json:"taskIds"`
}

func NewHandler(kubeService profile.KubeProfileService, cloudAccountService account.Service, provisioner Provisioner) *ProvisionHandler {
	return &ProvisionHandler{
		kubeService:         kubeService,
		cloudAccountService: cloudAccountService,
		provisioner:         provisioner,
	}
}

func (h *ProvisionHandler) Provision(w http.ResponseWriter, r *http.Request) {
	req := &ProvisionRequest{}
	err := json.NewDecoder(r.Body).Decode(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	kubeProfile, err := h.kubeService.Get(r.Context(), req.ProfileId)

	if err != nil {
		http.NotFound(w, r)
	}

	tasks, err := h.provisioner.Provision(context.Background(), kubeProfile)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	taskIds := make([]string, len(tasks))

	for _, task := range tasks {
		taskIds = append(taskIds, task.ID)
	}

	resp := ProvisionResponse{
		TaskIds: taskIds,
	}

	// Respond to client side that request has been accepted
	w.WriteHeader(http.StatusAccepted)

	err = json.NewEncoder(w).Encode(&resp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
