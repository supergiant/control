package api

import (
	"net/http"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

func ListEntrypointListeners(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	return handleList(core, r, new(model.EntrypointListener), new(model.EntrypointListenerList))
}

func CreateEntrypointListener(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.EntrypointListener)
	if err := decodeBodyInto(r, item); err != nil {
		return nil, err
	}
	if err := core.EntrypointListeners.Create(item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusCreated)
}

func GetEntrypointListener(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.EntrypointListener)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.EntrypointListeners.Get(id, item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusOK)
}

func UpdateEntrypointListener(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	item := new(model.EntrypointListener)
	if err := decodeBodyInto(r, item); err != nil {
		return nil, err
	}
	if err := core.EntrypointListeners.Update(id, new(model.EntrypointListener), item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}

func DeleteEntrypointListener(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.EntrypointListener)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.EntrypointListeners.Delete(id, item).Async(); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}
