package api

import (
	"github.com/supergiant/supergiant/core"

	"github.com/gorilla/mux"
)

func NewRouter(core *core.Core) *mux.Router {
	// StrictSlash will redirect /apps to /apps/
	// otherwise mux will simply not match /apps/
	r := mux.NewRouter()
	r.StrictSlash(true)

	s := r.PathPrefix("/v0").Subrouter()

	imageRepos := &ImageRepoController{core}
	entrypoints := &EntrypointController{core}
	apps := &AppController{core}
	components := &ComponentController{core}
	releases := &ReleaseController{core}
	instances := &InstanceController{core}
	tasks := &TaskController{core}

	s.HandleFunc("/registries/dockerhub/repos", imageRepos.Create).Methods("POST")
	s.HandleFunc("/registries/dockerhub/repos/{name}", imageRepos.Delete).Methods("DELETE")

	s.HandleFunc("/entrypoints", entrypoints.Create).Methods("POST")
	s.HandleFunc("/entrypoints", entrypoints.Index).Methods("GET")
	s.HandleFunc("/entrypoints/{domain}", entrypoints.Show).Methods("GET")
	s.HandleFunc("/entrypoints/{domain}", entrypoints.Delete).Methods("DELETE")

	s.HandleFunc("/apps", apps.Create).Methods("POST")
	s.HandleFunc("/apps", apps.Index).Methods("GET")
	s.HandleFunc("/apps/{app_name}", apps.Show).Methods("GET")
	s.HandleFunc("/apps/{app_name}", apps.Delete).Methods("DELETE")

	s.HandleFunc("/apps/{app_name}/components", components.Create).Methods("POST")
	s.HandleFunc("/apps/{app_name}/components", components.Index).Methods("GET")
	s.HandleFunc("/apps/{app_name}/components/{comp_name}", components.Show).Methods("GET")
	s.HandleFunc("/apps/{app_name}/components/{comp_name}", components.Delete).Methods("DELETE")

	s.HandleFunc("/apps/{app_name}/components/{comp_name}/releases", releases.Create).Methods("POST")
	s.HandleFunc("/apps/{app_name}/components/{comp_name}/releases", releases.Index).Methods("GET")
	s.HandleFunc("/apps/{app_name}/components/{comp_name}/releases/{release_timestamp}", releases.Show).Methods("GET")

	// Integration

	s.HandleFunc("/apps/{app_name}/components/{comp_name}/releases/{release_timestamp}/instances", instances.Index).Methods("GET")
	s.HandleFunc("/apps/{app_name}/components/{comp_name}/releases/{release_timestamp}/instances/{instance_id}", instances.Show).Methods("GET")

	s.HandleFunc("/apps/{app_name}/components/{comp_name}/releases/{release_timestamp}/instances/{instance_id}/start", instances.Start).Methods("POST")
	s.HandleFunc("/apps/{app_name}/components/{comp_name}/releases/{release_timestamp}/instances/{instance_id}/stop", instances.Stop).Methods("POST")

	// Misc

	s.HandleFunc("/tasks", tasks.Index).Methods("GET")
	s.HandleFunc("/tasks/{id}", tasks.Show).Methods("GET")

	return r
}
