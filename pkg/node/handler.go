package node

import "net/http"

type Handler struct {
	service Service
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) Put(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {

}
