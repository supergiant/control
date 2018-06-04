package api

import (
	"net/http"

	"strings"

	sgjwt "github.com/supergiant/supergiant/pkg/jwt"
	"github.com/supergiant/supergiant/internal/user"
)

type middleware struct {
	TokenService *sgjwt.TokenService
	UserService  *user.Service
}

func (m *middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get rid of Bearer
		tokenString := strings.Split(r.Header.Get("Authorization"), " ")[1]
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
