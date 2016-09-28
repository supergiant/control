package api

import (
	"github.com/supergiant/supergiant/pkg/core"

	"github.com/gorilla/mux"
)

func NewRouter(core *core.Core) *mux.Router {
	r := mux.NewRouter()

	s := r.PathPrefix("/api/v0").Subrouter()

	// Login request can't be authenticated
	s.HandleFunc("/sessions", openHandler(core, CreateSession)).Methods("POST")

	s.HandleFunc("/sessions/{id}", restrictedHandler(core, GetSession)).Methods("GET")
	s.HandleFunc("/sessions", restrictedHandler(core, ListSessions)).Methods("GET")
	s.HandleFunc("/sessions/{id}", restrictedHandler(core, DeleteSession)).Methods("DELETE")

	s.HandleFunc("/users", restrictedHandler(core, CreateUser)).Methods("POST")
	s.HandleFunc("/users", restrictedHandler(core, ListUsers)).Methods("GET")
	s.HandleFunc("/users/{id}", restrictedHandler(core, GetUser)).Methods("GET")
	s.HandleFunc("/users/{id}", restrictedHandler(core, UpdateUser)).Methods("PATCH", "PUT")
	s.HandleFunc("/users/{id}", restrictedHandler(core, DeleteUser)).Methods("DELETE")
	s.HandleFunc("/users/{id}/regenerate_api_token", restrictedHandler(core, RegenerateUserAPIToken)).Methods("POST")

	s.HandleFunc("/cloud_accounts", restrictedHandler(core, CreateCloudAccount)).Methods("POST")
	s.HandleFunc("/cloud_accounts", restrictedHandler(core, ListCloudAccounts)).Methods("GET")
	s.HandleFunc("/cloud_accounts/{id}", restrictedHandler(core, GetCloudAccount)).Methods("GET")
	s.HandleFunc("/cloud_accounts/{id}", restrictedHandler(core, UpdateCloudAccount)).Methods("PATCH", "PUT")
	s.HandleFunc("/cloud_accounts/{id}", restrictedHandler(core, DeleteCloudAccount)).Methods("DELETE")

	s.HandleFunc("/kubes", restrictedHandler(core, CreateKube)).Methods("POST")
	s.HandleFunc("/kubes", restrictedHandler(core, ListKubes)).Methods("GET")
	s.HandleFunc("/kubes/{id}", restrictedHandler(core, GetKube)).Methods("GET")
	s.HandleFunc("/kubes/{id}", restrictedHandler(core, UpdateKube)).Methods("PATCH", "PUT")
	s.HandleFunc("/kubes/{id}", restrictedHandler(core, DeleteKube)).Methods("DELETE")

	s.HandleFunc("/kube_resources", restrictedHandler(core, CreateKubeResource)).Methods("POST")
	s.HandleFunc("/kube_resources", restrictedHandler(core, ListKubeResources)).Methods("GET")
	s.HandleFunc("/kube_resources/{id}", restrictedHandler(core, GetKubeResource)).Methods("GET")
	s.HandleFunc("/kube_resources/{id}", restrictedHandler(core, UpdateKubeResource)).Methods("PATCH", "PUT")
	s.HandleFunc("/kube_resources/{id}/start", restrictedHandler(core, StartKubeResource)).Methods("POST")
	s.HandleFunc("/kube_resources/{id}/stop", restrictedHandler(core, StopKubeResource)).Methods("POST")
	s.HandleFunc("/kube_resources/{id}", restrictedHandler(core, DeleteKubeResource)).Methods("DELETE")

	s.HandleFunc("/nodes", restrictedHandler(core, CreateNode)).Methods("POST")
	s.HandleFunc("/nodes", restrictedHandler(core, ListNodes)).Methods("GET")
	s.HandleFunc("/nodes/{id}", restrictedHandler(core, GetNode)).Methods("GET")
	s.HandleFunc("/nodes/{id}", restrictedHandler(core, UpdateNode)).Methods("PATCH", "PUT")
	s.HandleFunc("/nodes/{id}", restrictedHandler(core, DeleteNode)).Methods("DELETE")

	s.HandleFunc("/volumes", restrictedHandler(core, CreateVolume)).Methods("POST")
	s.HandleFunc("/volumes", restrictedHandler(core, ListVolumes)).Methods("GET")
	s.HandleFunc("/volumes/{id}", restrictedHandler(core, GetVolume)).Methods("GET")
	s.HandleFunc("/volumes/{id}", restrictedHandler(core, UpdateVolume)).Methods("PATCH", "PUT")
	s.HandleFunc("/volumes/{id}", restrictedHandler(core, DeleteVolume)).Methods("DELETE")

	s.HandleFunc("/entrypoints", restrictedHandler(core, CreateEntrypoint)).Methods("POST")
	s.HandleFunc("/entrypoints", restrictedHandler(core, ListEntrypoints)).Methods("GET")
	s.HandleFunc("/entrypoints/{id}", restrictedHandler(core, GetEntrypoint)).Methods("GET")
	s.HandleFunc("/entrypoints/{id}", restrictedHandler(core, UpdateEntrypoint)).Methods("PATCH", "PUT")
	s.HandleFunc("/entrypoints/{id}", restrictedHandler(core, DeleteEntrypoint)).Methods("DELETE")

	s.HandleFunc("/entrypoint_listeners", restrictedHandler(core, CreateEntrypointListener)).Methods("POST")
	s.HandleFunc("/entrypoint_listeners", restrictedHandler(core, ListEntrypointListeners)).Methods("GET")
	s.HandleFunc("/entrypoint_listeners/{id}", restrictedHandler(core, GetEntrypointListener)).Methods("GET")
	s.HandleFunc("/entrypoint_listeners/{id}", restrictedHandler(core, DeleteEntrypointListener)).Methods("DELETE")

	s.HandleFunc("/log", logHandler(core)).Methods("GET")

	return r
}
