package api

import (
	"net/http"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

func ListInstances(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	return handleList(core, r, new(model.Instance), new(model.InstanceList))
}

func GetInstance(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.Instance)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.Instances.Get(id, item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusOK)
}

func StartInstance(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.Instance)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.Instances.Start(id, item).Async(); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}

func StopInstance(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.Instance)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.Instances.Stop(id, item).Async(); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}

func DeleteInstance(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.Instance)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.Instances.Delete(id, item).Async(); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}

//------------------------------------------------------------------------------

func ViewInstanceLog(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.Instance)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	includes := []string{"Component.App.Kube"}
	if err := core.Instances.GetWithIncludes(id, item, includes); err != nil {
		return nil, err
	}
	log, err := core.Instances.Log(item)
	if err != nil {
		return nil, err
	}
	return &Response{http.StatusOK, log}, nil
}
