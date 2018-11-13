package provisioner

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"gopkg.in/asaskevich/govalidator.v8"

	"github.com/supergiant/supergiant/pkg/account"
	"github.com/supergiant/supergiant/pkg/message"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

type AccountGetter interface {
	Get(context.Context, string) (*model.CloudAccount, error)
}

type KubeGetter interface {
	Get(ctx context.Context, name string) (*model.Kube, error)
}

type TokenGetter interface {
	GetToken(context.Context, int) (string, error)
}

type Handler struct {
	accountGetter AccountGetter
	tokenGetter   TokenGetter
	kubeGetter    KubeGetter
	provisioner   ClusterProvisioner
}

type ProvisionRequest struct {
	ClusterName      string          `json:"clusterName" valid:"matches(^[A-Za-z0-9-]+$)"`
	Profile          profile.Profile `json:"profile" valid:"-"`
	CloudAccountName string          `json:"cloudAccountName" valid:"-"`
}

type ProvisionResponse struct {
	Tasks map[string][]string `json:"tasks"`
}

type ClusterProvisioner interface {
	ProvisionCluster(context.Context, *profile.Profile, *steps.Config) (map[string][]*workflows.Task, error)
}

func NewHandler(kubeService KubeGetter, cloudAccountService *account.Service, tokenGetter TokenGetter, provisioner ClusterProvisioner) *Handler {
	return &Handler{
		kubeGetter:    kubeService,
		accountGetter: cloudAccountService,
		tokenGetter:   tokenGetter,
		provisioner:   provisioner,
	}
}

func (h *Handler) Register(m *mux.Router) {
	m.HandleFunc("/provision", h.Provision).Methods(http.MethodPost)
}

// TODO(stgleb): Move this to KubeHandler create kube
func (h *Handler) Provision(w http.ResponseWriter, r *http.Request) {
	req := &ProvisionRequest{}
	err := json.NewDecoder(r.Body).Decode(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logrus.Error(errors.Wrap(err, "unmarshal json"))
		return
	}

	ok, err := govalidator.ValidateStruct(req)
	if !ok {
		logrus.Errorf("Validation error %v", err.Error())
		message.SendValidationFailed(w, err)
		return
	}

	discoveryUrl, err := h.tokenGetter.GetToken(r.Context(),
		len(req.Profile.MasterProfiles))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Error(errors.Wrap(err, "get discoveryUrl"))
		return
	}

	logrus.Infof("Got discoveryUrl %s", discoveryUrl)
	config := steps.NewConfig(req.ClusterName, discoveryUrl,
		req.CloudAccountName, req.Profile)

	acc, err := h.accountGetter.Get(r.Context(), req.CloudAccountName)

	if err != nil {
		if sgerrors.IsNotFound(err) {
			http.NotFound(w, r)
			return
		}

		message.SendUnknownError(w, err)
		return
	}

	// Fill config with appropriate cloud account credentials
	err = util.FillCloudAccountCredentials(r.Context(), acc, config)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		logrus.Error(errors.Wrap(err, "fill cloud account"))
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), config.Timeout)
	taskMap, err := h.provisioner.ProvisionCluster(ctx, &req.Profile, config)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Error(errors.Wrap(err, "provisionCluster"))
		return
	}

	roleTaskIdMap := make(map[string][]string, len(taskMap))

	for role, taskSet := range taskMap {
		roleTaskIdMap[role] = make([]string, 0, len(taskSet))

		for _, task := range taskSet {
			roleTaskIdMap[role] = append(roleTaskIdMap[role], task.ID)
		}
	}

	resp := ProvisionResponse{
		Tasks: roleTaskIdMap,
	}

	// Respond to client side that request has been accepted
	w.WriteHeader(http.StatusAccepted)

	err = json.NewEncoder(w).Encode(&resp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Error(errors.Wrap(err, "marshal json"))
	}
}
