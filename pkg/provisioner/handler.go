package provisioner

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/supergiant/pkg/account"
	"github.com/supergiant/supergiant/pkg/model"
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
	ClusterName      string          `json:"clusterName"`
	Profile          profile.Profile `json:"profile"`
	CloudAccountName string          `json:"cloudAccountName"`
}

type ProvisionResponse struct {
	Tasks map[string][]string `json:"tasks"`
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

	discoveryUrl, err := h.tokenGetter.GetToken(r.Context(), len(req.Profile.MasterProfiles))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Error(errors.Wrap(err, "get discoveryUrl"))
		return
	}

	logrus.Infof("Got discoveryUrl %s", discoveryUrl)

	config := steps.NewConfig(req.ClusterName, discoveryUrl, req.CloudAccountName, req.Profile)
	// Fill config with appropriate cloud account credentials
	err = workflows.FillCloudAccountCredentials(r.Context(), h.accountGetter, config)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		logrus.Error(errors.Wrap(err, "fill cloud account"))
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), config.Timeout)
	taskMap, err := h.provisioner.Provision(ctx, &req.Profile, config)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Error(errors.Wrap(err, "provision"))
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
