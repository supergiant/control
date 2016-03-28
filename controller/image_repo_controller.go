package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/supergiant/supergiant/core"

	"github.com/gorilla/mux"
)

type ImageRepoController struct {
	client *core.Client
}

func (c *ImageRepoController) Create(w http.ResponseWriter, r *http.Request) {
	m := new(core.ImageRepo)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m, err = c.client.ImageRepos().Create(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(out))
}

// func (c *ImageRepoController) Index(w http.ResponseWriter, r *http.Request) {
// 	list, err := c.client.ImageRepos().List()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
//
// 	out, err := json.Marshal(list)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
//
// 	fmt.Fprint(w, string(out))
// }

// func (c *ImageRepoController) Show(w http.ResponseWriter, r *http.Request) {
// 	name := mux.Vars(r)["name"]
// 	m, err := c.client.ImageRepos().Get(name)
// 	if err != nil {
// 		http.Error(w, "Not Found", http.StatusNotFound)
// 		return
// 	}
//
// 	out, err := json.Marshal(m)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
//
// 	fmt.Fprint(w, string(out))
// }

func (c *ImageRepoController) Delete(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	if err := c.client.ImageRepos().Delete(name); err != nil { // TODO ------------- should do the same found, err := thing we did for Guber deletes
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// fmt.Fprint(w, string(out))
}
