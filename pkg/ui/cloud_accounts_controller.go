package ui

import (
	"net/http"

	"github.com/supergiant/supergiant/pkg/client"
	"github.com/supergiant/supergiant/pkg/model"
)

func NewCloudAccount(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	var m map[string]interface{}
	switch r.URL.Query().Get("option") {
	// case "aws":
	case "digitalocean":
		m = map[string]interface{}{
			"name":     "",
			"provider": "digitalocean",
			"credentials": map[string]interface{}{
				"token": "",
			},
		}
	default: // just default to AWS if option not provided, or mismatched
		m = map[string]interface{}{
			"name":     "",
			"provider": "aws",
			"credentials": map[string]interface{}{
				"access_key": "",
				"secret_key": "",
			},
		}
	}

	return renderTemplate(w, "new", map[string]interface{}{
		"title":      "Cloud Accounts",
		"formAction": "/ui/cloud_accounts",
		"model":      m,
	})
}

func CreateCloudAccount(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	m := new(model.CloudAccount)
	if err := unmarshalFormInto(r, m); err != nil {
		return err
	}
	if err := sg.CloudAccounts.Create(m); err != nil {
		return renderTemplate(w, "new", map[string]interface{}{
			"title":      "Cloud Accounts",
			"formAction": "/ui/cloud_accounts",
			"model":      m,
			"error":      err.Error(),
		})
	}
	http.Redirect(w, r, "/ui/cloud_accounts", http.StatusFound)
	return nil
}

func ListCloudAccounts(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	fields := []map[string]interface{}{
		{
			"title": "Name",
			"type":  "field_value",
			"field": "name",
		},
		{
			"title": "Provider",
			"type":  "field_value",
			"field": "provider",
		},
	}
	return renderTemplate(w, "index", map[string]interface{}{
		"title":       "Cloud Accounts",
		"uiBasePath":  "/ui/cloud_accounts",
		"apiListPath": "/api/v0/cloud_accounts",
		"fields":      fields,
		"showNewLink": true,
		"newOptions": map[string]string{
			"aws":          "AWS",
			"digitalocean": "DigitalOcean",
		},
		"batchActionPaths": map[string]string{
			"Delete": "/delete",
		},
	})
}

func GetCloudAccount(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return err
	}
	item := new(model.CloudAccount)
	if err := sg.CloudAccounts.Get(id, item); err != nil {
		return err
	}
	return renderTemplate(w, "show", map[string]interface{}{
		"title": "Cloud Accounts",
		"model": item,
	})
}

func DeleteCloudAccount(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return err
	}
	item := new(model.CloudAccount)
	if err := sg.CloudAccounts.Delete(id, item); err != nil {
		return err
	}
	// http.Redirect(w, r, "/ui/cloud_accounts", http.StatusFound)
	return nil
}
