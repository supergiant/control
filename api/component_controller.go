package api

import (
	"net/http"

	"github.com/supergiant/supergiant/api/task"
	"github.com/supergiant/supergiant/core"
	"github.com/supergiant/supergiant/types"
)

type ComponentController struct {
	core *core.Core
}

func (c *ComponentController) Create(w http.ResponseWriter, r *http.Request) {
	app, err := loadApp(c.core, w, r)
	if err != nil {
		return
	}

	component := app.Components().New()
	if err := unmarshalBodyInto(w, r, component); err != nil {
		return
	}

	component, err = app.Components().Create(component)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := marshalBody(w, component)
	if err != nil {
		return
	}
	renderWithStatusCreated(w, body)
}

func (c *ComponentController) Index(w http.ResponseWriter, r *http.Request) {
	app, err := loadApp(c.core, w, r)
	if err != nil {
		return
	}

	components, err := app.Components().List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := marshalBody(w, components)
	if err != nil {
		return
	}
	renderWithStatusOK(w, body)
}

func (c *ComponentController) Show(w http.ResponseWriter, r *http.Request) {
	component, err := loadComponent(c.core, w, r)
	if err != nil {
		return
	}

	body, err := marshalBody(w, component)
	if err != nil {
		return
	}
	renderWithStatusOK(w, body)
}

func (c *ComponentController) Delete(w http.ResponseWriter, r *http.Request) {
	component, err := loadComponent(c.core, w, r)
	if err != nil {
		return
	}

	msg := &task.DeleteComponentMessage{
		AppName:       component.App().Name,
		ComponentName: component.Name,
	}
	_, err = c.core.Tasks().Start(types.TaskTypeDeleteComponent, msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := marshalBody(w, component)
	if err != nil {
		return
	}
	renderWithStatusAccepted(w, body)
}
