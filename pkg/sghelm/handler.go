package sghelm

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"k8s.io/helm/pkg/repo"

	"github.com/supergiant/control/pkg/message"
	"github.com/supergiant/control/pkg/sgerrors"
)

// Handler is a http controller for a helm repositories.
type Handler struct {
	svc Servicer
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
	r.HandleFunc("/helm/repositories", h.listRepos).Methods(http.MethodGet)
	r.HandleFunc("/helm/repositories/{repoName}", h.deleteRepo).Methods(http.MethodDelete)

	r.HandleFunc("/helm/repositories/{repoName}/charts", h.listCharts).Methods(http.MethodGet)
	r.HandleFunc("/helm/repositories/{repoName}/charts/{chartName}", h.getChartData).Methods(http.MethodGet)
}

func (h *Handler) createRepo(w http.ResponseWriter, r *http.Request) {
	repoConf := &repo.Entry{}
	if err := json.NewDecoder(r.Body).Decode(repoConf); err != nil {
		log.Errorf("helm: create repository: decode: %s", err)
		message.SendValidationFailed(w, err)
		return
	}

	// TODO: use a custom struct instead of repo.Entry
	repoConf.Name, repoConf.URL = strings.TrimSpace(repoConf.Name), strings.TrimSpace(repoConf.URL)
	if repoConf.Name == "" || repoConf.URL == "" {
		log.Errorf("helm: create repository: validation failed: %+v", repoConf)
		message.SendValidationFailed(w, errors.New("helm repository: name and url should be provided"))
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
		log.Errorf("helm: create repository: %s: write resp: %s", repoConf.Name, err)
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

func (h *Handler) listRepos(w http.ResponseWriter, r *http.Request) {
	repos, err := h.svc.ListRepos(r.Context())
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

	hrepo, err := h.svc.DeleteRepo(r.Context(), repoName)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, repoName, err)
			return
		}
		log.Errorf("helm: delete repository: %s: %s", repoName, err)
		message.SendUnknownError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(hrepo); err != nil {
		log.Errorf("helm: delete repository: %s: encode: %s", hrepo.Config.Name, err)
		message.SendUnknownError(w, err)
		return
	}
}

func (h *Handler) getChartData(w http.ResponseWriter, r *http.Request) {
	repoName := mux.Vars(r)["repoName"]
	chartName := mux.Vars(r)["chartName"]

	version := r.URL.Query().Get("version")
	chrt, err := h.svc.GetChartData(r.Context(), repoName, chartName, version)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, repoName+"/"+chartName, err)
			return
		}
		log.Errorf("helm: get %s/%s chart: %s", repoName, chartName, err)
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
