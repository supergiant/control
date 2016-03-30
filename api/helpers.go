package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/supergiant/supergiant/core"

	"github.com/gorilla/mux"
)

// LoadApp loads an App resource from URL params, or renders an HTTP Bad Request
// error.
func LoadApp(core *core.Core, w http.ResponseWriter, r *http.Request) (*core.AppResource, error) {
	name := mux.Vars(r)["app_name"]
	app, err := core.Apps().Get(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return nil, err
	}

	return app, nil
}

// LoadComponent loads an Component resource from URL params, or renders an HTTP
// Bad Request error.
func LoadComponent(core *core.Core, w http.ResponseWriter, r *http.Request) (*core.ComponentResource, error) {
	app, err := LoadApp(core, w, r)
	if err != nil {
		return nil, err
	}

	name := mux.Vars(r)["comp_name"]
	component, err := app.Components().Get(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return nil, err
	}

	return component, nil
}

// LoadRelease loads an Release resource from URL params, or renders an HTTP
// Bad Request error.
func LoadRelease(core *core.Core, w http.ResponseWriter, r *http.Request) (*core.ReleaseResource, error) {
	component, err := LoadComponent(core, w, r)
	if err != nil {
		return nil, err
	}

	id := mux.Vars(r)["release_id"]
	release, err := component.Releases().Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return nil, err
	}

	return release, nil
}

// LoadInstance loads an Instance resource from URL params, or renders an HTTP
// Bad Request error.
func LoadInstance(core *core.Core, w http.ResponseWriter, r *http.Request) (*core.InstanceResource, error) {
	release, err := LoadRelease(core, w, r)
	if err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(mux.Vars(r)["instance_id"]) // TODO instance ID should be a string to begin with
	if err != nil {
		return nil, err
	}
	instance, err := release.Instances().Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return nil, err
	}

	return instance, nil
}

// LoadImageRepo loads an ImageRepo resource from URL params, or renders an HTTP
// Bad Request error.
func LoadImageRepo(core *core.Core, w http.ResponseWriter, r *http.Request) (*core.ImageRepoResource, error) {
	name := mux.Vars(r)["name"]
	repo, err := core.ImageRepos().Get(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return nil, err
	}

	return repo, nil
}

// LoadTask loads an Task resource from URL params, or renders an HTTP
// Bad Request error.
func LoadTask(core *core.Core, w http.ResponseWriter, r *http.Request) (*core.TaskResource, error) {
	id := mux.Vars(r)["id"]
	task, err := core.Tasks().Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return nil, err
	}

	return task, nil
}

// UnmarshalBodyInto decodes a JSON body into an interface or renders an HTTP
// Bad Request error.
func UnmarshalBodyInto(w http.ResponseWriter, r *http.Request, out interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(out)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	return nil
}

// MarshalBody marshals an interface into a JSON body or renders an HTTP Bad
// Request error.
func MarshalBody(w http.ResponseWriter, in interface{}) (string, error) {
	out, err := json.MarshalIndent(in, "", "  ")
	// out, err := json.Marshal(in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}
	return string(out) + "\n", nil
}

// RenderWithStatusAccepted renders a response with HTTP status 202.
func RenderWithStatusAccepted(w http.ResponseWriter, body string) {
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprint(w, body)
}

// RenderWithStatusCreated renders a response with HTTP status 201.
func RenderWithStatusCreated(w http.ResponseWriter, body string) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, body)
}

// RenderWithStatusOK renders a response with HTTP status 200.
func RenderWithStatusOK(w http.ResponseWriter, body string) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, body)
}
