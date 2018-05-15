package ui

import (
	"fmt"
	"net/http"

	"github.com/supergiant/supergiant/pkg/client"
	"github.com/supergiant/supergiant/pkg/model"
)

func NewLoadBalancer(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	return renderTemplate(sg, w, "new", map[string]interface{}{
		"title":      "Load Balancers",
		"formAction": "/ui/load_balancers",
		"model": map[string]interface{}{
			"kube_name": "",
			"name":      "",
			"namespace": "default",
			"selector": map[string]string{
				"key": "value",
			},
			"ports": map[int]int{
				80: 8080,
			},
		},
	})
}

func ListLoadBalancers(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	fields := []map[string]interface{}{
		{
			"title": "Kube name",
			"type":  "field_value",
			"field": "kube_name",
		},
		{
			"title": "Name",
			"type":  "field_value",
			"field": "name",
		},
		{
			"title": "Address",
			"type":  "field_value",
			"field": "address",
		},
	}
	return renderTemplate(sg, w, "index", map[string]interface{}{
		"title":         "Load Balancers",
		"uiBasePath":    "/ui/load_balancers",
		"apiBasePath":   "/api/v0/load_balancers",
		"fields":        fields,
		"showNewLink":   true,
		"showStatusCol": true,
		"actionPaths": map[string]string{
			"Edit": "/ui/load_balancers/{{ ID }}/edit",
		},
		"batchActionPaths": map[string]map[string]string{
			"Delete": {
				"method":       "DELETE",
				"relativePath": "",
			},
		},
	})
}

func CreateLoadBalancer(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	m := new(model.LoadBalancer)
	err := unmarshalFormInto(r, m)
	if err == nil {
		err = sg.LoadBalancers.Create(m)
	}
	if err != nil {
		return renderTemplate(sg, w, "new", map[string]interface{}{
			"title":      "Load Balancers",
			"formAction": "/ui/load_balancers",
			"model":      m,
			"error":      err.Error(),
		})
	}

	http.Redirect(w, r, "/ui/load_balancers", http.StatusFound)
	return nil
}

func GetLoadBalancer(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return err
	}
	item := new(model.LoadBalancer)
	if err := sg.LoadBalancers.Get(id, item); err != nil {
		return err
	}
	return renderTemplate(sg, w, "show", map[string]interface{}{
		"title": "Load Balancers",
		"model": item,
	})
}

func EditLoadBalancer(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return err
	}
	item := new(model.LoadBalancer)
	if err := sg.LoadBalancers.Get(id, item); err != nil {
		return err
	}
	return renderTemplate(sg, w, "new", map[string]interface{}{
		"title":      "Load Balancers",
		"formAction": fmt.Sprintf("/ui/load_balancers/%d", *id),
		"model": map[string]interface{}{
			"selector": item.Selector,
			"ports":    item.Ports,
		},
	})
}

func UpdateLoadBalancer(sg *client.Client, w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return err
	}
	m := new(model.LoadBalancer)
	err = unmarshalFormInto(r, m)
	if err == nil {
		err = sg.LoadBalancers.Update(id, m)
	}
	if err != nil {
		return renderTemplate(sg, w, "new", map[string]interface{}{
			"title":      "Load Balancers",
			"formAction": fmt.Sprintf("/ui/load_balancers/%d", *id),
			"model":      m,
			"error":      err.Error(),
		})
	}
	http.Redirect(w, r, "/ui/load_balancers", http.StatusFound)
	return nil
}
