package profile

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"gopkg.in/asaskevich/govalidator.v8"

	"github.com/supergiant/supergiant/pkg/storage"
)

const prefix = "/profile/"

type KubeProfileEndpoint struct {
	service *KubeProfileService
}

func NewKubeProfileEndpoint(prefix string, storage storage.Interface) *KubeProfileEndpoint {
	return &KubeProfileEndpoint{
		service: NewKubeProfileService(prefix, storage),
	}
}

func (h *KubeProfileEndpoint) GetProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profileId := vars["id"]

	kubeProfile, err := h.service.Get(r.Context(), profileId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(kubeProfile)
}

func (h *KubeProfileEndpoint) CreateProfile(w http.ResponseWriter, r *http.Request) {
	profile := &KubeProfile{}

	err := json.NewDecoder(r.Body).Decode(&profile)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ok, err := govalidator.ValidateStruct(profile)

	if !ok {
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

func (h *KubeProfileEndpoint) GetProfiles(w http.ResponseWriter, r *http.Request) {
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
