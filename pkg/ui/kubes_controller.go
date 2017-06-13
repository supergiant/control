package ui

import (
	"net/http"

	"github.com/supergiant/supergiant/pkg/client"
	"github.com/supergiant/supergiant/pkg/model"
)

//NewKube holds default inffo for the UI kubes object.
func NewKube(sg *client.Client, w http.ResponseWriter, r *http.Request) error {

	// TODO we shouldn't use a map for the model, because attr ordering is screwed
	// up. But it's difficult to initialize blank values with omitemptys (which
	// are needed for certain things), so we should probably have special structs.

	var m map[string]interface{}
	switch r.URL.Query().Get("option") {
	// case "aws":
	case "digitalocean":
		m = map[string]interface{}{
			"cloud_account_name": "",
			"name":               "",
			"master_node_size":   "1gb",
			"node_sizes": []string{
				"1gb",
				"2gb",
				"4gb",
				"8gb",
				"16gb",
				"32gb",
				"48gb",
				"64gb",
			},
			"digitalocean_config": map[string]interface{}{
				"region":              "nyc1",
				"ssh_key_fingerprint": "",
			},
		}
	case "openstack":
		m = map[string]interface{}{
			"cloud_account_name": "",
			"name":               "",
			"master_node_size":   "m1.smaller",
			"node_sizes": []string{
				"m1.smaller",
				"m1.small",
			},
			"openstack_config": map[string]interface{}{
				"image_name":          "CoreOS",
				"region":              "RegionOne",
				"ssh_key_fingerprint": "",
				"ssh_pub_key":         "",
			},
		}
	case "gce":
		m = map[string]interface{}{
			"cloud_account_name": "",
			"name":               "",
			"master_node_size":   "n1-standard-1",
			"node_sizes": []string{
				"n1-standard-1",
				"n1-standard-2",
				"n1-standard-4",
				"n1-standard-8",
			},
			"gce_config": map[string]interface{}{
				"zone":        "us-east1-b",
				"ssh_pub_key": "",
			},
		}
	default: // just default to AWS if option not provided, or mismatched
		m = map[string]interface{}{
			"cloud_account_name": "",
			"name":               "",
			"master_node_size":   "m4.large",
			"node_sizes": []string{
				"m4.large",
				"m4.xlarge",
				"m4.2xlarge",
				"m4.4xlarge",
			},
			"aws_config": map[string]interface{}{
				"region":       "us-east-1",
				"vpc_ip_range": "172.20.0.0/16",
			},
		}
	}

	return renderTemplate(sg, w, "new", map[string]interface{}{
		"title":      "Kubes",
		"formAction": "/ui/kubes",
		"model":      m,
	})
}

func CreateKube(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	m := new(model.Kube)
	err := unmarshalFormInto(r, m)
	if err == nil {
		err = sg.Kubes.Create(m)
	}
	if err != nil {
		return renderTemplate(sg, w, "new", map[string]interface{}{
			"title":      "Kubes",
			"formAction": "/ui/kubes",
			"model":      m,
			"error":      err.Error(),
		})
	}
	http.Redirect(w, r, "/ui/kubes", http.StatusFound)
	return nil
}

func ListKubes(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	fields := []map[string]interface{}{
		{
			"title": "Name",
			"type":  "field_value",
			"field": "name",
		},
		{
			"title": "Master Size",
			"type":  "field_value",
			"field": "master_node_size",
		},
	}
	return renderTemplate(sg, w, "index", map[string]interface{}{
		"title":         "Kubes",
		"uiBasePath":    "/ui/kubes",
		"apiBasePath":   "/api/v0/kubes",
		"fields":        fields,
		"showNewLink":   true,
		"showStatusCol": true,
		"newOptions": map[string]string{
			"aws":          "AWS",
			"digitalocean": "DigitalOcean",
			"openstack":    "OpenStack",
			"gce":          "GCE",
		},
		"batchActionPaths": map[string]map[string]string{
			"Reprovision": map[string]string{
				"method":       "POST",
				"relativePath": "/provision",
			},
			"Delete": map[string]string{
				"method":       "DELETE",
				"relativePath": "",
			},
		},
	})
}

func GetKube(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return err
	}
	item := new(model.Kube)
	if err := sg.Kubes.Get(id, item); err != nil {
		return err
	}
	return renderTemplate(sg, w, "show", map[string]interface{}{
		"title": "Kubes",
		"model": item,
	})
}
