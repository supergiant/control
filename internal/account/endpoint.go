package account

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gopkg.in/asaskevich/govalidator.v8"
)

// Endpoint is a http controller for account entity
type Endpoint struct {
	Service *Service
}

// Create register new cloud account
func (e *Endpoint) Create(rw http.ResponseWriter, r *http.Request) {
	account := new(CloudAccount)
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	ok, err := govalidator.ValidateStruct(account)
	if !ok {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if err = e.Service.Create(r.Context(), account); err != nil {
		logrus.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ListAll retrieves all cloud accounts
func (e *Endpoint) ListAll(rw http.ResponseWriter, r *http.Request) {
	accounts, err := e.Service.GetAll(r.Context())
	if err != nil {
		logrus.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(rw).Encode(accounts)
}

// Get retrieves individual account by name
func (e *Endpoint) Get(rw http.ResponseWriter, r *http.Request) {
	accountName := mux.Vars(r)["accountName"]
	if accountName == "" {
		http.Error(rw, "account name can't be blank", http.StatusBadRequest)
		return
	}
	account, err := e.Service.Get(r.Context(), accountName)
	if err != nil {
		logrus.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if account == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(rw).Encode(account)
}

// Update saves updated state of an cloud account, account name can't be changed
func (e *Endpoint) Update(rw http.ResponseWriter, r *http.Request) {
	account := new(CloudAccount)
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	ok, err := govalidator.ValidateStruct(account)
	if !ok {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	err = e.Service.Update(r.Context(), account)
	if err != nil {
		logrus.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Delete cloud account
func (e *Endpoint) Delete(rw http.ResponseWriter, r *http.Request) {
	accountName := mux.Vars(r)["accountName"]

	if accountName == "" {
		http.Error(rw, "account name can't be blank", http.StatusBadRequest)
		return
	}
	err := e.Service.Delete(r.Context(), accountName)
	if err != nil {
		logrus.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
