package api

import (
	"net/http"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

func ListApps(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	return handleList(core, r, new(model.App), new(model.AppList))
}

func CreateApp(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.App)
	if err := decodeBodyInto(r, item); err != nil {
		return nil, err
	}
	if err := core.Apps.Create(item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusCreated)
}

func UpdateApp(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	item := new(model.App)
	if err := decodeBodyInto(r, item); err != nil {
		return nil, err
	}
	if err := core.Apps.Update(id, new(model.App), item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}

func GetApp(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.App)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.Apps.Get(id, item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusOK)
}

func DeleteApp(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.App)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.Apps.Delete(id, item).Async(); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}
