package ui

import (
	"net/http"
	"strconv"

	"github.com/supergiant/supergiant/pkg/client"
	"github.com/supergiant/supergiant/pkg/model"
)

func NewHelmRelease(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	chart := new(model.HelmChart)

	if chartIDStr := r.URL.Query().Get("chart_id"); chartIDStr != "" {
		id, err := strconv.Atoi(chartIDStr)
		if err != nil {
			return err
		}
		id64 := int64(id)

		if err = sg.HelmCharts.Get(id64, chart); err != nil {
			return err
		}
	}

	return renderTemplate(sg, w, "new", map[string]interface{}{
		"title":      "HelmReleases",
		"formAction": "/ui/helm_releases",
		"model": map[string]interface{}{
			"kube_name":     "",
			"repo_name":     chart.RepoName,
			"chart_name":    chart.Name,
			"chart_version": chart.Version,
			"config":        chart.DefaultConfig,
			"name":          "",
		},
	})
}

func CreateHelmRelease(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	m := new(model.HelmRelease)
	err := unmarshalFormInto(r, m)
	if err == nil {
		err = sg.HelmReleases.Create(m)
	}
	if err != nil {
		return renderTemplate(sg, w, "new", map[string]interface{}{
			"title":      "Helm Releases",
			"formAction": "/ui/helm_releases",
			"model":      m,
			"error":      err.Error(),
		})
	}
	http.Redirect(w, r, "/ui/helm_releases", http.StatusFound)
	return nil
}

func ListHelmReleases(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
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
			"title": "Chart name",
			"type":  "field_value",
			"field": "chart_name",
		},
		{
			"title": "Chart version",
			"type":  "field_value",
			"field": "chart_version",
		},
		{
			"title": "Revision",
			"type":  "field_value",
			"field": "revision",
		},
		{
			"title": "Updated",
			"type":  "field_value",
			"field": "updated_value",
		},
		// {
		// 	"title": "Status",
		// 	"type":  "field_value",
		// 	"field": "status_value",
		// },
	}
	return renderTemplate(sg, w, "index", map[string]interface{}{
		"title":         "Helm Releases",
		"uiBasePath":    "/ui/helm_releases",
		"apiBasePath":   "/api/v0/helm_releases",
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

func GetHelmRelease(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return err
	}
	item := new(model.HelmRelease)
	if err := sg.HelmReleases.Get(id, item); err != nil {
		return err
	}
	return renderTemplate(sg, w, "show", map[string]interface{}{
		"title": "Helm Releases",
		"model": item,
	})
}
