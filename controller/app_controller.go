package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"supergiant/core"

	"github.com/gorilla/mux"
)

type AppController struct {
	client *core.Client
}

func (c *AppController) Create(w http.ResponseWriter, r *http.Request) {
	app := c.client.Apps().New()

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app, err = c.client.Apps().Create(app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// namespace := &guber.Namespace{
	// 	Metadata: &guber.Metadata{
	// 		Name: app.Name,
	// 	},
	// }
	// namespace, err = e.kube.Namespaces().Create(namespace)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	out, err := json.Marshal(app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(out))
}

func (c *AppController) Index(w http.ResponseWriter, r *http.Request) {
	apps, err := c.client.Apps().List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(apps)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(out))
}

func (c *AppController) Show(w http.ResponseWriter, r *http.Request) {
	appName := mux.Vars(r)["name"]
	app, err := c.client.Apps().Get(appName)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	out, err := json.Marshal(app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(out))
}

func (c *AppController) Delete(w http.ResponseWriter, r *http.Request) {
	appName := mux.Vars(r)["name"]

	// // NOTE we don't want to return on error here, because if not found, we are
	// // okay -- if the error is at the system level, we should be panicking, not returning errors
	// e.kube.Namespaces().Delete(appName) // TODO though it should probably not run this request if the App does not exist in storage

	if err := c.client.Apps().Delete(appName); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// fmt.Fprint(w, string(out))
}
