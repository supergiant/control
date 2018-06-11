package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/asaskevich/govalidator.v8"

	sgjwt "github.com/supergiant/supergiant/pkg/jwt"
)

type Endpoint struct {
	userService  Service
	tokenService sgjwt.TokenService
}

type authRequest struct {
	UserName string
	Password string
}

/**
s.HandleFunc("/users", restrictedHandler(core, CreateUser)).Methods("POST")
s.HandleFunc("/users", restrictedHandler(core, ListUsers)).Methods("GET")
s.HandleFunc("/users/{id}", restrictedHandler(core, GetUser)).Methods("GET")
s.HandleFunc("/users/{id}", restrictedHandler(core, UpdateUser)).Methods("PATCH", "PUT")
s.HandleFunc("/users/{id}", restrictedHandler(core, DeleteUser)).Methods("DELETE")
*/

func NewEndpoint(userService Service, tokenService sgjwt.TokenService) *Endpoint {
	return &Endpoint{
		userService:  userService,
		tokenService: tokenService,
	}
}

func (e *Endpoint) Authenticate(w http.ResponseWriter, r *http.Request) {
	var ar authRequest
	err := json.NewDecoder(r.Body).Decode(&ar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := e.userService.Authenticate(r.Context(), ar.UserName, ar.Password); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	if token, err := e.tokenService.Issue(ar.UserName); err == nil {
		w.Header().Set("Authorization", token)
		return
	} else {
		http.Error(w, fmt.Sprintf("Error while generating token %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (e *Endpoint) Create(rw http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	ok, err := govalidator.ValidateStruct(user)
	if !ok {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if err := e.userService.Create(r.Context(), &user); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
