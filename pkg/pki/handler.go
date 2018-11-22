package pki

import (
	"encoding/json"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/message"
	"github.com/supergiant/control/pkg/sgerrors"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r mux.Router) {
	// TODO: move it to the kubes handler
	r.HandleFunc("/certificates", h.GetAll).Methods(http.MethodGet)
	r.HandleFunc("/certificates/{id}", h.Get).Methods(http.MethodGet)
	r.HandleFunc("/certificates/{id}", h.Delete).Methods(http.MethodDelete)
	r.HandleFunc("/certificates", h.Generate).Methods(http.MethodPost)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	pkis, err := h.svc.GetAll(r.Context())
	if err != nil {
		logrus.Errorf("pki: %v", err)
		message.SendUnknownError(w, err)
	}
	if err := json.NewEncoder(w).Encode(pkis); err != nil {
		message.SendUnknownError(w, err)
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	m := mux.Vars(r)
	if ID, ok := m["id"]; ok {

		pki, err := h.svc.Get(r.Context(), ID)

		if err != nil {
			if sgerrors.IsNotFound(err) {
				message.SendNotFound(w, "certificates", err)
				return
			}

			logrus.Errorf("pki: %v", err)
			message.SendUnknownError(w, err)
		}

		if err := json.NewEncoder(w).Encode(pki); err != nil {
			message.SendUnknownError(w, err)
		}
	} else {
		message.SendNotFound(w, "certificate", errors.New("certificate ID is not provided"))
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	message.SendUnknownError(w, errors.New("not implemented"))
	//TODO
}

func (h *Handler) Generate(w http.ResponseWriter, r *http.Request) {
	caReq := new(CARequest)
	if err := json.NewDecoder(r.Body).Decode(&caReq); err != nil {
		message.SendInvalidJSON(w, err)
		return
	}

	//TODO ADD VALIDATION
	ips := make([]net.IP, 0)
	for _, v := range caReq.IPs {
		ip := net.ParseIP(v)
		ips = append(ips, ip)
	}

	if caReq.CA != nil {
		pi, err := h.svc.GenerateFromCA(r.Context(), caReq.CA)
		if err != nil {
			logrus.Errorf("pki: %v", err)
			message.SendUnknownError(w, err)
			return
		}

		if err := json.NewEncoder(w).Encode(pi); err != nil {
			message.SendUnknownError(w, err)
			return
		}
	} else {
		pi, err := h.svc.GenerateSelfSigned(r.Context())
		if err != nil {
			logrus.Errorf("pki: %v", err)
			message.SendUnknownError(w, err)
			return
		}

		if err := json.NewEncoder(w).Encode(pi); err != nil {
			message.SendUnknownError(w, err)
			return
		}
	}
}
