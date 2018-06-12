package repository

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/supergiant/pkg/message"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"gopkg.in/asaskevich/govalidator.v8"

	"github.com/supergiant/supergiant/pkg/model/helm"
	"github.com/supergiant/supergiant/pkg/storage"
)

// Handler is a http controller for a helm repositories.
type Handler struct {
	svc *Service
}

// New constructs a Handler for helm repositories.
func NewHandler(s storage.Interface) *Handler {
	return &Handler{
		svc: NewService(s),
	}
}

// Create stores a new helm repository.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	repo := new(helm.Repository)
	if err := json.NewDecoder(r.Body).Decode(repo); err != nil {
		message.SendInvalidJSON(w, err)
		return
	}

	ok, err := govalidator.ValidateStruct(repo)
	if !ok {
		message.SendValidationFailed(w, err)
		return
	}

	if err = h.svc.Create(r.Context(), repo); err != nil {
		message.SendUnknownError(w, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// Get retrieves a helm repository.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	repoName := mux.Vars(r)["repoName"]

	repo, err := h.svc.Get(r.Context(), repoName)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, "helm repository", err)
			return
		}
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(repo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ListAll retrieves all helm repositories.
func (h *Handler) ListAll(w http.ResponseWriter, r *http.Request) {
	repos, err := h.svc.GetAll(r.Context())
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(repos); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Delete removes a helm repository from the storage.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	repoName := mux.Vars(r)["repoName"]

	if err := h.svc.Delete(r.Context(), repoName); err != nil {
		logrus.Error(err)
		message.SendUnknownError(w, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
