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

	if err := unmarshalBodyInto(w, r, app); err != nil {
		return
	}

	app, err := c.core.Apps().Create(app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := marshalBody(w, app)
	if err != nil {
		return
	}
	renderWithStatusCreated(w, body)
}

func (c *AppController) Index(w http.ResponseWriter, r *http.Request) {
	apps, err := c.core.Apps().List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := marshalBody(w, apps)
	if err != nil {
		return
	}
	renderWithStatusOK(w, body)
}

func (c *AppController) Show(w http.ResponseWriter, r *http.Request) {
	app, err := loadApp(c.core, w, r)
	if err != nil {
		return
	}

	body, err := marshalBody(w, app)
	if err != nil {
		return
	}
	renderWithStatusOK(w, body)
}

func (c *AppController) Delete(w http.ResponseWriter, r *http.Request) {
	app, err := loadApp(c.core, w, r)
	if err != nil {
		return
	}

	msg := &task.DeleteAppMessage{AppName: app.Name}
	_, err = c.core.Tasks().Start(types.TaskTypeDeleteApp, msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := marshalBody(w, app)
	if err != nil {
		return
	}
	renderWithStatusAccepted(w, body)
}
