package user

import (
	"net/http"

	"strings"

	sgjwt "github.com/supergiant/supergiant/pkg/jwt"
)

func AuthMiddleware(tokenService sgjwt.TokenService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get rid of Bearer
		tokenString := strings.Split(r.Header.Get("Authorization"), " ")[0]
		claims, err := tokenService.Validate(tokenString)

		// TODO(stgleb): Do something with claims
		userId := claims["user_id"].(string)

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
