package ui

import (
	"net/http"

	"github.com/supergiant/supergiant/pkg/client"
	"github.com/supergiant/supergiant/pkg/model"
)

func NewHelmRepo(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	return renderTemplate(sg, w, "new", map[string]interface{}{
		"title":      "HelmRepos",
		"formAction": "/ui/helm_repos",
		"model": map[string]interface{}{
			"name": "",
			"url":  "",
		},
	})
}

func CreateHelmRepo(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	m := new(model.HelmRepo)
	err := unmarshalFormInto(r, m)
	if err == nil {
		err = sg.HelmRepos.Create(m)
	}
	if err != nil {
		return renderTemplate(sg, w, "new", map[string]interface{}{
			"title":      "Helm Repos",
			"formAction": "/ui/helm_repos",
			"model":      m,
			"error":      err.Error(),
		})
	}
	http.Redirect(w, r, "/ui/helm_repos", http.StatusFound)
	return nil
}

func ListHelmRepos(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	fields := []map[string]interface{}{
		{
			"title": "Name",
			"type":  "field_value",
			"field": "name",
		},
		{
			"title": "URL",
			"type":  "field_value",
			"field": "url",
		},
	}
	return renderTemplate(sg, w, "index", map[string]interface{}{
		"title":       "Helm Repos",
		"uiBasePath":  "/ui/helm_repos",
		"apiBasePath": "/api/v0/helm_repos",
		"fields":      fields,
		"showNewLink": true,
		"batchActionPaths": map[string]map[string]string{
			"Delete": {
				"method":       "DELETE",
				"relativePath": "",
			},
		},
	})
}

func GetHelmRepo(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return err
	}
	item := new(model.HelmRepo)
	if err := sg.HelmRepos.Get(id, item); err != nil {
		return err
	}
	return renderTemplate(sg, w, "show", map[string]interface{}{
		"title": "Helm Repos",
		"model": item,
	})
}
