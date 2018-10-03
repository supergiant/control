package helm

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"k8s.io/helm/pkg/repo"

	"github.com/supergiant/supergiant/pkg/message"
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

	r.HandleFunc("/helm/repositories/{repoName}/charts", h.listCharts).Methods(http.MethodGet)
	r.HandleFunc("/helm/repositories/{repoName}/charts/{chartName}", h.getChart).Methods(http.MethodGet)
}

func (h *Handler) createRepo(w http.ResponseWriter, r *http.Request) {
	repoConf := &repo.Entry{}
	if err := json.NewDecoder(r.Body).Decode(repoConf); err != nil {
		log.Errorf("helm: create repository: decode: %s", err)
		message.SendInvalidJSON(w, err)
		return
	}

	hrepo, err := h.svc.CreateRepo(r.Context(), repoConf)
	if err != nil {
		if sgerrors.IsAlreadyExists(err) {
			message.SendAlreadyExists(w, repoConf.Name, err)
			return
		}
		log.Errorf("helm: create repository: %s: %s", repoConf.Name, err)
		message.SendUnknownError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(hrepo); err != nil {
		log.Errorf("helm: create repository: %s: encode: %s", repoConf.Name, err)
		message.SendUnknownError(w, err)
		return
	}
}

func (h *Handler) getRepo(w http.ResponseWriter, r *http.Request) {
	repoName := mux.Vars(r)["repoName"]

	hrepo, err := h.svc.GetRepo(r.Context(), repoName)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, repoName, err)
			return
		}
		log.Errorf("helm: get repository: %s: %s", repoName, err)
		message.SendUnknownError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(hrepo); err != nil {
		log.Errorf("helm: get repository: %s: encode: %s", repoName, err)
		message.SendUnknownError(w, err)
		return
	}
}

func (h *Handler) listAllRepos(w http.ResponseWriter, r *http.Request) {
	repos, err := h.svc.GetAllRepos(r.Context())
	if err != nil {
		log.Errorf("helm: list repositories: %s", err)
		message.SendUnknownError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(repos); err != nil {
		log.Errorf("helm: list repositories: encode: %s", err)
		message.SendUnknownError(w, err)
		return
	}
}

func (h *Handler) deleteRepo(w http.ResponseWriter, r *http.Request) {
	repoName := mux.Vars(r)["repoName"]

	if err := h.svc.DeleteRepo(r.Context(), repoName); err != nil {
		if sgerrors.IsAlreadyExists(err) {
			message.SendAlreadyExists(w, repoName, err)
			return
		}
		log.Errorf("helm: delete repository: %s: %s", repoName, err)
		message.SendUnknownError(w, err)
		return
	}
}

func (h *Handler) getChart(w http.ResponseWriter, r *http.Request) {
	repoName := mux.Vars(r)["repoName"]
	chartName := mux.Vars(r)["chartName"]

	chrt, err := h.svc.GetChart(r.Context(), repoName, chartName)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, repoName+"/"+chartName, err)
			return
		}
		log.Errorf("helm: get chart: %s/%s: %s", repoName, chartName, err)
		message.SendUnknownError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(chrt); err != nil {
		log.Errorf("helm: get chart: %s/%s: encode: %s", repoName, chartName, err)
		message.SendUnknownError(w, err)
		return
	}
}

func (h *Handler) listCharts(w http.ResponseWriter, r *http.Request) {
	repoName := mux.Vars(r)["repoName"]

	chrtList, err := h.svc.ListCharts(r.Context(), repoName)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, repoName, err)
			return
		}
		log.Errorf("helm: list charts: %s: %s", repoName, err)
		message.SendUnknownError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(chrtList); err != nil {
		log.Errorf("helm: list chart: %s: encode: %s", repoName, err)
		message.SendUnknownError(w, err)
		return
	}
}
