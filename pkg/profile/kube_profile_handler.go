package profile

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gopkg.in/asaskevich/govalidator.v8"

	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/storage"
)

type KubeProfileHandler struct {
	service *KubeProfileService
}

func NewKubeProfileEndpoint(prefix string, storage storage.Interface) *KubeProfileHandler {
	return &KubeProfileHandler{
		service: NewKubeProfileService(prefix, storage),
	}
}

func (h *KubeProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profileId := vars["id"]

	kubeProfile, err := h.service.Get(r.Context(), profileId)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(kubeProfile); err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *KubeProfileHandler) CreateProfile(w http.ResponseWriter, r *http.Request) {
	profile := &KubeProfile{}

	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ok, err := govalidator.ValidateStruct(profile)
	if !ok {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.Create(r.Context(), profile); err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *KubeProfileHandler) GetProfiles(w http.ResponseWriter, r *http.Request) {
	profiles, err := h.service.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(profiles); err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
