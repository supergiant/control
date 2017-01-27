package ui

import (
	"net/http"

	"github.com/supergiant/supergiant/pkg/client"
	"github.com/supergiant/supergiant/pkg/model"
)

func NewNode(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	return renderTemplate(sg, w, "new", map[string]interface{}{
		"title":      "Nodes",
		"formAction": "/ui/nodes",
		"model": map[string]interface{}{
			"kube_name": "",
			"size":      "",
		},
	})
}

func CreateNode(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	m := new(model.Node)
	err := unmarshalFormInto(r, m)
	if err == nil {
		err = sg.Nodes.Create(m)
	}
	if err != nil {
		return renderTemplate(sg, w, "new", map[string]interface{}{
			"title":      "Nodes",
			"formAction": "/ui/nodes",
			"model":      m,
			"error":      err.Error(),
		})
	}

	http.Redirect(w, r, "/ui/nodes", http.StatusFound)
	return nil
}

func ListNodes(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	fields := []map[string]interface{}{
		{
			"title": "Kube name",
			"type":  "field_value",
			"field": "kube_name",
		},
		{
			"title": "Size",
			"type":  "field_value",
			"field": "size",
		},
		{
			"title": "Provider ID",
			"type":  "field_value",
			"field": "provider_id",
		},
		{
			"title":             "CPU Usage",
			"type":              "percentage",
			"numerator_field":   "cpu_usage",
			"denominator_field": "cpu_limit",
		},
		{
			"title":             "RAM usage",
			"type":              "percentage",
			"numerator_field":   "ram_usage",
			"denominator_field": "ram_limit",
		},
	}
	return renderTemplate(sg, w, "index", map[string]interface{}{
		"title":         "Nodes",
		"uiBasePath":    "/ui/nodes",
		"apiBasePath":   "/api/v0/nodes",
		"fields":        fields,
		"showNewLink":   true,
		"showStatusCol": true,
		"batchActionPaths": map[string]map[string]string{
			"Delete": map[string]string{
				"method":       "DELETE",
				"relativePath": "",
			},
		},
	})
}

func GetNode(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return err
	}
	item := new(model.Node)
	if err := sg.Nodes.Get(id, item); err != nil {
		return err
	}
	return renderTemplate(sg, w, "show", map[string]interface{}{
		"title": "Nodes",
		"model": item,
	})
}
