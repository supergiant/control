package kube

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Handler is a http controller for a kube entity.
type Handler struct {
	svc *Service
}

// NewHandler constructs a Handler for kubes.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc}
}

// Register adds kube handlers to a router.
func (h *Handler) Register(r mux.Router) {
	r.HandleFunc("/kubes", h.createKube).Methods(http.MethodPost)
	r.HandleFunc("/kubes", h.listKubes).Methods(http.MethodGet)
	r.HandleFunc("/kubes/{id}", h.getKube).Methods(http.MethodGet)
	r.HandleFunc("/kubes/{id}", h.deleteKube).Methods(http.MethodDelete)
	r.HandleFunc("/kubes/{id}/list", h.listResources).Methods(http.MethodGet)
	r.HandleFunc("/kubes/{id}/resources/{resource}", h.getResource).Methods(http.MethodGet)
}

func (h *Handler) createKube(w http.ResponseWriter, r *http.Request) {
	k := &Kube{}
	err := json.NewDecoder(r.Body).Decode(k)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if _, err = h.svc.CreateKube(r.Context(), k); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	_, err = w.Write([]byte("kube has been added: " + k.Name))
	handle(err)
}

func (h *Handler) getKube(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars == nil {
		http.Error(w, "invalid url path", http.StatusBadRequest)
	}

	id := vars["id"]
	k, err := h.svc.Get(r.Context(), id)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	handle(json.NewEncoder(w).Encode(k))
}

func (h *Handler) deleteKube(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars == nil {
		http.Error(w, "invalid url path", http.StatusBadRequest)
	}

	id := vars["id"]
	if err := h.svc.Delete(r.Context(), id); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) listKubes(w http.ResponseWriter, r *http.Request) {
	kubes, err := h.svc.ListAll(r.Context())
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	handle(json.NewEncoder(w).Encode(kubes))
}

func (h *Handler) listResources(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars == nil {
		http.Error(w, "invalid url path", http.StatusBadRequest)
		return
	}

	id := vars["id"]

	rawResources, err := h.svc.ListKubeResources(r.Context(), id)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(rawResources)
	handle(err)
}

func (h *Handler) getResource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars == nil {
		http.Error(w, "invalid url path", http.StatusBadRequest)
	}

	id := vars["id"]
	rs := vars["resource"]
	ns := r.URL.Query().Get("namespace")
	name := r.URL.Query().Get("name")

	rawResources, err := h.svc.GetKubeResources(r.Context(), id, rs, ns, name)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(rawResources)
	handle(err)
}

func handle(err error) {
	if err != nil {
		logrus.Errorf("kube handler: http write: %v", err)
	}
}
