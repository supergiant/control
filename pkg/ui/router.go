package ui

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
	"github.com/supergiant/supergiant/bindata"
	"github.com/supergiant/supergiant/pkg/client"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

var templates = make(map[string]*template.Template)

func init() {
	// AssetNames() is like a filepath Glob of our generated assets
	viewFiles := bindata.AssetNames()
	var partials []string
	var layouts []string
	var views []string
	for _, viewFile := range viewFiles {
		if !strings.HasPrefix(viewFile, "ui/views") {
			continue // don't want to load any non-HTML assets into this
		}
		if strings.HasPrefix(viewFile, "ui/views/partials") {
			partials = append(partials, viewFile)
		} else if strings.HasPrefix(viewFile, "ui/views/layouts") {
			layouts = append(layouts, viewFile)
		} else {
			views = append(views, viewFile)
		}
	}

	for _, view := range views {
		// https://golang.org/src/html/template/template.go?s=11938:12007#L354
		var t *template.Template
		for _, view := range append(layouts, append(partials, view)...) {
			name := filepath.Base(view)
			if t == nil {
				t = template.New(name)
			}
			var tmpl *template.Template
			if name == t.Name() {
				tmpl = t
			} else {
				tmpl = t.New(name)
			}

			src, err := bindata.Asset(view)
			if err != nil {
				panic(err)
			}

			if _, err = tmpl.Parse(string(src)); err != nil {
				panic(err)
			}
		}

		key := regexp.MustCompile(`([^/]+)\.html$`).FindStringSubmatch(view)[1]
		templates[key] = t
	}
}

func NewRouter(c *core.Core, baseRouter *mux.Router) *mux.Router {
	base := baseRouter.StrictSlash(true)

	// Redirect / to /ui
	base.HandleFunc("/", uiRedirect).Methods("GET")

	r := base.PathPrefix("/ui").Subrouter()

	// Assets
	assetDir := &assetfs.AssetFS{Asset: bindata.Asset, AssetDir: bindata.AssetDir, AssetInfo: bindata.AssetInfo}
	r.PathPrefix("/assets/").Handler(http.FileServer(assetDir))

	r.HandleFunc("/", restrictedHandler(c, Root)).Methods("GET")

	r.HandleFunc("/sessions/new", openHandler(c, NewSession)).Methods("GET")
	r.HandleFunc("/sessions", openHandler(c, CreateSession)).Methods("POST")
	r.HandleFunc("/sessions/{id}", openHandler(c, GetSession)).Methods("GET")

	r.HandleFunc("/sessions", restrictedHandler(c, ListSessions)).Methods("GET")

	r.HandleFunc("/users/new", restrictedHandler(c, NewUser)).Methods("GET")
	r.HandleFunc("/users", restrictedHandler(c, CreateUser)).Methods("POST")
	r.HandleFunc("/users", restrictedHandler(c, ListUsers)).Methods("GET")
	r.HandleFunc("/users/{id}", restrictedHandler(c, GetUser)).Methods("GET")
	r.HandleFunc("/users/{id}/edit", restrictedHandler(c, EditUser)).Methods("GET")
	r.HandleFunc("/users/{id}", restrictedHandler(c, UpdateUser)).Methods("POST")

	r.HandleFunc("/cloud_accounts/new", restrictedHandler(c, NewCloudAccount)).Methods("GET")
	r.HandleFunc("/cloud_accounts", restrictedHandler(c, CreateCloudAccount)).Methods("POST")
	r.HandleFunc("/cloud_accounts", restrictedHandler(c, ListCloudAccounts)).Methods("GET")
	r.HandleFunc("/cloud_accounts/{id}", restrictedHandler(c, GetCloudAccount)).Methods("GET")

	r.HandleFunc("/kubes/new", restrictedHandler(c, NewKube)).Methods("GET")
	r.HandleFunc("/kubes", restrictedHandler(c, CreateKube)).Methods("POST")
	r.HandleFunc("/kubes", restrictedHandler(c, ListKubes)).Methods("GET")
	r.HandleFunc("/kubes/{id}", restrictedHandler(c, GetKube)).Methods("GET")

	r.HandleFunc("/nodes/new", restrictedHandler(c, NewNode)).Methods("GET")
	r.HandleFunc("/nodes", restrictedHandler(c, CreateNode)).Methods("POST")
	r.HandleFunc("/nodes", restrictedHandler(c, ListNodes)).Methods("GET")
	r.HandleFunc("/nodes/{id}", restrictedHandler(c, GetNode)).Methods("GET")

	r.HandleFunc("/kube_resources/new", restrictedHandler(c, NewKubeResource)).Methods("GET")
	r.HandleFunc("/kube_resources", restrictedHandler(c, CreateKubeResource)).Methods("POST")
	r.HandleFunc("/kube_resources", restrictedHandler(c, ListKubeResources)).Methods("GET")
	r.HandleFunc("/kube_resources/{id}", restrictedHandler(c, GetKubeResource)).Methods("GET")
	r.HandleFunc("/kube_resources/{id}/edit", restrictedHandler(c, EditKubeResource)).Methods("GET")
	r.HandleFunc("/kube_resources/{id}", restrictedHandler(c, UpdateKubeResource)).Methods("POST")

	r.HandleFunc("/volumes/new", restrictedHandler(c, NewVolume)).Methods("GET")
	r.HandleFunc("/volumes", restrictedHandler(c, CreateVolume)).Methods("POST")
	r.HandleFunc("/volumes", restrictedHandler(c, ListVolumes)).Methods("GET")
	r.HandleFunc("/volumes/{id}", restrictedHandler(c, GetVolume)).Methods("GET")
	r.HandleFunc("/volumes/{id}/edit", restrictedHandler(c, EditVolume)).Methods("GET")
	r.HandleFunc("/volumes/{id}", restrictedHandler(c, UpdateVolume)).Methods("POST")

	r.HandleFunc("/entrypoints/new", restrictedHandler(c, NewEntrypoint)).Methods("GET")
	r.HandleFunc("/entrypoints", restrictedHandler(c, CreateEntrypoint)).Methods("POST")
	r.HandleFunc("/entrypoints", restrictedHandler(c, ListEntrypoints)).Methods("GET")
	r.HandleFunc("/entrypoints/{id}", restrictedHandler(c, GetEntrypoint)).Methods("GET")

	r.HandleFunc("/entrypoint_listeners/new", restrictedHandler(c, NewEntrypointListener)).Methods("GET")
	r.HandleFunc("/entrypoint_listeners", restrictedHandler(c, CreateEntrypointListener)).Methods("POST")
	r.HandleFunc("/entrypoint_listeners", restrictedHandler(c, ListEntrypointListeners)).Methods("GET")
	r.HandleFunc("/entrypoint_listeners/{id}", restrictedHandler(c, GetEntrypointListener)).Methods("GET")

	return baseRouter
}

func restrictedHandler(c *core.Core, fn func(*client.Client, http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Load Client by Session ID stored in cookie. If either cookie or client
		// does not exist, redirect to login page with 401.
		var sessionID string
		var client *client.Client

		if sessionCookie, err := r.Cookie(core.SessionCookieName); err == nil {
			sessionID = sessionCookie.Value
			client = c.Sessions.Client(sessionID)
		}

		if client == nil {
			http.Redirect(w, r, "/ui/sessions/new", http.StatusFound) // can't do 401 here unless you want browser behavior
			return
		}
		if err := fn(client, w, r); err != nil {

			status := http.StatusInternalServerError
			if strings.Contains(err.Error(), "404") {
				status = http.StatusNotFound
			}

			w.WriteHeader(status)
			w.Write([]byte(err.Error()))
		}
	}
}

func openHandler(c *core.Core, fn func(*client.Client, http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Unauthenticated client
		if err := fn(c.APIClient("", ""), w, r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}
}

//------------------------------------------------------------------------------

func parseID(r *http.Request) (*int64, error) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return nil, err
	}
	id64 := int64(id)
	return &id64, nil
}

func renderTemplate(sg *client.Client, w http.ResponseWriter, name string, data map[string]interface{}) error {

	// TODO
	if mi := data["model"]; mi != nil {
		if m, isModel := mi.(model.Model); isModel {
			model.ZeroPrivateFields(m)
			data["model"] = m
		}
	}

	modelJSON, _ := json.Marshal(data["model"])
	data["modelJSON"] = string(modelJSON)

	if fields, ok := data["fields"]; ok {
		fieldsJSON, _ := json.Marshal(fields)
		data["fieldsJSON"] = string(fieldsJSON)
	}

	// Version number
	data["supergiantVersion"] = sg.Version

	if err := templates[name].ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	return nil
}

func unmarshalFormInto(r *http.Request, out interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	return json.Unmarshal([]byte(r.PostForm.Get("json_input")), out)
}

//------------------------------------------------------------------------------

func uiRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/ui", http.StatusFound)
}

//------------------------------------------------------------------------------

func Root(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	http.Redirect(w, r, "/ui/cloud_accounts", http.StatusFound)
	return nil
}
