package account

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/asaskevich/govalidator.v8"

	"github.com/supergiant/control/pkg/message"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows/steps"
)

// Handler is a http controller for account entity
type Handler struct {
	validator util.CloudAccountValidator
	service   *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		validator: util.NewCloudAccountValidator(),
		service:   service,
	}
}

func (h *Handler) Register(r *mux.Router) {
	r.HandleFunc("/accounts", h.Create).Methods(http.MethodPost)
	r.HandleFunc("/accounts", h.ListAll).Methods(http.MethodGet)
	r.HandleFunc("/accounts/{accountName}", h.Get).Methods(http.MethodGet)
	r.HandleFunc("/accounts/{accountName}", h.Update).Methods(http.MethodPut)
	r.HandleFunc("/accounts/{accountName}", h.Delete).Methods(http.MethodDelete)
	r.HandleFunc("/accounts/{accountName}/regions", h.GetRegions).Methods(http.MethodGet)
	r.HandleFunc("/accounts/{accountName}/regions/{region}/az", h.GetAZs).Methods(http.MethodGet)
	r.HandleFunc("/accounts/{accountName}/regions/{region}/az/{az}/types", h.GetTypes).Methods(http.MethodGet)
}

// Create register new cloud account
func (h *Handler) Create(rw http.ResponseWriter, r *http.Request) {
	account := new(model.CloudAccount)
	if err := json.NewDecoder(r.Body).Decode(account); err != nil {
		message.SendInvalidJSON(rw, err)
		return
	}

	ok, err := govalidator.ValidateStruct(account)
	if !ok {
		logrus.Errorf("Error validating account struct %v", err)
		message.SendValidationFailed(rw, err)
		return
	}

	// Check account data for validity
	if err := h.validator.ValidateCredentials(account); err != nil {
		logrus.Errorf("error validating credentials %v", err)
		message.SendInvalidCredentials(rw, err)
		return
	}

	if err = h.service.Create(r.Context(), account); err != nil {
		if sgerrors.IsUnsupportedProvider(err) {
			message.SendMessage(rw, message.New(fmt.Sprintf("Unsupported provider %s", account.Provider),
				err.Error(), sgerrors.UnsupportedProvider, ""), http.StatusBadRequest)
			return
		}

		if sgerrors.IsAlreadyExists(err) {
			message.SendAlreadyExists(rw, account.Name, sgerrors.ErrAlreadyExists)
			return
		}

		logrus.Errorf("account handler: create %v", err)
		message.SendUnknownError(rw, err)
		return
	}
}

// ListAll retrieves all cloud accounts
func (h *Handler) ListAll(rw http.ResponseWriter, r *http.Request) {
	accounts, err := h.service.GetAll(r.Context())
	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(rw, "accounts", err)
			return
		}

		logrus.Errorf("account handler: list all %v", err)
		message.SendUnknownError(rw, err)
		return
	}
	if err := json.NewEncoder(rw).Encode(accounts); err != nil {
		logrus.Errorf("account handler: list all %v", err)
		message.SendUnknownError(rw, err)
		return
	}
}

// Get retrieves individual account by name
func (h *Handler) Get(rw http.ResponseWriter, r *http.Request) {
	accountName := mux.Vars(r)["accountName"]
	account, err := h.service.Get(r.Context(), accountName)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(rw, "account", err)
			return
		}
		logrus.Errorf("account handler: get %v", err)
		message.SendUnknownError(rw, err)
		return
	}

	if err := json.NewEncoder(rw).Encode(account); err != nil {
		logrus.Errorf("account handler: get %v", err)
		message.SendUnknownError(rw, err)
		return
	}
}

// TODO(stgleb): Use patch for updating
// Update saves updated state of an cloud account, account name can't be changed
func (h *Handler) Update(rw http.ResponseWriter, r *http.Request) {
	account := new(model.CloudAccount)
	if err := json.NewDecoder(r.Body).Decode(account); err != nil {
		message.SendInvalidJSON(rw, err)
		return
	}

	ok, err := govalidator.ValidateStruct(account)
	if !ok {
		message.SendValidationFailed(rw, err)
		return
	}
	if err := h.service.Update(r.Context(), account); err != nil {
		logrus.Errorf("account handler: update: %v", err)
		message.SendUnknownError(rw, err)
		return
	}
}

// Delete cloud account
func (h *Handler) Delete(rw http.ResponseWriter, r *http.Request) {
	accountName := mux.Vars(r)["accountName"]
	if accountName == "" {
		msg := message.New("account name can't be blank", "", sgerrors.InvalidJSON, "")
		message.SendMessage(rw, msg, http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), accountName); err != nil {
		logrus.Errorf("account handler: delete %v", err)
		message.SendUnknownError(rw, err)
		return
	}
}

func (h *Handler) GetRegions(w http.ResponseWriter, r *http.Request) {
	accountName, ok := mux.Vars(r)["accountName"]
	if !ok || accountName == "" {
		message.SendValidationFailed(w, errors.New("clouds: preconditions failed"))
		return
	}

	acc, err := h.service.Get(r.Context(), accountName)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, "account", err)
			return
		}
		logrus.Errorf("clouds: get regions %v", err)
		message.SendUnknownError(w, err)
		return
	}

	config := &steps.Config{}
	getter, err := NewRegionsGetter(acc, config)
	if err != nil {
		logrus.Errorf("clouds: get regions %v", err)
		message.SendUnknownError(w, err)
		return
	}

	aggregate, err := getter.GetRegions(r.Context())
	if err != nil {
		logrus.Errorf("clouds: get regions %v", err)
		message.SendUnknownError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(aggregate); err != nil {
		logrus.Errorf("clouds: get regions %v", err)
		message.SendUnknownError(w, err)
		return
	}
}

func (h *Handler) GetAZs(w http.ResponseWriter, r *http.Request) {
	accountName, ok := mux.Vars(r)["accountName"]
	if !ok || accountName == "" {
		message.SendValidationFailed(w, errors.New("clouds: " +
			"preconditions failed"))
		return
	}

	region, ok := mux.Vars(r)["region"]
	if region == "" {
		message.SendValidationFailed(w, errors.New("clouds: " +
			"preconditions failed"))
		return
	}

	acc, err := h.service.Get(r.Context(), accountName)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, "account", err)
			return
		}

		logrus.Errorf("clouds: get account %s error: %v",
			accountName, err)
		message.SendUnknownError(w, err)
		return
	}

	acc.Credentials["region"] = region
	config := &steps.Config{}
	getter, err := NewZonesGetter(acc, config)
	if err != nil {
		logrus.Errorf("clouds: get %s availability zones %v",
			acc.Provider, err)
		message.SendUnknownError(w, err)
		return
	}

	azs, err := getter.GetZones(r.Context(), *config)
	if err != nil {
		logrus.Errorf("clouds: get %s availability zones %v",
			acc.Provider, err)
		message.SendUnknownError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(azs); err != nil {
		logrus.Errorf("clouds: get %s availability zones %v",
			acc.Provider, err)
		message.SendUnknownError(w, err)
		return
	}
}

func (h *Handler) GetTypes(w http.ResponseWriter, r *http.Request) {
	accountName, ok := mux.Vars(r)["accountName"]
	if !ok || accountName == "" {
		message.SendValidationFailed(w, errors.New("clouds: " +
			"preconditions failed"))
		return
	}

	region, ok := mux.Vars(r)["region"]
	if region == "" {
		message.SendValidationFailed(w, errors.New("clouds: " +
			"preconditions failed"))
		return
	}

	az, ok := mux.Vars(r)["az"]
	if az == "" {
		message.SendValidationFailed(w, errors.New("clouds: " +
			"preconditions failed"))
		return
	}

	acc, err := h.service.Get(r.Context(), accountName)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, "account", err)
			return
		}

		logrus.Errorf("clouds: get types %s %v", accountName, err)
		message.SendUnknownError(w, err)
		return
	}

	acc.Credentials["availabilityZone"] = az
	acc.Credentials["region"] = region

	config := &steps.Config{}
	getter, err := NewTypesGetter(acc, config)
	if err != nil {
		logrus.Errorf("clouds: get %s types %v", acc.Provider, err)
		message.SendUnknownError(w, err)
		return
	}

	types, err := getter.GetTypes(r.Context(), *config)
	if err != nil {
		logrus.Errorf("clouds: get %s types %v", acc.Provider, err)
		message.SendUnknownError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(types); err != nil {
		logrus.Errorf("clouds: get %s aws types %v", acc.Provider, err)
		message.SendUnknownError(w, err)
		return
	}
}
