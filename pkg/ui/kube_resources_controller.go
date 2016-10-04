package ui

import (
	"fmt"
	"net/http"

	"github.com/supergiant/supergiant/pkg/client"
	"github.com/supergiant/supergiant/pkg/model"
)

func NewKubeResource(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	var m map[string]interface{}
	switch r.URL.Query().Get("option") {
	case "pod":
		m = map[string]interface{}{
			"kube_name": "",
			"kind":      "Pod",
			"namespace": "default",
			"name":      "",
			"template": map[string]interface{}{
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
	case "service":
		m = map[string]interface{}{
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
							"SUPERGIANT_ENTRYPOINT_LISTENER": map[string]interface{}{
								"entrypoint_name": "",
								"entrypoint_port": 80,
							},
						},
					},
				},
			},
		}
	default:
		m = map[string]interface{}{
			"kube_name": "",
			"kind":      "",
			"namespace": "default",
			"name":      "",
			"template":  map[string]interface{}{},
		}
	}

	return renderTemplate(sg, w, "new", map[string]interface{}{
		"title":      "Kube Resources",
		"formAction": "/ui/kube_resources",
		"model":      m,
	})
}

func CreateKubeResource(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	m := new(model.KubeResource)
	err := unmarshalFormInto(r, m)
	if err == nil {
		err = sg.KubeResources.Create(m)
	}
	if err != nil {
		return renderTemplate(sg, w, "new", map[string]interface{}{
			"title":      "Kube Resources",
			"formAction": "/ui/kube_resources",
			"model":      m,
			"error":      err.Error(),
		})
	}

	http.Redirect(w, r, "/ui/kube_resources", http.StatusFound)
	return nil
}

func ListKubeResources(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	fields := []map[string]interface{}{
		{
			"title": "Kube name",
			"type":  "field_value",
			"field": "kube_name",
		},
		{
			"title": "Kind",
			"type":  "field_value",
			"field": "kind",
		},
		{
			"title": "Namespace",
			"type":  "field_value",
			"field": "namespace",
		},
		{
			"title": "Name",
			"type":  "field_value",
			"field": "name",
		},
	}
	return renderTemplate(sg, w, "kube_resources", map[string]interface{}{
		"title":       "Kube Resources",
		"uiBasePath":  "/ui/kube_resources",
		"apiBasePath": "/api/v0/kube_resources",
		"fields":      fields,
		"showNewLink": true,
		"newOptions": map[string]string{
			"pod":     "Pod",
			"service": "Service",
			"other":   "Other",
		},
		"actionPaths": map[string]string{
			"Edit": "/edit",
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

func GetKubeResource(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return err
	}
	item := new(model.KubeResource)
	if err := sg.KubeResources.Get(id, item); err != nil {
		return err
	}
	return renderTemplate(sg, w, "show", map[string]interface{}{
		"title": "Kube Resources",
		"model": item,
	})
}

func EditKubeResource(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return err
	}
	item := new(model.KubeResource)
	if err := sg.KubeResources.Get(id, item); err != nil {
		return err
	}
	return renderTemplate(sg, w, "new", map[string]interface{}{
		"title":      "Kube Resources",
		"formAction": fmt.Sprintf("/ui/kube_resources/%d", *id),
		"model": map[string]interface{}{
			"template": item.Template,
		},
	})
}

func UpdateKubeResource(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return err
	}
	m := new(model.KubeResource)
	err = unmarshalFormInto(r, m)
	if err == nil {
		err = sg.KubeResources.Update(id, m)
	}
	if err != nil {
		return renderTemplate(sg, w, "new", map[string]interface{}{
			"title":      "Kube Resources",
			"formAction": fmt.Sprintf("/ui/kube_resources/%d", *id),
			"model":      m,
			"error":      err.Error(),
		})
	}
	http.Redirect(w, r, "/ui/kube_resources", http.StatusFound)
	return nil
}
