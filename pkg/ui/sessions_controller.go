package ui

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/supergiant/supergiant/pkg/client"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

func NewSession(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	return renderTemplate(w, "login", map[string]interface{}{
		"title":      "Sessions",
		"formAction": "/ui/sessions",
	})
}

func CreateSession(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	m := &model.Session{
		User: &model.User{
			Username: r.PostForm.Get("username"),
			Password: r.PostForm.Get("password"),
		},
	}
	if err := sg.Sessions.Create(m); err != nil {
		return renderTemplate(w, "login", map[string]interface{}{
			"title":      "Sessions",
			"formAction": "/ui/sessions",
			"error":      err.Error(),
		})
	}

	// Store Session ID in Cookie
	cookie := &http.Cookie{
		Name:  core.SessionCookieName,
		Value: m.ID,
		Path:  "/",
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/ui/sessions", http.StatusFound)
	return nil
}

func ListSessions(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	fields := []map[string]interface{}{
		{
			"title": "User ID",
			"type":  "field_value",
			"field": "user_id",
		},
		{
			"title": "Created at",
			"type":  "field_value",
			"field": "created_at",
		},
	}
	return renderTemplate(w, "index", map[string]interface{}{
		"title":       "Sessions",
		"uiBasePath":  "/ui/sessions",
		"apiListPath": "/api/v0/sessions",
		"fields":      fields,
		"showNewLink": true,
		"batchActionPaths": map[string]string{
			"Delete": "/delete",
		},
	})
}

func GetSession(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	item := new(model.Session)
	if err := sg.Sessions.Get(id, item); err != nil {
		return err
	}
	return renderTemplate(w, "show", map[string]interface{}{
		"title": "Sessions",
		"model": item,
	})
}

func DeleteSession(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	item := new(model.Session)
	item.ID = id
	if err := sg.Sessions.Delete(id, item); err != nil {
		return err
	}
	return nil
}
