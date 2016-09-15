package ui

import (
	"fmt"
	"net/http"

	"github.com/supergiant/supergiant/pkg/client"
	"github.com/supergiant/supergiant/pkg/model"
)

func NewUser(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	return renderTemplate(w, "new", map[string]interface{}{
		"title":      "Users",
		"formAction": "/ui/users",
		"formMethod": "POST",
		"model": map[string]interface{}{
			"username": "",
			"password": "",
			"role":     "user",
		},
	})
}

func CreateUser(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	m := new(model.User)
	if err := unmarshalFormInto(r, m); err != nil {
		return err
	}
	if err := sg.Users.Create(m); err != nil {
		return renderTemplate(w, "new", map[string]interface{}{
			"title":      "Users",
			"formAction": "/ui/users",
			"formMethod": "POST",
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
	return renderTemplate(w, "index", map[string]interface{}{
		"title":       "Users",
		"uiBasePath":  "/ui/users",
		"apiListPath": "/api/v0/users",
		"fields":      fields,
		"showNewLink": true,
		"actionPaths": map[string]string{
			"Edit": "/edit",
		},
		"batchActionPaths": map[string]string{
			"Regenerate API token": "/regenerate_api_token",
			"Delete":               "/delete",
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
	return renderTemplate(w, "show", map[string]interface{}{
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
	return renderTemplate(w, "new", map[string]interface{}{
		"title":      "Users",
		"formAction": fmt.Sprintf("/ui/users/%d", *id),
		"formMethod": "PUT",
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
	if err := unmarshalFormInto(r, m); err != nil {
		return err
	}
	if err := sg.Users.Update(id, m); err != nil {
		return renderTemplate(w, "new", map[string]interface{}{
			"title":      "Users",
			"formAction": fmt.Sprintf("/ui/users/%d", *id),
			"formMethod": "PUT",
			"model":      m,
			"error":      err.Error(),
		})
	}
	http.Redirect(w, r, "/ui/users", http.StatusFound)
	return nil
}

func DeleteUser(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return err
	}
	item := new(model.User)
	if err := sg.Users.Delete(id, item); err != nil {
		return err
	}
	return nil
}

func RegenerateUserAPIToken(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return err
	}
	item := new(model.User)
	if err := sg.Users.RegenerateAPIToken(id, item); err != nil {
		return err
	}
	return nil
}
