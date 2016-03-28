package controller

import (
	"github.com/supergiant/supergiant/core"

	"github.com/gorilla/mux"
)

func NewRouter(client *core.Client) *mux.Router {
	// StrictSlash will redirect /apps to /apps/
	// otherwise mux will simply not match /apps/
	r := mux.NewRouter()
	r.StrictSlash(true)

	apps := &AppController{client}
	components := &ComponentController{client}
	// deployments := &DeploymentController{client}
	imageRepos := &ImageRepoController{client}
	// instances := &InstanceController{client}
	// releases := &ReleaseController{client}

	r.HandleFunc("/registries/dockerhub/repos", imageRepos.Create).Methods("POST")
	r.HandleFunc("/registries/dockerhub/repos/{name}", imageRepos.Delete).Methods("DELETE")

	r.HandleFunc("/apps", apps.Create).Methods("POST")
	r.HandleFunc("/apps", apps.Index).Methods("GET")
	r.HandleFunc("/apps/{name}", apps.Show).Methods("GET")
	r.HandleFunc("/apps/{name}", apps.Delete).Methods("DELETE")

	r.HandleFunc("/apps/{app_name}/components", components.Create).Methods("POST")
	r.HandleFunc("/apps/{app_name}/components", components.Index).Methods("GET")
	r.HandleFunc("/apps/{app_name}/components/{name}", components.Show).Methods("GET")
	r.HandleFunc("/apps/{app_name}/components/{name}", components.Delete).Methods("DELETE")

	// // --------------------------------------------------------------------------------------------------------------------------
	// //
	// // RELEASE
	// // {
	// //   strategy: "rolling_restart",
	// //   change: {
	// //     ...partial attrs (to merge) ---- for arrays, all members must be included
	// //   }
	// // }
	//
	//
	// r.HandleFunc("/apps/{app_name}/components/{comp_name}/releases", releases.Create).Methods("POST")
	// r.HandleFunc("/apps/{app_name}/components/{comp_name}/releases", releases.Index).Methods("GET")
	// r.HandleFunc("/apps/{app_name}/components/{comp_name}/releases/{id}", releases.Show).Methods("GET")
	//
	//
	//
	//
	// // This is where all the integration happens -----------------------------------------------------------------------------------------
	// r.HandleFunc("/deployments/{id}", deployments.Delete).Methods("DELETE")
	//
	// r.HandleFunc("/deployments/{deployment_id}/instances", instances.Index).Methods("GET")
	// r.HandleFunc("/deployments/{deployment_id}/instances/{id}", instances.Show).Methods("GET")
	// r.HandleFunc("/deployments/{deployment_id}/instances/{id}", instances.Delete).Methods("DELETE") // restart

	return r
}
