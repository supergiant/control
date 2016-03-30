package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/supergiant/supergiant/core"

	"github.com/gorilla/mux"
)

// loadApp loads an App resource from URL params, or renders an HTTP Not Found
// error.
func loadApp(core *core.Core, w http.ResponseWriter, r *http.Request) (*core.AppResource, error) {
	name := mux.Vars(r)["app_name"]
	app, err := core.Apps().Get(&name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return nil, err
	}

	return app, nil
}

// loadComponent loads an Component resource from URL params, or renders an HTTP
// Not Found error.
func loadComponent(core *core.Core, w http.ResponseWriter, r *http.Request) (*core.ComponentResource, error) {
	app, err := loadApp(core, w, r)
	if err != nil {
		return nil, err
	}

	name := mux.Vars(r)["comp_name"]
	component, err := app.Components().Get(&name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return nil, err
	}

	return component, nil
}

// loadRelease loads an Release resource from URL params, or renders an HTTP
// Not Found error.
func loadRelease(core *core.Core, w http.ResponseWriter, r *http.Request) (*core.ReleaseResource, error) {
	component, err := loadComponent(core, w, r)
	if err != nil {
		return nil, err
	}

	timestamp := mux.Vars(r)["release_timestamp"]
	release, err := component.Releases().Get(&timestamp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return nil, err
	}

	return release, nil
}

// loadInstance loads an Instance resource from URL params, or renders an HTTP
// Not Found error.
func loadInstance(core *core.Core, w http.ResponseWriter, r *http.Request) (*core.InstanceResource, error) {
	release, err := loadRelease(core, w, r)
	if err != nil {
		return nil, err
	}

	id := mux.Vars(r)["instance_id"]
	if err != nil {
		return nil, err
	}
	instance, err := release.Instances().Get(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return nil, err
	}

	return instance, nil
}

// loadImageRepo loads an ImageRepo resource from URL params, or renders an HTTP
// Not Found error.
func loadImageRepo(core *core.Core, w http.ResponseWriter, r *http.Request) (*core.ImageRepoResource, error) {
	name := mux.Vars(r)["name"]
	repo, err := core.ImageRepos().Get(&name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return nil, err
	}

	return repo, nil
}

// loadEntrypoint loads an Entrypoint resource from URL params, or renders an HTTP
// Not Found error.
func loadEntrypoint(core *core.Core, w http.ResponseWriter, r *http.Request) (*core.EntrypointResource, error) {
	domain := mux.Vars(r)["domain"]
	entrypoint, err := core.Entrypoints().Get(&domain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return nil, err
	}

	return entrypoint, nil
}

// loadTask loads an Task resource from URL params, or renders an HTTP
// Not Found error.
func loadTask(core *core.Core, w http.ResponseWriter, r *http.Request) (*core.TaskResource, error) {
	id := mux.Vars(r)["id"]
	task, err := core.Tasks().Get(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return nil, err
	}

	return task, nil
}

// unmarshalBodyInto decodes a JSON body into an interface or renders an HTTP
// Not Found error.
func unmarshalBodyInto(w http.ResponseWriter, r *http.Request, out interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(out)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	return nil
}

// marshalBody marshals an interface into a JSON body or renders an HTTP Bad
// Request error.
func marshalBody(w http.ResponseWriter, in interface{}) (string, error) {
	out, err := json.MarshalIndent(in, "", "  ")
	// out, err := json.Marshal(in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}
	return string(out) + "\n", nil
}

// renderWithStatusAccepted renders a response with HTTP status 202.
func renderWithStatusAccepted(w http.ResponseWriter, body string) {
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprint(w, body)
}

// renderWithStatusCreated renders a response with HTTP status 201.
func renderWithStatusCreated(w http.ResponseWriter, body string) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, body)
}

// renderWithStatusOK renders a response with HTTP status 200.
func renderWithStatusOK(w http.ResponseWriter, body string) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, body)
}
