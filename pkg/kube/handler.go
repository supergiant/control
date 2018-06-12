package kube

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/supergiant/pkg/message"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"gopkg.in/asaskevich/govalidator.v8"
)

// Handler exposes operations over kubernetes cluster via http
type Handler struct {
	service *Service
}

func (e *Handler) GetKube(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	kubeId := vars["id"]
	kube, err := e.service.Get(r.Context(), kubeId)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, "kubernetes cluster", err)
			return
		}
		logrus.Errorf("error while GetKube %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(kube)
}

func (e *Handler) GetAllKubes(w http.ResponseWriter, r *http.Request) {
	kubes, err := e.service.GetAll(r.Context())
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(kubes); err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (e *Handler) CreateKube(w http.ResponseWriter, r *http.Request) {
	kube := &Kube{}
	err := json.NewDecoder(r.Body).Decode(&kube)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ok, err := govalidator.ValidateStruct(kube)
	if !ok {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = e.service.Create(r.Context(), kube); err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (e *Handler) DeleteKube(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	kubeId := vars["id"]

	if err := e.service.Delete(r.Context(), kubeId); err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusAccepted)
}
