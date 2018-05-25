package account

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/gorilla/mux"
	"gopkg.in/asaskevich/govalidator.v8"
)

type Endpoint struct {
	Service *Service
}

func (e *Endpoint) Create(rw http.ResponseWriter, r *http.Request) {
	account := new(CloudAccount)
	err := json.NewDecoder(r.Body).Decode(account)
	//TODO Revise error handling
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

func (e *Endpoint) ListAll(rw http.ResponseWriter, r *http.Request) {
	accounts, err := e.Service.GetAll(r.Context())
	if err != nil {
		logrus.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(rw).Encode(accounts)
}

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

func (e *Endpoint) Update(rw http.ResponseWriter, r *http.Request) {
	acc := new(CloudAccount)
	err := json.NewDecoder(r.Body).Decode(acc)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	err = e.Service.Update(r.Context(), acc)
	if err != nil {
		logrus.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

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
