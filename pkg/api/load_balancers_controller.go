package api

import (
	"net/http"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

func ListLoadBalancers(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	return handleList(core, r, new(model.LoadBalancer), new(model.LoadBalancerList))
}

func CreateLoadBalancer(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.LoadBalancer)
	if err := decodeBodyInto(r, item); err != nil {
		return nil, err
	}
	if err := core.LoadBalancers.Create(item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusCreated)
}

func UpdateLoadBalancer(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	item := new(model.LoadBalancer)
	if err := decodeBodyInto(r, item); err != nil {
		return nil, err
	}
	if err := core.LoadBalancers.Update(id, new(model.LoadBalancer), item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}

func GetLoadBalancer(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.LoadBalancer)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.LoadBalancers.Get(id, item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusOK)
}

func DeleteLoadBalancer(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.LoadBalancer)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.LoadBalancers.Delete(id, item).Async(); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}
