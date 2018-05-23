package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	sgjwt "github.com/supergiant/supergiant/pkg/jwt"
	"github.com/supergiant/supergiant/pkg/user"
)

const supergiantAuthHeader = "SGTOKEN"

type AuthMiddleware struct {
	userService  user.Service
	tokenService sgjwt.TokenService
}

type authRequest struct {
	UserName string
	Password string
}

// TODO(stgleb): move to separate handlers module
func (m *AuthMiddleware) LoginHandler(w http.ResponseWriter, r *http.Request) {
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
		w.Header().Set(supergiantAuthHeader, token)
		return
	} else {
		http.Error(w, fmt.Sprintf("Error while generating token %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get(supergiantAuthHeader)
		err := m.tokenService.Validate(tokenString)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
