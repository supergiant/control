package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	sgjwt "github.com/supergiant/supergiant/pkg/jwt"
)

type AuthHandler struct {
	userService  Service
	tokenService sgjwt.TokenService
}

type authRequest struct {
	UserName string
	Password string
}

func NewAuthHandler(userService Service, tokenService sgjwt.TokenService) *AuthHandler {
	return &AuthHandler{
		userService:  userService,
		tokenService: tokenService,
	}
}

func (m *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var ar authRequest

	err := json.NewDecoder(r.Body).Decode(&ar)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := m.userService.Authenticate(r.Context(), ar.UserName, ar.Password); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	if token, err := m.tokenService.Issue(ar.UserName); err == nil {
		w.Header().Set("Authorization", token)
		return
	} else {
		http.Error(w, fmt.Sprintf("Error while generating token %s", err.Error()), http.StatusInternalServerError)
		return
	}
}
