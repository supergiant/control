package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/supergiant/control/pkg/message"
	"github.com/supergiant/control/pkg/sgerrors"
	"gopkg.in/asaskevich/govalidator.v8"
)

type TokenIssuer interface {
	Issue(string) (string, error)
}

type Handler struct {
	userService  *Service
	tokenService TokenIssuer
}

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func NewHandler(userService *Service, tokenService TokenIssuer) *Handler {
	return &Handler{
		userService:  userService,
		tokenService: tokenService,
	}
}

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func (h *Handler) Authenticate(w http.ResponseWriter, r *http.Request) {
	var ar AuthRequest

	// not sure if this is needed after route gets protected again
	// (when default user is created)
	enableCors(w)

	if err := json.NewDecoder(r.Body).Decode(&ar); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.userService.Authenticate(r.Context(), ar.Login, ar.Password); err != nil {
		if sgerrors.IsInvalidCredentials(err) {
			http.Error(w, sgerrors.ErrInvalidCredentials.Error(), http.StatusForbidden)
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if token, err := h.tokenService.Issue(ar.Login); err == nil {
		w.Header().Set("Authorization", token)
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")
		return
	} else {
		http.Error(w, fmt.Sprintf("Error while generating token %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) RegisterRootUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		message.SendInvalidJSON(w, err)
		return
	}

	ok, err := govalidator.ValidateStruct(user)
	if !ok {
		message.SendValidationFailed(w, err)
		return
	}

	coldstart, err := h.userService.IsColdStart(r.Context())
	if err != nil {
		message.SendUnknownError(w, err)
		return
	}

	if coldstart {
		if err := h.userService.Create(r.Context(), &user); err != nil {
			message.SendUnknownError(w, err)
			return
		}
	} else {
		message.SendAlreadyExists(w, "root user", err)
		return
	}
}

func (h *Handler) IsColdStart(w http.ResponseWriter, r *http.Request) {
	coldstart, err := h.userService.IsColdStart(r.Context())
	if err != nil {
		message.SendUnknownError(w, err)
		return
	}

	resp := &struct {
		IsColdStart bool `json:"isColdStart"`
	}{
		coldstart,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		message.SendUnknownError(w, err)
	}
}

func (h *Handler) Create(rw http.ResponseWriter, r *http.Request) {
	var user User

	// not sure if this is needed after route gets protected again
	// (when default user is created)
	enableCors(rw)

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	ok, err := govalidator.ValidateStruct(user)
	if !ok {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.userService.Create(r.Context(), &user); err != nil {
		if sgerrors.IsAlreadyExists(err) {
			msg := message.New(fmt.Sprintf("login %s is already occupied", user.Login), "", sgerrors.EntityAlreadyExists, "")
			message.SendMessage(rw, msg, http.StatusBadRequest)
			return
		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
