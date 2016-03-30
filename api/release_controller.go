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
	component, err := loadComponent(c.core, w, r)
	if err != nil {
		return
	}

	release := component.Releases().New()
	if err := unmarshalBodyInto(w, r, release); err != nil {
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

	component.TargetReleaseTimestamp = release.Timestamp
	// (this may should be elsewhere)
	// Set the task ID of the deploy on Component
	component.DeployTaskID = task.ID
	if err := component.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := marshalBody(w, release)
	if err != nil {
		return
	}
	renderWithStatusCreated(w, body)
}

func (c *ReleaseController) Index(w http.ResponseWriter, r *http.Request) {
	component, err := loadComponent(c.core, w, r)
	if err != nil {
		return
	}

	releases, err := component.Releases().List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := marshalBody(w, releases)
	if err != nil {
		return
	}
	renderWithStatusOK(w, body)
}

func (c *ReleaseController) Show(w http.ResponseWriter, r *http.Request) {
	release, err := loadRelease(c.core, w, r)
	if err != nil {
		return
	}

	body, err := marshalBody(w, release)
	if err != nil {
		return
	}
	renderWithStatusOK(w, body)
}
