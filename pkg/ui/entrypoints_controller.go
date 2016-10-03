package ui

import (
	"net/http"

	"github.com/supergiant/supergiant/pkg/client"
	"github.com/supergiant/supergiant/pkg/model"
)

func NewEntrypoint(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	return renderTemplate(w, "new", map[string]interface{}{
		"title":      "Entrypoints",
		"formAction": "/ui/entrypoints",
		"model": map[string]interface{}{
			"kube_name": "",
			"name":      "",
		},
	})
}

func CreateEntrypoint(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	m := new(model.Entrypoint)
	err := unmarshalFormInto(r, m)
	if err == nil {
		err = sg.Entrypoints.Create(m)
	}
	if err != nil {
		return renderTemplate(w, "new", map[string]interface{}{
			"title":      "Entrypoints",
			"formAction": "/ui/entrypoints",
			"model":      m,
			"error":      err.Error(),
		})
	}
	http.Redirect(w, r, "/ui/entrypoints", http.StatusFound)
	return nil
}

func ListEntrypoints(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	fields := []map[string]interface{}{
		{
			"title": "Kube ID",
			"type":  "field_value",
			"field": "kube_name",
		},
		{
			"title": "Name",
			"type":  "field_value",
			"field": "name",
		},
		{
			"title": "Address",
			"type":  "field_value",
			"field": "address",
		},
	}
	return renderTemplate(w, "index", map[string]interface{}{
		"title":       "Entrypoints",
		"uiBasePath":  "/ui/entrypoints",
		"apiBasePath": "/api/v0/entrypoints",
		"fields":      fields,
		"showNewLink": true,
		"batchActionPaths": map[string]map[string]string{
			"Delete": map[string]string{
				"method":       "DELETE",
				"relativePath": "",
			},
		},
	})
}

func GetEntrypoint(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return err
	}
	item := new(model.Entrypoint)
	if err := sg.Entrypoints.Get(id, item); err != nil {
		return err
	}
	return renderTemplate(w, "show", map[string]interface{}{
		"title": "Entrypoints",
		"model": item,
	})
}
