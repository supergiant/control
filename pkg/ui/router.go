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
	"github.com/supergiant/supergiant/pkg/client"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

var templates = make(map[string]*template.Template)

func init() {
	// AssetNames() is like a filepath Glob of our generated assets
	viewFiles := AssetNames()
	var partials []string
	var layouts []string
	for _, viewFile := range viewFiles {
		if !strings.HasPrefix(viewFile, "ui/views") {
			continue // don't want to load any non-HTML assets into this
		}
		if strings.HasPrefix(viewFile, "ui/views/partials") {
			partials = append(partials, viewFile)
		} else {
			layouts = append(layouts, viewFile)
		}
	}

	for _, partial := range partials {
		// https://golang.org/src/html/template/template.go?s=11938:12007#L354
		var t *template.Template
		for _, view := range append(layouts, partial) {
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

			src, err := Asset(view)
			if err != nil {
				panic(err)
			}

			if _, err = tmpl.Parse(string(src)); err != nil {
				panic(err)
			}
		}

		key := regexp.MustCompile(`([^/]+)\.html$`).FindStringSubmatch(partial)[1]
		templates[key] = t
	}
}

func NewRouter(c *core.Core, baseRouter *mux.Router) *mux.Router {
	base := baseRouter.StrictSlash(true)

	// Redirect / to /ui
	base.HandleFunc("/", uiRedirect).Methods("GET")

	r := base.PathPrefix("/ui").Subrouter()

	// Assets
	assetDir := &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo}
	r.PathPrefix("/assets/").Handler(http.FileServer(assetDir))

	r.HandleFunc("/", restrictedHandler(c, Root)).Methods("GET")

	// Login-based functions can't be authenticated
	sharedClient := c.NewAPIClient("", "")
	r.HandleFunc("/sessions/new", openHandler(sharedClient, NewSession)).Methods("GET")
	r.HandleFunc("/sessions", openHandler(sharedClient, CreateSession)).Methods("POST")
	r.HandleFunc("/sessions/{id}", openHandler(sharedClient, GetSession)).Methods("GET")

	r.HandleFunc("/sessions", restrictedHandler(c, ListSessions)).Methods("GET")
	r.HandleFunc("/sessions/{id}/delete", restrictedHandler(c, DeleteSession)).Methods("PUT")

	r.HandleFunc("/users/new", restrictedHandler(c, NewUser)).Methods("GET")
	r.HandleFunc("/users", restrictedHandler(c, CreateUser)).Methods("POST")
	r.HandleFunc("/users", restrictedHandler(c, ListUsers)).Methods("GET")
	r.HandleFunc("/users/{id}", restrictedHandler(c, GetUser)).Methods("GET")
	r.HandleFunc("/users/{id}/edit", restrictedHandler(c, EditUser)).Methods("GET")
	r.HandleFunc("/users/{id}", restrictedHandler(c, UpdateUser)).Methods("PUT")
	r.HandleFunc("/users/{id}/delete", restrictedHandler(c, DeleteUser)).Methods("PUT")
	r.HandleFunc("/users/{id}/regenerate_api_token", restrictedHandler(c, RegenerateUserAPIToken)).Methods("PUT")

	r.HandleFunc("/cloud_accounts/new", restrictedHandler(c, NewCloudAccount)).Methods("GET")
	r.HandleFunc("/cloud_accounts", restrictedHandler(c, CreateCloudAccount)).Methods("POST")
	r.HandleFunc("/cloud_accounts", restrictedHandler(c, ListCloudAccounts)).Methods("GET")
	r.HandleFunc("/cloud_accounts/{id}", restrictedHandler(c, GetCloudAccount)).Methods("GET")
	r.HandleFunc("/cloud_accounts/{id}/delete", restrictedHandler(c, DeleteCloudAccount)).Methods("PUT")

	r.HandleFunc("/kubes/new", restrictedHandler(c, NewKube)).Methods("GET")
	r.HandleFunc("/kubes", restrictedHandler(c, CreateKube)).Methods("POST")
	r.HandleFunc("/kubes", restrictedHandler(c, ListKubes)).Methods("GET")
	r.HandleFunc("/kubes/{id}", restrictedHandler(c, GetKube)).Methods("GET")
	r.HandleFunc("/kubes/{id}/delete", restrictedHandler(c, DeleteKube)).Methods("PUT")

	r.HandleFunc("/apps/new", restrictedHandler(c, NewApp)).Methods("GET")
	r.HandleFunc("/apps", restrictedHandler(c, CreateApp)).Methods("POST")
	r.HandleFunc("/apps", restrictedHandler(c, ListApps)).Methods("GET")
	r.HandleFunc("/apps/{id}", restrictedHandler(c, GetApp)).Methods("GET")
	r.HandleFunc("/apps/{id}/delete", restrictedHandler(c, DeleteApp)).Methods("PUT")

	r.HandleFunc("/components/new", restrictedHandler(c, NewComponent)).Methods("GET")
	r.HandleFunc("/components", restrictedHandler(c, CreateComponent)).Methods("POST")
	r.HandleFunc("/components", restrictedHandler(c, ListComponents)).Methods("GET")
	r.HandleFunc("/components/{id}", restrictedHandler(c, GetComponent)).Methods("GET")
	r.HandleFunc("/components/{id}/delete", restrictedHandler(c, DeleteComponent)).Methods("PUT")
	r.HandleFunc("/components/{id}/deploy", restrictedHandler(c, DeployComponent)).Methods("PUT")
	r.HandleFunc("/components/{id}/configure", restrictedHandler(c, ConfigureComponent)).Methods("GET")

	r.HandleFunc("/releases", restrictedHandler(c, CreateRelease)).Methods("POST")
	r.HandleFunc("/releases/{id}", restrictedHandler(c, UpdateRelease)).Methods("PUT")

	r.HandleFunc("/instances", restrictedHandler(c, ListInstances)).Methods("GET")
	r.HandleFunc("/instances/{id}", restrictedHandler(c, GetInstance)).Methods("GET")

	r.HandleFunc("/volumes", restrictedHandler(c, ListVolumes)).Methods("GET")
	r.HandleFunc("/volumes/{id}", restrictedHandler(c, GetVolume)).Methods("GET")

	r.HandleFunc("/private_image_keys/new", restrictedHandler(c, NewPrivateImageKey)).Methods("GET")
	r.HandleFunc("/private_image_keys", restrictedHandler(c, CreatePrivateImageKey)).Methods("POST")
	r.HandleFunc("/private_image_keys", restrictedHandler(c, ListPrivateImageKeys)).Methods("GET")
	r.HandleFunc("/private_image_keys/{id}", restrictedHandler(c, GetPrivateImageKey)).Methods("GET")
	r.HandleFunc("/private_image_keys/{id}/delete", restrictedHandler(c, DeletePrivateImageKey)).Methods("PUT")

	r.HandleFunc("/entrypoints/new", restrictedHandler(c, NewEntrypoint)).Methods("GET")
	r.HandleFunc("/entrypoints", restrictedHandler(c, CreateEntrypoint)).Methods("POST")
	r.HandleFunc("/entrypoints", restrictedHandler(c, ListEntrypoints)).Methods("GET")
	r.HandleFunc("/entrypoints/{id}", restrictedHandler(c, GetEntrypoint)).Methods("GET")
	r.HandleFunc("/entrypoints/{id}/delete", restrictedHandler(c, DeleteEntrypoint)).Methods("PUT")

	r.HandleFunc("/nodes/new", restrictedHandler(c, NewNode)).Methods("GET")
	r.HandleFunc("/nodes", restrictedHandler(c, CreateNode)).Methods("POST")
	r.HandleFunc("/nodes", restrictedHandler(c, ListNodes)).Methods("GET")
	r.HandleFunc("/nodes/{id}", restrictedHandler(c, GetNode)).Methods("GET")
	r.HandleFunc("/nodes/{id}/delete", restrictedHandler(c, DeleteNode)).Methods("PUT")

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
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}
}

func openHandler(sharedClient *client.Client, fn func(*client.Client, http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(sharedClient, w, r); err != nil {
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

func renderTemplate(w http.ResponseWriter, name string, data map[string]interface{}) error {

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
	var cloudAccounts []*model.CloudAccount
	if err := sg.CloudAccounts.List(&cloudAccounts); err != nil {
		return err
	}
	if len(cloudAccounts) == 0 {
		http.Redirect(w, r, "/ui/cloud_accounts/new", http.StatusFound)
		return nil
	}

	var kubes []*model.Kube
	if err := sg.Kubes.List(&kubes); err != nil {
		return err
	}
	if len(kubes) == 0 {
		http.Redirect(w, r, "/ui/kubes/new", http.StatusFound)
		return nil
	}

	http.Redirect(w, r, "/ui/apps", http.StatusFound)
	return nil
}
