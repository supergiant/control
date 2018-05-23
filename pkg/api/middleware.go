package api

import (
	"net/http"
	"github.com/supergiant/supergiant/pkg/user"
)


type middleware struct {
	Users user.Service
}

func (m *middleware) authorisationHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		//TODO Implement JWT
		next.ServeHTTP(rw, r)
	})
}
