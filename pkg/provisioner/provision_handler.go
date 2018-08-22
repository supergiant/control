package provisioner

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
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

func NewHandler(kubeService *profile.KubeProfileService, cloudAccountService *account.Service,
	tokenGetter TokenGetter, provisioner Provisioner) *ProvisionHandler {
	return &ProvisionHandler{
		profileGetter: kubeService,
		accountGetter: cloudAccountService,
		tokenGetter:   tokenGetter,
		provisioner:   provisioner,
	}
}

func (h *ProvisionHandler) Register(m *mux.Router) {
	m.HandleFunc("/provision", h.Provision).Methods(http.MethodPost)
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
		DigitalOceanConfig: steps.DOConfig{
			Region: "fra1",
			Size:   "s-1vcpu-2gb",
			Image:  "ubuntu-18-04-x64",
		},
		AWSConfig:    steps.AWSConfig{},
		GCEConfig:    steps.GCEConfig{},
		OSConfig:     steps.OSConfig{},
		PacketConfig: steps.PacketConfig{},

		DockerConfig: steps.DockerConfig{
			Version:        kubeProfile.DockerVersion,
			ReleaseVersion: kubeProfile.UbuntuVersion,
			Arch:           kubeProfile.Arch,
		},
		DownloadK8sBinary: steps.DownloadK8sBinary{
			K8SVersion:      kubeProfile.K8SVersion,
			Arch:            kubeProfile.Arch,
			OperatingSystem: kubeProfile.OperatingSystem,
		},
		CertificatesConfig: steps.CertificatesConfig{
			KubernetesConfigDir: "/etc/kubernetes",
			Username:            "root",
			Password:            "1234",
		},
		FlannelConfig: steps.FlannelConfig{
			Arch:    kubeProfile.Arch,
			Version: kubeProfile.FlannelVersion,
			// TODO(stgleb): this should be configurable from user side
			Network:     "10.0.0.0",
			NetworkType: kubeProfile.NetworkType,
		},
		KubeletConfig: steps.KubeletConfig{
			MasterPrivateIP: "localhost",
			ProxyPort:       "8080",
			EtcdClientPort:  "2379",
			K8SVersion:      kubeProfile.K8SVersion,
		},
		ManifestConfig: steps.ManifestConfig{
			K8SVersion:          kubeProfile.K8SVersion,
			KubernetesConfigDir: "/etc/kubernetes",
			RBACEnabled:         kubeProfile.RBACEnabled,
			ProviderString:      "todo",
			MasterHost:          "localhost",
			MasterPort:          "8080",
		},
		PostStartConfig: steps.PostStartConfig{
			Host:        "localhost",
			Port:        "8080",
			Username:    "root",
			RBACEnabled: false,
		},
		TillerConfig: steps.TillerConfig{
			HelmVersion:     kubeProfile.HelmVersion,
			OperatingSystem: kubeProfile.OperatingSystem,
			Arch:            kubeProfile.Arch,
		},
		SshConfig: steps.SshConfig{
			Port: "22",
			User: "root",
			PrivateKey: []byte(``),
			Timeout: 10,
		},
		EtcdConfig: steps.EtcdConfig{
			Token: token,
		},

		CloudAccountName: req.CloudAccountName,
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
