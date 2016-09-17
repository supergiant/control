package api

import (
	"net/http"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

func ListVolumes(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	return handleList(core, r, new(model.Volume), new(model.VolumeList))
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
