package node

import (
	"encoding/json"
	"net/http"

	"gopkg.in/asaskevich/govalidator.v8"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/supergiant/pkg/message"
	"github.com/supergiant/supergiant/pkg/sgerrors"
)

type Handler struct {
	service Service
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profileId := vars["id"]

	node, err := h.service.Get(r.Context(), profileId)
	if err != nil {
		if sgerrors.IsNotFound(err) {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(node); err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) Put(w http.ResponseWriter, r *http.Request) {
	k := &Node{}
	err := json.NewDecoder(r.Body).Decode(k)
	if err != nil {
		message.SendInvalidJSON(w, err)
		return
	}

	ok, err := govalidator.ValidateStruct(k)
	if !ok {
		message.SendValidationFailed(w, err)
		return
	}

	if err = h.service.Create(r.Context(), k); err != nil {
		message.SendUnknownError(w, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	kname := vars["id"]
	if err := h.service.Delete(r.Context(), kname); err != nil {
		if sgerrors.IsNotFound(err) {
			message.SendNotFound(w, kname, err)
			return
		}
		message.SendUnknownError(w, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	kubes, err := h.service.ListAll(r.Context())
	if err != nil {
		message.SendUnknownError(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(kubes); err != nil {
		message.SendUnknownError(w, err)
	}
}
