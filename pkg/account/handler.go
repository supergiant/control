package account

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/supergiant/pkg/message"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"gopkg.in/asaskevich/govalidator.v8"
)

// Handler is a http controller for account entity
type Handler struct {
	service *Service
}

// Create register new cloud account
func (e *Handler) Create(rw http.ResponseWriter, r *http.Request) {
	account := new(CloudAccount)
	if err := json.NewDecoder(r.Body).Decode(account); err != nil {
		message.SendInvalidJSON(rw, err)
		return
	}

	ok, err := govalidator.ValidateStruct(account)
	if !ok {
		message.SendValidationFailed(rw, err)
		return
	}

	if err = e.service.Create(r.Context(), account); err != nil {
		logrus.Errorf("create account: %v", err)
		message.SendUnknownError(rw, err)
		return
	}
}

// ListAll retrieves all cloud accounts
func (e *Handler) ListAll(rw http.ResponseWriter, r *http.Request) {
	accounts, err := e.service.GetAll(r.Context())
	if err != nil {
		logrus.Errorf("accounts list all %v", err)
		message.SendUnknownError(rw, err)
		return
	}
	if err := json.NewEncoder(rw).Encode(accounts); err != nil {
		logrus.Errorf("accounts list all %v", err)
		message.SendUnknownError(rw, err)
		return
	}
}

// Get retrieves individual account by name
func (e *Handler) Get(rw http.ResponseWriter, r *http.Request) {
	accountName := mux.Vars(r)["accountName"]
	if accountName == "" {
		msg := message.New("account name can't be blank", "", sgerrors.CantChangeID, "")
		message.SendMessage(rw, msg, http.StatusBadRequest)
		return
	}
	account, err := e.service.Get(r.Context(), accountName)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(rw, "account", err)
			return
		}
		logrus.Errorf("account get %v", err)
		message.SendUnknownError(rw, err)
		return
	}

	if err := json.NewEncoder(rw).Encode(account); err != nil {
		logrus.Errorf("account get %v", err)
		message.SendUnknownError(rw, err)
		return
	}
}

// Update saves updated state of an cloud account, account name can't be changed
func (e *Handler) Update(rw http.ResponseWriter, r *http.Request) {
	account := new(CloudAccount)
	if err := json.NewDecoder(r.Body).Decode(account); err != nil {
		message.SendInvalidJSON(rw, err)
		return
	}

	ok, err := govalidator.ValidateStruct(account)
	if !ok {
		message.SendValidationFailed(rw, err)
		return
	}
	if err := e.service.Update(r.Context(), account); err != nil {
		logrus.Errorf("account update: %v", err)
		message.SendUnknownError(rw, err)
		return
	}
}

// Delete cloud account
func (e *Handler) Delete(rw http.ResponseWriter, r *http.Request) {
	accountName := mux.Vars(r)["accountName"]
	if accountName == "" {
		msg := message.New("account name can't be blank", "", sgerrors.CantChangeID, "")
		message.SendMessage(rw, msg, http.StatusBadRequest)
		return
	}

	if err := e.service.Delete(r.Context(), accountName); err != nil {
		logrus.Errorf("account delete %v", err)
		message.SendUnknownError(rw, err)
		return
	}
}
