package api

import (
	"net/http"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

func ListReleases(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	return handleList(core, r, new(model.Release), new(model.ReleaseList))
}

func CreateRelease(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.Release)
	if err := decodeBodyInto(r, item); err != nil {
		return nil, err
	}
	if err := core.Releases.Create(item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusCreated)
}

func UpdateRelease(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	item := new(model.Release)
	if err := decodeBodyInto(r, item); err != nil {
		return nil, err
	}
	if err := core.Releases.Update(id, new(model.Release), item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}

func GetRelease(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.Release)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.Releases.Get(id, item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusOK)
}

func DeleteRelease(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.Release)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.Releases.Delete(id, item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}
