package ui

import (
	"fmt"
	"net/http"

	"github.com/supergiant/supergiant/pkg/client"
	"github.com/supergiant/supergiant/pkg/model"
)

func NewVolume(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	return renderTemplate(w, "new", map[string]interface{}{
		"title":      "Volumes",
		"formAction": "/ui/volumes",
		"model": map[string]interface{}{
			"kube_name": "",
			"name":      "",
			"type":      "gp2",
			"size":      10,
		},
	})
}

func CreateVolume(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	m := new(model.Volume)
	err := unmarshalFormInto(r, m)
	if err == nil {
		err = sg.Volumes.Create(m)
	}
	if err != nil {
		return renderTemplate(w, "new", map[string]interface{}{
			"title":      "Volumes",
			"formAction": "/ui/volumes",
			"model":      m,
			"error":      err.Error(),
		})
	}

	http.Redirect(w, r, "/ui/volumes", http.StatusFound)
	return nil
}

func ListVolumes(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	fields := []map[string]interface{}{
		{
			"title": "Kube ID",
			"type":  "field_value",
			"field": "kube_name",
		},
		{
			"title": "Name",
			"type":  "field_value",
			"field": "name",
		},
		{
			"title": "Type",
			"type":  "field_value",
			"field": "type",
		},
		{
			"title": "Size",
			"type":  "field_value",
			"field": "size",
		},
	}
	return renderTemplate(w, "index", map[string]interface{}{
		"title":       "Volumes",
		"uiBasePath":  "/ui/volumes",
		"apiBasePath": "/api/v0/volumes",
		"fields":      fields,
		"showNewLink": true,
		"actionPaths": map[string]string{
			"Edit": "/edit",
		},
		"batchActionPaths": map[string]map[string]string{
			"Delete": map[string]string{
				"method":       "DELETE",
				"relativePath": "",
			},
		},
	})
}

func GetVolume(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return err
	}
	item := new(model.Volume)
	if err := sg.Volumes.Get(id, item); err != nil {
		return err
	}
	return renderTemplate(w, "show", map[string]interface{}{
		"title": "Volumes",
		"model": item,
	})
}

func EditVolume(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return err
	}
	item := new(model.Volume)
	if err := sg.Volumes.Get(id, item); err != nil {
		return err
	}
	return renderTemplate(w, "new", map[string]interface{}{
		"title":      "Volumes",
		"formAction": fmt.Sprintf("/ui/volumes/%d", *id),
		"model": map[string]interface{}{
			"size": item.Size,
		},
	})
}

func UpdateVolume(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return err
	}
	m := new(model.Volume)
	err = unmarshalFormInto(r, m)
	if err == nil {
		err = sg.Volumes.Update(id, m)
	}
	if err != nil {
		return renderTemplate(w, "new", map[string]interface{}{
			"title":      "Volumes",
			"formAction": fmt.Sprintf("/ui/volumes/%d", *id),
			"model":      m,
			"error":      err.Error(),
		})
	}
	http.Redirect(w, r, "/ui/volumes", http.StatusFound)
	return nil
}
