package ui

import (
	"net/http"

	"github.com/supergiant/supergiant/pkg/client"
)

func NewService(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	m := map[string]interface{}{
		"kube_name": "",
		"kind":      "Service",
		"namespace": "default",
		"name":      "",
		"template": map[string]interface{}{
			"spec": map[string]interface{}{
				"type":     "NodePort",
				"selector": map[string]string{},
				"ports": []map[string]interface{}{
					{
						"name": "jenkins",
						"port": 8080,
					},
				},
			},
		},
	}

	return renderTemplate(sg, w, "new", map[string]interface{}{
		"title":      "Services",
		"formAction": "/ui/services",
		"model":      m,
	})
}

func ListServices(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
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
			"title": "Namespace",
			"type":  "field_value",
			"field": "namespace",
		},
	}
	return renderTemplate(sg, w, "kube_resources", map[string]interface{}{
		"title":         "Services",
		"uiBasePath":    "/ui/services",
		"apiBasePath":   "/api/v0/kube_resources?filter.kind=Service",
		"fields":        fields,
		"showNewLink":   true,
		"showStatusCol": true,
		"actionPaths": map[string]string{
			"Edit": "/ui/services/{{ ID }}/edit",
		},
		"batchActionPaths": map[string]map[string]string{
			"Delete": map[string]string{
				"method":       "DELETE",
				"relativePath": "",
			},
			"Start": map[string]string{
				"method":       "POST",
				"relativePath": "/start",
			},
			"Stop": map[string]string{
				"method":       "POST",
				"relativePath": "/stop",
			},
		},
	})
}
