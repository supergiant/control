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

	imageRepos := &ImageRepoController{core}
	apps := &AppController{core}
	components := &ComponentController{core}
	releases := &ReleaseController{core}
	instances := &InstanceController{core}
	tasks := &TaskController{core}

	// deploys := &DeployController{core}

	r.HandleFunc("/registries/dockerhub/repos", imageRepos.Create).Methods("POST")
	r.HandleFunc("/registries/dockerhub/repos/{name}", imageRepos.Delete).Methods("DELETE")

	r.HandleFunc("/apps", apps.Create).Methods("POST")
	r.HandleFunc("/apps", apps.Index).Methods("GET")
	r.HandleFunc("/apps/{app_name}", apps.Show).Methods("GET")
	r.HandleFunc("/apps/{app_name}", apps.Delete).Methods("DELETE")

	r.HandleFunc("/apps/{app_name}/components", components.Create).Methods("POST")
	r.HandleFunc("/apps/{app_name}/components", components.Index).Methods("GET")
	r.HandleFunc("/apps/{app_name}/components/{comp_name}", components.Show).Methods("GET")
	r.HandleFunc("/apps/{app_name}/components/{comp_name}", components.Delete).Methods("DELETE")

	r.HandleFunc("/apps/{app_name}/components/{comp_name}/releases", releases.Create).Methods("POST")
	r.HandleFunc("/apps/{app_name}/components/{comp_name}/releases", releases.Index).Methods("GET")
	r.HandleFunc("/apps/{app_name}/components/{comp_name}/releases/{release_id}", releases.Show).Methods("GET")

	// Below is where all the integration happens
	//============================================================================

	r.HandleFunc("/apps/{app_name}/components/{comp_name}/releases/{release_id}/instances", instances.Index).Methods("GET")
	r.HandleFunc("/apps/{app_name}/components/{comp_name}/releases/{release_id}/instances/{instance_id}", instances.Show).Methods("GET")

	r.HandleFunc("/apps/{app_name}/components/{comp_name}/releases/{release_id}/instances/{instance_id}/start", instances.Start).Methods("POST")
	r.HandleFunc("/apps/{app_name}/components/{comp_name}/releases/{release_id}/instances/{instance_id}/stop", instances.Stop).Methods("POST")

	r.HandleFunc("/tasks", tasks.Index).Methods("GET")
	r.HandleFunc("/tasks/{id}", tasks.Show).Methods("GET")

	return r
}
