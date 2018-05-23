package api

import (
	"net/http"
	"regexp"

	sgjwt "github.com/supergiant/supergiant/pkg/jwt"
	"github.com/supergiant/supergiant/pkg/user"
	"encoding/json"
	"fmt"
)

const supergiantAuthHeader = "SGTOKEN"

type middleware struct {
	userService  user.Service
	tokenService sgjwt.TokenService
}

type authRequest struct {
	UserName string
	Password string
}


// TODO(stgleb): move to separate handlers module
func (m *middleware) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var ar authRequest

	err := json.NewDecoder(r.Body).Decode(&ar)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ok, err := m.userService.Authenticate(r.Context(), ar.UserName, ar.Password); !ok || err != nil {
		http.Error(w, "unknown user", http.StatusForbidden)
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

func (m *middleware) AuthMiddleware(next http.Handler) http.Handler {
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

//FIXME Reusing the old code to keep frontend from being broken
func (m *middleware) authorisationHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		tokenMatch := regexp.MustCompile(`^SGAPI (token|session)="([A-Za-z0-9]{32})"$`).FindStringSubmatch(auth)

		if len(tokenMatch) != 3 {
			respond(rw, nil, errorBadAuthHeader)
		}
		switch tokenMatch[1] {
		case "token":
			token := tokenMatch[2]
			_, err := m.userService.GetByToken(r.Context(), token)
			if err != nil {
				respond(rw, nil, errorUnauthorized)
				return
			}
			//return user
		case "session":
			session := tokenMatch[3]
			if err := m.userService.GetBySession(r.Context(), session); err != nil {
				respond(rw, nil, errorUnauthorized)
				return
			}
			return
		default:
			respond(rw, nil, errorBadAuthHeader)
			return
		}
		next.ServeHTTP(rw, r)
	})
}
