package profile

import (
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
)

const prefix  = "/profile/"

type KubeProfileHandler struct{
	service KubeProfileService
}

func (h *KubeProfileHandler) GetProfile(w http.ResponseWriter,r *http.Request) {
	vars := mux.Vars(r)
	profileId := vars["id"]

	kubeProfile, err := h.service.Get(profileId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(kubeProfile)
}

func (h *KubeProfileHandler) CreateProfile(w http.ResponseWriter,r *http.Request) {
	profile := &KubeProfile{}

	err := json.NewDecoder(r.Body).Decode(&profile)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.Create(profile)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *KubeProfileHandler) GetProfiles(w http.ResponseWriter,r *http.Request) {
	profiles, err := h.service.GetAll()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(profiles)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
