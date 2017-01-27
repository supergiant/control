package api

import (
	"net/http"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

func ListHelmRepos(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	return handleList(core, r, new(model.HelmRepo), new(model.HelmRepoList))
}

func CreateHelmRepo(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.HelmRepo)
	if err := decodeBodyInto(r, item); err != nil {
		return nil, err
	}
	if err := core.HelmRepos.Create(item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusCreated)
}

func UpdateHelmRepo(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	item := new(model.HelmRepo)
	if err := decodeBodyInto(r, item); err != nil {
		return nil, err
	}
	if err := core.HelmRepos.Update(id, new(model.HelmRepo), item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}

func GetHelmRepo(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.HelmRepo)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.HelmRepos.Get(id, item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusOK)
}

func DeleteHelmRepo(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.HelmRepo)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.HelmRepos.Delete(id, item).Async(); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}
