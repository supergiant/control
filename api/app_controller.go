package api

import (
	"net/http"

	"github.com/supergiant/supergiant/api/task"
	"github.com/supergiant/supergiant/core"
	"github.com/supergiant/supergiant/types"
)

type AppController struct {
	core *core.Core
}

func (c *AppController) Create(w http.ResponseWriter, r *http.Request) {
	app := c.core.Apps().New()

	if err := UnmarshalBodyInto(w, r, app); err != nil {
		return
	}

	app, err := c.core.Apps().Create(app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := MarshalBody(w, app)
	if err != nil {
		return
	}
	RenderWithStatusCreated(w, body)
}

func (c *AppController) Index(w http.ResponseWriter, r *http.Request) {
	apps, err := c.core.Apps().List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := MarshalBody(w, apps)
	if err != nil {
		return
	}
	RenderWithStatusOK(w, body)
}

func (c *AppController) Show(w http.ResponseWriter, r *http.Request) {
	app, err := LoadApp(c.core, w, r)
	if err != nil {
		return
	}

	body, err := MarshalBody(w, app)
	if err != nil {
		return
	}
	RenderWithStatusOK(w, body)
}

func (c *AppController) Delete(w http.ResponseWriter, r *http.Request) {
	app, err := LoadApp(c.core, w, r)
	if err != nil {
		return
	}

	msg := &task.DeleteAppMessage{AppName: app.Name}
	_, err = c.core.Tasks().Start(types.TaskTypeDeleteApp, msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := MarshalBody(w, app)
	if err != nil {
		return
	}
	RenderWithStatusAccepted(w, body)
}
