package profile

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gopkg.in/asaskevich/govalidator.v8"

	"github.com/supergiant/supergiant/pkg/sgerrors"
)

type NodeProfileHandler struct {
	service *NodeProfileService
}

func NewNodeProfileHandler(svc *NodeProfileService) *NodeProfileHandler {
	return &NodeProfileHandler{
		service: svc,
	}
}

func (h *NodeProfileHandler) Register(r *mux.Router) {
	r.HandleFunc("/nodeprofiles/{id}", h.GetProfile).Methods(http.MethodGet)
	r.HandleFunc("/nodeprofiles", h.CreateProfile).Methods(http.MethodPost)
	r.HandleFunc("/nodeprofiles", h.GetProfiles).Methods(http.MethodGet)
}

func (h *NodeProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profileId := vars["id"]

	nodeProfile, err := h.service.Get(r.Context(), profileId)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(nodeProfile)
}

func (h *NodeProfileHandler) CreateProfile(w http.ResponseWriter, r *http.Request) {
	profile := &NodeProfile{}

	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		logrus.Error(err)
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

func (h *NodeProfileHandler) GetProfiles(w http.ResponseWriter, r *http.Request) {
	profiles, err := h.service.GetAll(r.Context())
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(profiles); err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
