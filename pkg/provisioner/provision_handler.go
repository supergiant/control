package provisioner

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/supergiant/supergiant/pkg/account"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/sgerrors"
)

type ProfileGetter interface {
	Get(context.Context, string) (*profile.KubeProfile, error)
}

type AccountGetter interface {
	Get(context.Context, string) (*model.CloudAccount, error)
}

type ProvisionHandler struct {
	profileGetter ProfileGetter
	accountGetter AccountGetter
	provisioner   Provisioner
}

type ProvisionRequest struct {
	ProfileId        string `json:"profileId"`
	CloudAccountName string `json:"cloudAccountName"`
}

type ProvisionResponse struct {
	TaskIds []string `json:"taskIds"`
}

func NewHandler(kubeService *profile.KubeProfileService, cloudAccountService *account.Service, provisioner Provisioner) *ProvisionHandler {
	return &ProvisionHandler{
		profileGetter: kubeService,
		accountGetter: cloudAccountService,
		provisioner:   provisioner,
	}
}

func (h *ProvisionHandler) Provision(w http.ResponseWriter, r *http.Request) {
	req := &ProvisionRequest{}
	err := json.NewDecoder(r.Body).Decode(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	kubeProfile, err := h.profileGetter.Get(r.Context(), req.ProfileId)

	if err != nil {
		if sgerrors.IsNotFound(err) {
			http.NotFound(w, r)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tasks, err := h.provisioner.Provision(context.Background(), kubeProfile)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
