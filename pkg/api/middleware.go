package api

import (
	"net/http"
	"regexp"

	"github.com/supergiant/supergiant/pkg/user"
)

const tokenRegexp = `SGAPI (token|session)="([A-Za-z0-9]{32})"$`

type middleware struct {
	Users user.Service
}

//FIXME Reusing the old code to keep frontend from being broken
func (m *middleware) authorisationHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		tokenMatch := regexp.MustCompile(tokenRegexp).FindStringSubmatch(auth)

		if len(tokenMatch) != 3 {
			respond(rw, nil, errorBadAuthHeader)
		}
		switch tokenMatch[1] {
		case "token":
			token := tokenMatch[2]
			_, err := m.Users.GetByToken(r.Context(), token)
			if err != nil {
				respond(rw, nil, errorUnauthorized)
				return
			}
			//return user
		case "session":
			session := tokenMatch[3]
			if err := m.Users.GetBySession(r.Context(), session); err != nil {
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
