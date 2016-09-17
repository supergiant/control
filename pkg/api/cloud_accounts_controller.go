package api

import (
	"net/http"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

func ListCloudAccounts(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	return handleList(core, r, new(model.CloudAccount), new(model.CloudAccountList))
}

func CreateCloudAccount(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.CloudAccount)
	if err := decodeBodyInto(r, item); err != nil {
		return nil, err
	}
	if err := core.CloudAccounts.Create(item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusCreated)
}

func UpdateCloudAccount(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	item := new(model.CloudAccount)
	if err := decodeBodyInto(r, item); err != nil {
		return nil, err
	}
	if err := core.CloudAccounts.Update(id, new(model.CloudAccount), item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}

func GetCloudAccount(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.CloudAccount)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.CloudAccounts.Get(id, item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusOK)
}

func DeleteCloudAccount(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.CloudAccount)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.CloudAccounts.Delete(id, item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}
