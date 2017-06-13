package ui

import (
	"net/http"

	"github.com/supergiant/supergiant/pkg/client"
	"github.com/supergiant/supergiant/pkg/model"
)

// NewCloudAccount holds template info for UI cloud accounts.
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
	case "openstack":
		m = map[string]interface{}{
			"name":     "",
			"provider": "openstack",
			"credentials": map[string]interface{}{
				"identity_endpoint": "",
				"username":          "",
				"password":          "",
				"tenant_id":         "",
				"domain_id":         "",
				"domain_name":       "",
			},
		}
	case "gce":
		m = map[string]interface{}{
			"name":     "",
			"provider": "gce",
			"credentials": map[string]interface{}{
				"type":                        "",
				"project_id":                  "",
				"private_key_id":              "",
				"private_key":                 "",
				"client_email":                "",
				"client_id":                   "",
				"auth_uri":                    "",
				"token_uri":                   "",
				"auth_provider_x509_cert_url": "",
				"client_x509_cert_url":        "",
			},
		}
	case "packet":
		m = map[string]interface{}{
			"name":     "",
			"provider": "packet",
			"credentials": map[string]interface{}{
				"api_token": "",
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

	return renderTemplate(sg, w, "new", map[string]interface{}{
		"title":      "Cloud Accounts",
		"formAction": "/ui/cloud_accounts",
		"model":      m,
	})
}

func CreateCloudAccount(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	m := new(model.CloudAccount)
	err := unmarshalFormInto(r, m)
	if err == nil {
		err = sg.CloudAccounts.Create(m)
	}
	if err != nil {
		return renderTemplate(sg, w, "new", map[string]interface{}{
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
	return renderTemplate(sg, w, "index", map[string]interface{}{
		"title":         "Cloud Accounts",
		"uiBasePath":    "/ui/cloud_accounts",
		"apiBasePath":   "/api/v0/cloud_accounts",
		"fields":        fields,
		"showNewLink":   true,
		"showStatusCol": false,
		"newOptions": map[string]string{
			"aws":          "AWS",
			"digitalocean": "DigitalOcean",
			"openstack":    "OpenStack",
			"gce":          "GCE",
			"packet":       "Packet",
		},
		"batchActionPaths": map[string]map[string]string{
			"Delete": map[string]string{
				"method":       "DELETE",
				"relativePath": "",
			},
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
	return renderTemplate(sg, w, "show", map[string]interface{}{
		"title": "Cloud Accounts",
		"model": item,
	})
}
