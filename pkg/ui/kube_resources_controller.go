package ui

import (
	"fmt"
	"net/http"

	"github.com/supergiant/supergiant/pkg/client"
	"github.com/supergiant/supergiant/pkg/model"
)

func NewKubeResource(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	m := map[string]interface{}{
		"kube_name": "",
		"kind":      "",
		"namespace": "default",
		"name":      "",
		"template":  map[string]interface{}{},
	}

	return renderTemplate(sg, w, "new", map[string]interface{}{
		"title":      "Kube Resources",
		"formAction": "/ui/kube_resources",
		"model":      m,
	})
}

func ListKubeResources(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
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
			"title": "Kind",
			"type":  "field_value",
			"field": "kind",
		},
		{
			"title": "Namespace",
			"type":  "field_value",
			"field": "namespace",
		},
	}
	return renderTemplate(sg, w, "kube_resources", map[string]interface{}{
		"title":         "Kube Resources",
		"uiBasePath":    "/ui/kube_resources",
		"apiBasePath":   "/api/v0/kube_resources",
		"fields":        fields,
		"showNewLink":   true,
		"showStatusCol": true,
		"actionPaths": map[string]string{
			"Edit": "/ui/kube_resources/{{ ID }}/edit",
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
			"resource": item.Resource,
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
