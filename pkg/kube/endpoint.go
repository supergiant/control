package kube

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/asaskevich/govalidator.v8"
)

type KubeEndpoint struct {
	service *KubeService
}

func (e *KubeEndpoint) GetKube(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	kubeId := vars["id"]

	kube, err := e.service.Get(r.Context(), kubeId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(kube)
}

func (e *KubeEndpoint) GetAllKubes(w http.ResponseWriter, r *http.Request) {
	kubes, err := e.service.GetAll(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(kubes)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (e *KubeEndpoint) CreateKube(w http.ResponseWriter, r *http.Request) {
	kube := &Kube{}
	err := json.NewDecoder(r.Body).Decode(&kube)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ok, err := govalidator.ValidateStruct(kube)

	if !ok {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = e.service.Create(r.Context(), kube)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (e *KubeEndpoint) DeleteKube(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	kubeId := vars["id"]

	err := e.service.Delete(r.Context(), kubeId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
