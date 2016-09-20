package api

import (
	"net/http"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

func ListComponents(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	return handleList(core, r, new(model.Component))
}

func CreateComponent(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.Component)
	if err := decodeBodyInto(r, item); err != nil {
		return nil, err
	}
	if err := core.Components.Create(item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusCreated)
}

func UpdateComponent(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	item := new(model.Component)
	if err := decodeBodyInto(r, item); err != nil {
		return nil, err
	}
	if err := core.Components.Update(id, new(model.Component), item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}

func GetComponent(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.Component)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.Components.GetWithIncludes(id, item, parseIncludes(r)); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusOK)
}

func DeployComponent(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.Component)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.Components.Deploy(user, id, item).Async(); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}

func DeleteComponent(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.Component)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.Components.Delete(id, item).Async(); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}
