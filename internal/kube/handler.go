package kube

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc}
}

func (h *Handler) Register(r mux.Router) {
	r.HandleFunc("/kubes", h.CreateKube).Methods(http.MethodPost)
	r.HandleFunc("/kubes", h.ListKubes).Methods(http.MethodGet)
	r.HandleFunc("/kubes/{id}", h.GetKube).Methods(http.MethodGet)
	r.HandleFunc("/kubes/{id}", h.DeleteKube).Methods(http.MethodDelete)
	r.HandleFunc("/kubes/{id}/list", h.ListResources).Methods(http.MethodGet)
	r.HandleFunc("/kubes/{id}/resources/{resource}", h.GetResource).Methods(http.MethodGet)
}

func (h *Handler) CreateKube(w http.ResponseWriter, r *http.Request) {
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
	w.Write([]byte("kube has been added: " + k.ID))
}

func (h *Handler) GetKube(w http.ResponseWriter, r *http.Request) {
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

	json.NewEncoder(w).Encode(k)
}

func (h *Handler) DeleteKube(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) ListKubes(w http.ResponseWriter, r *http.Request) {
	kubes, err := h.svc.ListAll(r.Context())
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(kubes)
}

func (h *Handler) ListResources(w http.ResponseWriter, r *http.Request) {
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

	w.Write(rawResources)
}

func (h *Handler) GetResource(w http.ResponseWriter, r *http.Request) {
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

	w.Write(rawResources)
}

func notImplemented(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "not implemented!\n")
}
