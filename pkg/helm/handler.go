package helm

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gopkg.in/asaskevich/govalidator.v8"

	"github.com/supergiant/supergiant/pkg/message"
	"github.com/supergiant/supergiant/pkg/model/helm"
	"github.com/supergiant/supergiant/pkg/sgerrors"
)

// Handler is a http controller for a helm repositories.
type Handler struct {
	svc *Service
}

// NewHandler constructs a Handler for helm repositories.
func NewHandler(svc *Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

// Register adds helm specific api to the main handler.
func (h *Handler) Register(r *mux.Router) {
	r.HandleFunc("/helm/repositories", h.createRepo).Methods(http.MethodPost)
	r.HandleFunc("/helm/repositories/{repoName}", h.getRepo).Methods(http.MethodGet)
	r.HandleFunc("/helm/repositories", h.listAllRepos).Methods(http.MethodGet)
	r.HandleFunc("/helm/repositories/{repoName}", h.deleteRepo).Methods(http.MethodDelete)

	//r.HandleFunc("/helm/releases", h.createRelease).Methods(http.MethodPost)
	//r.HandleFunc("/helm/releases/{relName}", h.getRelease).Methods(http.MethodGet)
	//r.HandleFunc("/helm/releases", h.listAllReleases).Methods(http.MethodGet)
	//r.HandleFunc("/helm/releases/{relName}", h.deleteRelease).Methods(http.MethodDelete)
}

func (h *Handler) createRepo(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) getRepo(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) listAllRepos(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) deleteRepo(w http.ResponseWriter, r *http.Request) {
	repoName := mux.Vars(r)["repoName"]

	if err := h.svc.Delete(r.Context(), repoName); err != nil {
		logrus.Error(err)
		message.SendUnknownError(w, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
