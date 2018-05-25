package profile

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type NodeProfileEndpoint struct {
	service *NodeProfileService
}

func (h *NodeProfileEndpoint) GetProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profileId := vars["id"]

	nodeProfile, err := h.service.Get(r.Context(), profileId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(nodeProfile)
}

func (h *NodeProfileEndpoint) CreateProfile(w http.ResponseWriter, r *http.Request) {
	profile := &NodeProfile{}

	err := json.NewDecoder(r.Body).Decode(&profile)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.Create(r.Context(), profile)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *NodeProfileEndpoint) GetProfiles(w http.ResponseWriter, r *http.Request) {
	profiles, err := h.service.GetAll(r.Context())

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
