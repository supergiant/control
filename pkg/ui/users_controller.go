package ui

import (
	"fmt"
	"net/http"

	"github.com/supergiant/supergiant/pkg/client"
	"github.com/supergiant/supergiant/pkg/model"
)

func NewUser(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	return renderTemplate(sg, w, "new", map[string]interface{}{
		"title":      "Users",
		"formAction": "/ui/users",
		"model": map[string]interface{}{
			"username": "",
			"password": "",
			"role":     "user",
		},
	})
}

func CreateUser(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	m := new(model.User)
	err := unmarshalFormInto(r, m)
	if err == nil {
		err = sg.Users.Create(m)
	}
	if err != nil {
		return renderTemplate(sg, w, "new", map[string]interface{}{
			"title":      "Users",
			"formAction": "/ui/users",
			"model":      m,
			"error":      err.Error(),
		})
	}
	http.Redirect(w, r, "/ui/users", http.StatusFound)
	return nil
}

func ListUsers(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	fields := []map[string]interface{}{
		{
			"title": "Username",
			"type":  "field_value",
			"field": "username",
		},
		{
			"title": "Role",
			"type":  "field_value",
			"field": "role",
		},
	}
	return renderTemplate(sg, w, "index", map[string]interface{}{
		"title":       "Users",
		"uiBasePath":  "/ui/users",
		"apiBasePath": "/api/v0/users",
		"fields":      fields,
		"showNewLink": true,
		"actionPaths": map[string]string{
			"Edit": "/ui/users/{{ ID }}/edit",
		},
		"batchActionPaths": map[string]map[string]string{
			"Delete": {
				"method":       "DELETE",
				"relativePath": "",
			},
			"Regenerate API token": {
				"method":       "POST",
				"relativePath": "/regenerate_api_token",
			},
		},
	})
}

func GetUser(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return err
	}
	item := new(model.User)
	if err := sg.Users.Get(id, item); err != nil {
		return err
	}
	return renderTemplate(sg, w, "show", map[string]interface{}{
		"title": "Users",
		"model": item,
	})
}

func EditUser(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return err
	}
	item := new(model.User)
	if err := sg.Users.Get(id, item); err != nil {
		return err
	}
	return renderTemplate(sg, w, "new", map[string]interface{}{
		"title":      "Users",
		"formAction": fmt.Sprintf("/ui/users/%d", *id),
		"model": map[string]interface{}{
			"password": "",
			"role":     item.Role,
		},
	})
}

func UpdateUser(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return err
	}
	m := new(model.User)
	err = unmarshalFormInto(r, m)
	if err == nil {
		err = sg.Users.Update(id, m)
	}
	if err != nil {
		return renderTemplate(sg, w, "new", map[string]interface{}{
			"title":      "Users",
			"formAction": fmt.Sprintf("/ui/users/%d", *id),
			"model":      m,
			"error":      err.Error(),
		})
	}
	http.Redirect(w, r, "/ui/users", http.StatusFound)
	return nil
}
