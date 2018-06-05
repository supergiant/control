package repository

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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
	err := json.NewDecoder(r.Body).Decode(repo)
	if err != nil {
		http.Error(w, "can't unmarshal request body", http.StatusBadRequest)
		return
	}

	ok, err := govalidator.ValidateStruct(repo)
	if !ok {
		http.Error(w, "invalid repository", http.StatusBadRequest)
		return
	}

	if err = h.svc.Create(r.Context(), repo); err != nil {
		logrus.Errorf("handler: create %s helm repo: %v", repo.Name, err)
		http.Error(w, "failed to create", http.StatusInternalServerError)
		return
	}
}

// Get retrieves a helm repository.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	repoName := mux.Vars(r)["repoName"]

	println("reponame:", repoName)

	repo, err := h.svc.Get(r.Context(), repoName)
	if err != nil {
		logrus.Errorf("handler: get %s helm repository: %v", repoName, err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if repo == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(repo)
}

// ListAll retrieves all helm repositories.
func (h *Handler) ListAll(w http.ResponseWriter, r *http.Request) {
	repos, err := h.svc.GetAll(r.Context())
	if err != nil {
		logrus.Errorf("handler: get all helm repositories: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(repos)
}

// Delete removes a helm repository from the storage.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	repoName := mux.Vars(r)["repoName"]

	err := h.svc.Delete(r.Context(), repoName)
	if err != nil {
		logrus.Errorf("handler: delete %s helm repository: %v", repoName, err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
