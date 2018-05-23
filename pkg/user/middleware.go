package user

import (
	"net/http"

	sgjwt "github.com/supergiant/supergiant/pkg/jwt"
)

const supergiantAuthHeader = "SGTOKEN"

func Authenticate(tokenService sgjwt.TokenService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get(supergiantAuthHeader)
		err := tokenService.Validate(tokenString)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
