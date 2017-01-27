package api

import (
	"net/http"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

func ListHelmReleases(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	return handleList(core, r, new(model.HelmRelease), new(model.HelmReleaseList))
}

func CreateHelmRelease(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.HelmRelease)
	if err := decodeBodyInto(r, item); err != nil {
		return nil, err
	}
	if err := core.HelmReleases.Create(item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusCreated)
}

func UpdateHelmRelease(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	item := new(model.HelmRelease)
	if err := decodeBodyInto(r, item); err != nil {
		return nil, err
	}
	if err := core.HelmReleases.Update(id, new(model.HelmRelease), item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}

func GetHelmRelease(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.HelmRelease)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.HelmReleases.Get(id, item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusOK)
}

func DeleteHelmRelease(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.HelmRelease)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.HelmReleases.Delete(id, item).Async(); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}
