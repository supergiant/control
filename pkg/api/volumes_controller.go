package api

import (
	"net/http"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

func ListVolumes(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	return handleList(core, r, new(model.Volume), new(model.VolumeList))
}

func CreateVolume(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.Volume)
	if err := decodeBodyInto(r, item); err != nil {
		return nil, err
	}
	if err := core.Volumes.Create(item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusCreated)
}

func UpdateVolume(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	item := new(model.Volume)
	if err := decodeBodyInto(r, item); err != nil {
		return nil, err
	}
	if err := core.Volumes.Update(id, new(model.Volume), item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}

func GetVolume(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.Volume)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.Volumes.Get(id, item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusOK)
}

func DeleteVolume(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.Volume)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}
	if err := core.Volumes.Delete(id, item).Async(); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}
