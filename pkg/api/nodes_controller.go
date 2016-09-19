package api

import (
	"net/http"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

func ListNodes(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	return handleList(core, r, new(model.Node))
}

func CreateNode(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.Node)
	if err := decodeBodyInto(r, item); err != nil {
		return nil, err
	}
	if err := core.Nodes.Create(item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusCreated)
}

func UpdateNode(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	item := new(model.Node)
	if err := decodeBodyInto(r, item); err != nil {
		return nil, err
	}
	if err := core.Nodes.Update(id, new(model.Node), item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}

func GetNode(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.Node)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.Nodes.Get(id, item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusOK)
}

func DeleteNode(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.Node)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.Nodes.Delete(id, item).Async(); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}
