package ui

import (
	"net/http"

	"github.com/supergiant/supergiant/pkg/client"
)

func NewPod(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	m := map[string]interface{}{
		"kube_name": "",
		"kind":      "Pod",
		"namespace": "default",
		"name":      "",
		"resource": map[string]interface{}{
			"metadata": map[string]interface{}{
				"labels": map[string]string{},
			},
			"spec": map[string]interface{}{
				"containers": []map[string]interface{}{
					{
						"name":  "jenkins",
						"image": "jenkins",
						"ports": []map[string]interface{}{
							{
								"containerPort": 8080,
							},
						},
					},
				},
			},
		},
	}
	// case "service":
	// 	m = map[string]interface{}{
	// 		"kube_name": "",
	// 		"kind":      "Service",
	// 		"namespace": "default",
	// 		"name":      "",
	// 		"template": map[string]interface{}{
	// 			"spec": map[string]interface{}{
	// 				"type":     "NodePort",
	// 				"selector": map[string]string{},
	// 				"ports": []map[string]interface{}{
	// 					{
	// 						"name": "jenkins",
	// 						"port": 8080,
	// 					},
	// 				},
	// 			},
	// 		},
	// 	}
	// default:
	// 	m = map[string]interface{}{
	// 		"kube_name": "",
	// 		"kind":      "",
	// 		"namespace": "default",
	// 		"name":      "",
	// 		"template":  map[string]interface{}{},
	// 	}
	// }

	return renderTemplate(sg, w, "new", map[string]interface{}{
		"title":      "Pods",
		"formAction": "/ui/pods",
		"model":      m,
	})
}

func ListPods(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
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
		"title":         "Pods",
		"uiBasePath":    "/ui/pods",
		"apiBasePath":   "/api/v0/kube_resources?filter.kind=Pod",
		"fields":        fields,
		"showNewLink":   true,
		"showStatusCol": true,
		"actionPaths": map[string]string{
			"Edit": "/ui/pods/{{ ID }}/edit",
		},
		"batchActionPaths": map[string]map[string]string{
			"Delete": {
				"method":       "DELETE",
				"relativePath": "",
			},
			"Start": {
				"method":       "POST",
				"relativePath": "/start",
			},
			"Stop": {
				"method":       "POST",
				"relativePath": "/stop",
			},
		},
	})
}
