package ui

import (
	"net/http"

	"github.com/supergiant/supergiant/pkg/client"
)

func ListVolumes(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	fields := []map[string]interface{}{
		{
			"title": "Name",
			"type":  "field_value",
			"field": "name",
		},
		{
			"title": "Kube name",
			"type":  "field_value",
			"field": "kube_name",
		},
		{
			"title": "Capacity",
			"type":  "field_value",
			"field": "resource.spec.capacity.storage",
		},
	}
	return renderTemplate(sg, w, "kube_resources", map[string]interface{}{
		"title":       "Volumes",
		"uiBasePath":  "/ui/volumes",
		"apiBasePath": "/api/v0/kube_resources?filter.kind=PersistentVolume",
		"fields":      fields,
		// "showNewLink":   true,
		"showStatusCol": true,
		"actionPaths": map[string]string{
			"Edit": "/ui/volumes/{{ ID }}/edit",
		},
		"batchActionPaths": map[string]map[string]string{
			"Delete": map[string]string{
				"method":       "DELETE",
				"relativePath": "",
			},
		},
	})
}
