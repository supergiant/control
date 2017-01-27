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
	s.HandleFunc("/kubes/{id}/provision", restrictedHandler(core, ProvisionKube)).Methods("POST")
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

	s.HandleFunc("/load_balancers", restrictedHandler(core, CreateLoadBalancer)).Methods("POST")
	s.HandleFunc("/load_balancers", restrictedHandler(core, ListLoadBalancers)).Methods("GET")
	s.HandleFunc("/load_balancers/{id}", restrictedHandler(core, GetLoadBalancer)).Methods("GET")
	s.HandleFunc("/load_balancers/{id}", restrictedHandler(core, UpdateLoadBalancer)).Methods("PATCH", "PUT")
	s.HandleFunc("/load_balancers/{id}", restrictedHandler(core, DeleteLoadBalancer)).Methods("DELETE")

	s.HandleFunc("/helm_repos", restrictedHandler(core, CreateHelmRepo)).Methods("POST")
	s.HandleFunc("/helm_repos", restrictedHandler(core, ListHelmRepos)).Methods("GET")
	s.HandleFunc("/helm_repos/{id}", restrictedHandler(core, GetHelmRepo)).Methods("GET")
	s.HandleFunc("/helm_repos/{id}", restrictedHandler(core, UpdateHelmRepo)).Methods("PATCH", "PUT")
	s.HandleFunc("/helm_repos/{id}", restrictedHandler(core, DeleteHelmRepo)).Methods("DELETE")

	s.HandleFunc("/helm_charts", restrictedHandler(core, CreateHelmChart)).Methods("POST")
	s.HandleFunc("/helm_charts", restrictedHandler(core, ListHelmCharts)).Methods("GET")
	s.HandleFunc("/helm_charts/{id}", restrictedHandler(core, GetHelmChart)).Methods("GET")
	s.HandleFunc("/helm_charts/{id}", restrictedHandler(core, UpdateHelmChart)).Methods("PATCH", "PUT")
	s.HandleFunc("/helm_charts/{id}", restrictedHandler(core, DeleteHelmChart)).Methods("DELETE")

	s.HandleFunc("/helm_releases", restrictedHandler(core, CreateHelmRelease)).Methods("POST")
	s.HandleFunc("/helm_releases", restrictedHandler(core, ListHelmReleases)).Methods("GET")
	s.HandleFunc("/helm_releases/{id}", restrictedHandler(core, GetHelmRelease)).Methods("GET")
	s.HandleFunc("/helm_releases/{id}", restrictedHandler(core, UpdateHelmRelease)).Methods("PATCH", "PUT")
	s.HandleFunc("/helm_releases/{id}", restrictedHandler(core, DeleteHelmRelease)).Methods("DELETE")

	s.HandleFunc("/log", logHandler(core)).Methods("GET")

	return r
}
