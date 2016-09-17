package api

import (
	"net/http"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

func ListPrivateImageKeys(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	return handleList(core, r, new(model.PrivateImageKey), new(model.PrivateImageKeyList))
}

func CreatePrivateImageKey(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.PrivateImageKey)
	if err := decodeBodyInto(r, item); err != nil {
		return nil, err
	}
	if err := core.PrivateImageKeys.Create(item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusCreated)
}

func UpdatePrivateImageKey(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	item := new(model.PrivateImageKey)
	if err := decodeBodyInto(r, item); err != nil {
		return nil, err
	}
	if err := core.PrivateImageKeys.Update(id, new(model.PrivateImageKey), item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}

func GetPrivateImageKey(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.PrivateImageKey)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.PrivateImageKeys.Get(id, item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusOK)
}

func DeletePrivateImageKey(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.PrivateImageKey)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.PrivateImageKeys.Delete(id, item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}
