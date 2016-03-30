package api

import (
	"net/http"

	"github.com/supergiant/supergiant/api/task"
	"github.com/supergiant/supergiant/core"
	"github.com/supergiant/supergiant/types"
)

type ReleaseController struct {
	core *core.Core
}

func (c *ReleaseController) Create(w http.ResponseWriter, r *http.Request) {
	component, err := LoadComponent(c.core, w, r)
	if err != nil {
		return
	}

	release := component.Releases().New()
	if err := UnmarshalBodyInto(w, r, release); err != nil {
		return
	}

	release, err = component.Releases().Create(release)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	msg := &task.DeployComponentMessage{
		AppName:       release.App().Name,
		ComponentName: release.Component().Name,
	}
	task, err := c.core.Tasks().Start(types.TaskTypeDeployComponent, msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	component.TargetReleaseID = release.ID
	// (this may should be elsewhere)
	// Set the task ID of the deploy on Component
	component.DeployTaskID = task.ID
	if err := component.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := MarshalBody(w, release)
	if err != nil {
		return
	}
	RenderWithStatusCreated(w, body)
}

func (c *ReleaseController) Index(w http.ResponseWriter, r *http.Request) {
	component, err := LoadComponent(c.core, w, r)
	if err != nil {
		return
	}

	releases, err := component.Releases().List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := MarshalBody(w, releases)
	if err != nil {
		return
	}
	RenderWithStatusOK(w, body)
}

func (c *ReleaseController) Show(w http.ResponseWriter, r *http.Request) {
	release, err := LoadRelease(c.core, w, r)
	if err != nil {
		return
	}

	body, err := MarshalBody(w, release)
	if err != nil {
		return
	}
	RenderWithStatusOK(w, body)
}
