package api

import (
	"net/http"
	"strings"

	sgjwt "github.com/supergiant/supergiant/pkg/jwt"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/user"
)

type Middleware struct {
	TokenService *sgjwt.TokenService
	UserService  *user.Service
}

func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, sgerrors.ErrInvalidCredentials.Error(), http.StatusForbidden)
			return
		}

		ts := strings.Split(authHeader, " ")
		if len(ts) <= 1 {
			http.Error(w, sgerrors.ErrInvalidCredentials.Error(), http.StatusForbidden)
			return
		}

		tokenString := ts[1]
		claims, err := m.TokenService.Validate(tokenString)

		// TODO(stgleb): Do something with claims
		userId, ok := claims["user_id"].(string)
		if !ok {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		if len(userId) == 0 {
			http.Error(w, "unknown user", http.StatusForbidden)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ContentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
