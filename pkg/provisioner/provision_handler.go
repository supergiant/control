package provisioner

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/supergiant/supergiant/pkg/account"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/workflows"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

type ProfileGetter interface {
	Get(context.Context, string) (*profile.KubeProfile, error)
}

type AccountGetter interface {
	Get(context.Context, string) (*model.CloudAccount, error)
}

type TokenGetter interface {
	GetToken(context.Context, int) (string, error)
}

type ProvisionHandler struct {
	profileGetter ProfileGetter
	accountGetter AccountGetter
	tokenGetter   TokenGetter
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

	token, err := h.tokenGetter.GetToken(r.Context(), len(kubeProfile.MasterProfiles))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	config := &steps.Config{
		DigitalOceanConfig: steps.DOConfig{},
		AWSConfig:          steps.AWSConfig{},
		GCEConfig:          steps.GCEConfig{},
		OSConfig:           steps.OSConfig{},
		PacketConfig:       steps.PacketConfig{},

		DockerConfig:                steps.DockerConfig{},
		DownloadK8sBinary:           steps.DownloadK8sBinary{},
		CertificatesConfig:          steps.CertificatesConfig{},
		FlannelConfig:               steps.FlannelConfig{},
		KubeletConfig:               steps.KubeletConfig{},
		ManifestConfig:              steps.ManifestConfig{},
		PostStartConfig:             steps.PostStartConfig{},
		KubeletSystemdServiceConfig: steps.KubeletSystemdServiceConfig{},
		TillerConfig:                steps.TillerConfig{},
		SshConfig:                   steps.SshConfig{},
		EtcdConfig: steps.EtcdConfig{
			Token: token,
		},
	}

	// Fill config with appropriate cloud account credentials
	err = workflows.FillCloudAccountCredentials(r.Context(), h.accountGetter, config)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	tasks, err := h.provisioner.Provision(context.Background(), kubeProfile, config)

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
