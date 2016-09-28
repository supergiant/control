package ui

import (
	"net/http"

	"github.com/supergiant/supergiant/pkg/client"
	"github.com/supergiant/supergiant/pkg/model"
)

func NewEntrypointListener(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	return renderTemplate(w, "new", map[string]interface{}{
		"title":      "Entrypoint Listeners",
		"formAction": "/ui/entrypoint_listeners",
		"model": map[string]interface{}{
			"entrypoint_name":     "",
			"entrypoint_port":     0,
			"entrypoint_protocol": "TCP",
			"node_port":           0,
			"node_protocol":       "TCP",
		},
	})
}

func CreateEntrypointListener(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	m := new(model.EntrypointListener)
	if err := unmarshalFormInto(r, m); err != nil {
		return err
	}
	if err := sg.EntrypointListeners.Create(m); err != nil {
		return renderTemplate(w, "new", map[string]interface{}{
			"title":      "Entrypoint Listeners",
			"formAction": "/ui/entrypoint_listeners",
			"model":      m,
			"error":      err.Error(),
		})
	}
	http.Redirect(w, r, "/ui/entrypoint_listeners", http.StatusFound)
	return nil
}

func ListEntrypointListeners(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	fields := []map[string]interface{}{
		{
			"title": "Entrypoint ID",
			"type":  "field_value",
			"field": "entrypoint_name",
		},
		{
			"title": "Entrypoint Port",
			"type":  "field_value",
			"field": "entrypoint_port",
		},
		{
			"title": "Entrypoint Protocol",
			"type":  "field_value",
			"field": "entrypoint_protocol",
		},
		{
			"title": "Node Port",
			"type":  "field_value",
			"field": "node_port",
		},
		{
			"title": "Node Protocol",
			"type":  "field_value",
			"field": "node_protocol",
		},
	}
	return renderTemplate(w, "index", map[string]interface{}{
		"title":       "Entrypoint Listeners",
		"uiBasePath":  "/ui/entrypoint_listeners",
		"apiListPath": "/api/v0/entrypoint_listeners",
		"fields":      fields,
		"showNewLink": true,
		"batchActionPaths": map[string]string{
			"Delete": "/delete",
		},
	})
}

func GetEntrypointListener(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return err
	}
	item := new(model.EntrypointListener)
	if err := sg.EntrypointListeners.Get(id, item); err != nil {
		return err
	}
	return renderTemplate(w, "show", map[string]interface{}{
		"title": "Entrypoint Listeners",
		"model": item,
	})
}

func DeleteEntrypointListener(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return err
	}
	item := new(model.EntrypointListener)
	if err := sg.EntrypointListeners.Delete(id, item); err != nil {
		return err
	}
	// http.Redirect(w, r, "/ui/entrypoint_listeners", http.StatusFound)
	return nil
}
