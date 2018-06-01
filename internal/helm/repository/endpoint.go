package repository

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gopkg.in/asaskevich/govalidator.v8"

	"github.com/supergiant/supergiant/pkg/model/helm"
)

// Controller is a http controller for a helm repositories.
type Controller struct {
	svc *Service
}

// New constructs a Controller for helm repositories.
func New(svc *Service) (*Controller, error) {
	return &Controller{
		svc: svc,
	}, nil
}

// Create stores a new helm repository.
func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	repo := new(helm.Repository)
	err := json.NewDecoder(r.Body).Decode(repo)
	if err != nil {
		http.Error(w, "can't unmarshal request body", http.StatusBadRequest)
		return
	}

	ok, err := govalidator.ValidateStruct(repo)
	if !ok {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = c.svc.Create(r.Context(), repo); err != nil {
		logrus.Errorf("controller: create %s helm repo: %v", repo.Name, err)
		http.Error(w, "failed to create", http.StatusInternalServerError)
		return
	}
}

// Get retrieves a helm repository.
func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	repoName := mux.Vars(r)["repoName"]
	if strings.TrimSpace(repoName) == "" {
		http.Error(w, "name can't be empty", http.StatusBadRequest)
		return
	}

	repo, err := c.svc.Get(r.Context(), repoName)
	if err != nil {
		logrus.Errorf("controller: get %s relm repository: %v", repo.Name, err)
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
func (c *Controller) ListAll(w http.ResponseWriter, r *http.Request) {
	repos, err := c.svc.GetAll(r.Context())
	if err != nil {
		logrus.Errorf("controller: get all relm repositories: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(repos)
}

// Delete removes a helm repository from the storage.
func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	repoName := mux.Vars(r)["repoName"]
	if strings.TrimSpace(repoName) == "" {
		http.Error(w, "name can't be empty", http.StatusBadRequest)
		return
	}

	err := c.svc.Delete(r.Context(), repoName)
	if err != nil {
		logrus.Errorf("controller: delete %s relm repository: %v", repoName, err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
