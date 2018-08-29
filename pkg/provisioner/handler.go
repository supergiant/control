package provisioner

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/supergiant/pkg/account"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/node"
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/workflows"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

type AccountGetter interface {
	Get(context.Context, string) (*model.CloudAccount, error)
}

type TokenGetter interface {
	GetToken(context.Context, int) (string, error)
}

type Handler struct {
	accountGetter AccountGetter
	tokenGetter   TokenGetter
	provisioner   Provisioner
}

type ProvisionRequest struct {
	ClusterName      string              `json:"clusterName"`
	Profile          profile.KubeProfile `json:"profile"`
	CloudAccountName string              `json:"cloudAccountName"`
}

type ProvisionResponse struct {
	TaskIDs []string `json:"taskIds"`
}

func NewHandler(cloudAccountService *account.Service, tokenGetter TokenGetter, provisioner Provisioner) *Handler {
	return &Handler{
		accountGetter: cloudAccountService,
		tokenGetter:   tokenGetter,
		provisioner:   provisioner,
	}
}

func (h *Handler) Register(m *mux.Router) {
	m.HandleFunc("/provision", h.Provision).Methods(http.MethodPost)
}

func (h *Handler) Provision(w http.ResponseWriter, r *http.Request) {
	req := &ProvisionRequest{}
	err := json.NewDecoder(r.Body).Decode(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logrus.Error(errors.Wrap(err, "unmarshal json"))
		return
	}

	token, err := h.tokenGetter.GetToken(r.Context(), len(req.Profile.MasterProfiles))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Error(errors.Wrap(err, "get token"))
		return
	}

	logrus.Infof("Got token %s", token)

	config := &steps.Config{
		ClusterName: req.ClusterName,
		DigitalOceanConfig: steps.DOConfig{
			Region: req.Profile.Region,
		},
		AWSConfig:    steps.AWSConfig{},
		GCEConfig:    steps.GCEConfig{},
		OSConfig:     steps.OSConfig{},
		PacketConfig: steps.PacketConfig{},

		DockerConfig: steps.DockerConfig{
			Version:        req.Profile.DockerVersion,
			ReleaseVersion: req.Profile.UbuntuVersion,
			Arch:           req.Profile.Arch,
		},
		DownloadK8sBinary: steps.DownloadK8sBinary{
			K8SVersion:      req.Profile.K8SVersion,
			Arch:            req.Profile.Arch,
			OperatingSystem: req.Profile.OperatingSystem,
		},
		CertificatesConfig: steps.CertificatesConfig{
			KubernetesConfigDir: "/etc/kubernetes",
			Username:            "root",
			Password:            "1234",
		},
		NetworkConfig: steps.NetworkConfig{
			EtcdRepositoryUrl: "https://github.com/coreos/etcd/releases/download",
			EtcdVersion:       "3.3.9",
			EtcdHost:          "0.0.0.0",

			Arch:            req.Profile.Arch,
			OperatingSystem: req.Profile.OperatingSystem,

			Network:     req.Profile.CIDR,
			NetworkType: req.Profile.NetworkType,
		},
		FlannelConfig: steps.FlannelConfig{
			Arch:    req.Profile.Arch,
			Version: req.Profile.FlannelVersion,
			// TODO(stgleb): this should be configurable from user side
			EtcdHost: "0.0.0.0",
		},
		KubeletConfig: steps.KubeletConfig{
			MasterPrivateIP: "localhost",
			ProxyPort:       "8080",
			EtcdClientPort:  "2379",
			K8SVersion:      req.Profile.K8SVersion,
		},
		ManifestConfig: steps.ManifestConfig{
			K8SVersion:          req.Profile.K8SVersion,
			KubernetesConfigDir: "/etc/kubernetes",
			RBACEnabled:         req.Profile.RBACEnabled,
			ProviderString:      "todo",
			MasterHost:          "localhost",
			MasterPort:          "8080",
		},
		PostStartConfig: steps.PostStartConfig{
			Host:        "localhost",
			Port:        "8080",
			Username:    "root",
			RBACEnabled: req.Profile.RBACEnabled,
			Timeout:     300,
		},
		TillerConfig: steps.TillerConfig{
			HelmVersion:     req.Profile.HelmVersion,
			OperatingSystem: req.Profile.OperatingSystem,
			Arch:            req.Profile.Arch,
		},
		SshConfig: steps.SshConfig{
			Port:       "22",
			User:       "root",
			PrivateKey: []byte(``),
			Timeout:    10,
		},
		EtcdConfig: steps.EtcdConfig{
			// TODO(stgleb): this field must be changed per node
			Name:           "etcd0",
			Version:        "3.3.9",
			Host:           "0.0.0.0",
			DataDir:        "/etcd-data",
			ServicePort:    "2379",
			ManagementPort: "2380",
			StartTimeout:   "0",
			RestartTimeout: "5",
			DiscoveryUrl:   token,
		},

		MasterNodes:      make(map[string]*node.Node, len(req.Profile.MasterProfiles)),
		Timeout:          time.Second * 1200,
		CloudAccountName: req.CloudAccountName,
	}

	// Fill config with appropriate cloud account credentials
	err = workflows.FillCloudAccountCredentials(r.Context(), h.accountGetter, config)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		logrus.Error(errors.Wrap(err, "fill cloud account"))
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), config.Timeout)
	tasks, err := h.provisioner.Provision(ctx, &req.Profile, config)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Error(errors.Wrap(err, "provision"))
		return
	}

	taskIds := make([]string, 0, len(tasks))

	for _, task := range tasks {
		taskIds = append(taskIds, task.ID)
	}

	resp := ProvisionResponse{
		TaskIDs: taskIds,
	}

	// Respond to client side that request has been accepted
	w.WriteHeader(http.StatusAccepted)

	err = json.NewEncoder(w).Encode(&resp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Error(errors.Wrap(err, "marshal json"))
	}
}
