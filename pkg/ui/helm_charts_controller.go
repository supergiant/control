package ui

import (
	"net/http"

	"github.com/supergiant/supergiant/pkg/client"
	"github.com/supergiant/supergiant/pkg/model"
)

func NewHelmChart(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	return renderTemplate(sg, w, "new", map[string]interface{}{
		"title":      "HelmCharts",
		"formAction": "/ui/helm_charts",
		"model": map[string]interface{}{
			"name": "",
			"url":  "",
		},
	})
}

func CreateHelmChart(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	m := new(model.HelmChart)
	err := unmarshalFormInto(r, m)
	if err == nil {
		err = sg.HelmCharts.Create(m)
	}
	if err != nil {
		return renderTemplate(sg, w, "new", map[string]interface{}{
			"title":      "Helm Charts",
			"formAction": "/ui/helm_charts",
			"model":      m,
			"error":      err.Error(),
		})
	}
	http.Redirect(w, r, "/ui/helm_charts", http.StatusFound)
	return nil
}

func ListHelmCharts(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	fields := []map[string]interface{}{
		{
			"title": "Name",
			"type":  "field_value",
			"field": "name",
		},
		{
			"title": "Repo name",
			"type":  "field_value",
			"field": "repo_name",
		},
		{
			"title": "Version",
			"type":  "field_value",
			"field": "version",
		},
		{
			"title": "Description",
			"type":  "field_value",
			"field": "description",
		},
	}
	return renderTemplate(sg, w, "index", map[string]interface{}{
		"title":       "Helm Charts",
		"uiBasePath":  "/ui/helm_charts",
		"apiBasePath": "/api/v0/helm_charts",
		"fields":      fields,
		"showNewLink": false,
		"actionPaths": map[string]string{
			"Launch": "/ui/helm_releases/new?chart_id={{ ID }}",
		},
		"batchActionPaths": map[string]map[string]string{
			"Delete": {
				"method":       "DELETE",
				"relativePath": "",
			},
		},
	})
}

func GetHelmChart(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return err
	}
	item := new(model.HelmChart)
	if err := sg.HelmCharts.Get(id, item); err != nil {
		return err
	}
	return renderTemplate(sg, w, "show", map[string]interface{}{
		"title": "Helm Charts",
		"model": item,
	})
}
