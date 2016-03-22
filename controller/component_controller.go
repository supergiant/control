package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"supergiant/core"

	"github.com/gorilla/mux"
)

type ComponentController struct {
	client *core.Client
}

func (c *ComponentController) Create(w http.ResponseWriter, r *http.Request) {
	appName := mux.Vars(r)["app_name"]
	app, err := c.client.Apps().Get(appName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// component := new(core.Component)
	component := app.Components().New()

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(component)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	component, err = app.Components().Create(component)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(component)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(out))
}

func (c *ComponentController) Index(w http.ResponseWriter, r *http.Request) {
	appName := mux.Vars(r)["app_name"]
	app, err := c.client.Apps().Get(appName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	components, err := app.Components().List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(components)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(out))
}

func (c *ComponentController) Show(w http.ResponseWriter, r *http.Request) {
	urlVars := mux.Vars(r)
	appName := urlVars["app_name"]
	compName := urlVars["name"]
	app, err := c.client.Apps().Get(appName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	component, err := app.Components().Get(compName)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	out, err := json.Marshal(component)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(out))
}

func (c *ComponentController) Delete(w http.ResponseWriter, r *http.Request) {
	urlVars := mux.Vars(r)
	appName := urlVars["app_name"]
	compName := urlVars["name"]
	app, err := c.client.Apps().Get(appName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	component, err := app.Components().Get(compName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// NOTE / TODO don't like the inconsistency here -- all other controllers call
	// Resource.Delete(), but in this case, that method only destroys the record.
	// We don't want to destroy the record until AFTER the job has run, so I made
	// a model-level method that is named distinctly.
	if err := component.TeardownAndDelete(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fmt.Fprint(w, string(out))
}
