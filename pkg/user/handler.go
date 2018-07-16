package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/asaskevich/govalidator.v8"

	sgjwt "github.com/supergiant/supergiant/pkg/jwt"
	"github.com/supergiant/supergiant/pkg/message"
	"github.com/supergiant/supergiant/pkg/sgerrors"
)

type Handler struct {
	userService  *Service
	tokenService *sgjwt.TokenService
}

type authRequest struct {
	UserName string
	Password string
}

func NewHandler(userService *Service, tokenService *sgjwt.TokenService) *Handler {
	return &Handler{
		userService:  userService,
		tokenService: tokenService,
	}
}

func (h *Handler) Authenticate(w http.ResponseWriter, r *http.Request) {
	var ar authRequest
	if err := json.NewDecoder(r.Body).Decode(&ar); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.userService.Authenticate(r.Context(), ar.UserName, ar.Password); err != nil {
		if sgerrors.IsInvalidCredentials(err) {
			http.Error(w, sgerrors.ErrInvalidCredentials.Error(), http.StatusForbidden)
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if token, err := h.tokenService.Issue(ar.UserName); err == nil {
		w.Header().Set("Authorization", token)
		return
	} else {
		http.Error(w, fmt.Sprintf("Error while generating token %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Create(rw http.ResponseWriter, r *http.Request) {
	var user User
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
