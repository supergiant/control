package api

import (
	"net/http"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

func ListKubeResources(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	return handleList(core, r, new(model.KubeResource), new(model.KubeResourceList))
}

func CreateKubeResource(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.KubeResource)
	if err := decodeBodyInto(r, item); err != nil {
		return nil, err
	}
	if err := core.KubeResources.Create(item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusCreated)
}

func UpdateKubeResource(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	item := new(model.KubeResource)
	if err := decodeBodyInto(r, item); err != nil {
		return nil, err
	}
	if err := core.KubeResources.Update(id, new(model.KubeResource), item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}

func GetKubeResource(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.KubeResource)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.KubeResources.Get(id, item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusOK)
}

func StartKubeResource(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.KubeResource)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.KubeResources.Start(id, item).Async(); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}

func StopKubeResource(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.KubeResource)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.KubeResources.Stop(id, item).Async(); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}

func DeleteKubeResource(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.KubeResource)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.KubeResources.Delete(id, item).Async(); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}
