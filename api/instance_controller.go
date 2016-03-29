package api

import (
	"net/http"

	"github.com/supergiant/supergiant/api/task"
	"github.com/supergiant/supergiant/core"
	"github.com/supergiant/supergiant/types"
)

type InstanceController struct {
	core *core.Core
}

func (c *InstanceController) Index(w http.ResponseWriter, r *http.Request) {
	release, err := LoadRelease(c.core, w, r)
	if err != nil {
		return
	}

	instances := release.Instances().List()

	body, err := MarshalBody(w, instances)
	if err != nil {
		return
	}
	RenderWithStatusOK(w, body)
}

func (c *InstanceController) Show(w http.ResponseWriter, r *http.Request) {
	instance, err := LoadInstance(c.core, w, r)
	if err != nil {
		return
	}

	body, err := MarshalBody(w, instance)
	if err != nil {
		return
	}
	RenderWithStatusOK(w, body)
}

func (c *InstanceController) Start(w http.ResponseWriter, r *http.Request) {
	instance, err := LoadInstance(c.core, w, r)
	if err != nil {
		return
	}

	msg := &task.StartInstanceMessage{
		AppName:          instance.App().Name,
		ComponentName:    instance.Component().Name,
		ReleaseTimestamp: instance.Release().Timestamp,
		ID:               instance.ID,
	}
	_, err = c.core.Tasks().Start(types.TaskTypeStartInstance, msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := MarshalBody(w, instance)
	if err != nil {
		return
	}
	RenderWithStatusAccepted(w, body)
}

func (c *InstanceController) Stop(w http.ResponseWriter, r *http.Request) {
	instance, err := LoadInstance(c.core, w, r)
	if err != nil {
		return
	}

	msg := &task.StopInstanceMessage{
		AppName:          instance.App().Name,
		ComponentName:    instance.Component().Name,
		ReleaseTimestamp: instance.Release().Timestamp,
		ID:               instance.ID,
	}
	_, err = c.core.Tasks().Start(types.TaskTypeStopInstance, msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := MarshalBody(w, instance)
	if err != nil {
		return
	}
	RenderWithStatusAccepted(w, body)
}
