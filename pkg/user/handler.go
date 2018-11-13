package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/asaskevich/govalidator.v8"

	"github.com/supergiant/supergiant/pkg/message"
	"github.com/supergiant/supergiant/pkg/sgerrors"
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
			return
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
			msg := message.New(fmt.Sprintf("Login %s is already occupied", user.Login), "", sgerrors.EntityAlreadyExists, "")
			message.SendMessage(rw, msg, http.StatusBadRequest)
			return
		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
