package api

import (
	"net/http"

	"github.com/supergiant/supergiant/core"
)

type EntrypointController struct {
	core *core.Core
}

func (c *EntrypointController) Create(w http.ResponseWriter, r *http.Request) {
	entrypoint := c.core.Entrypoints().New()

	if err := UnmarshalBodyInto(w, r, entrypoint); err != nil {
		return
	}

	entrypoint, err := c.core.Entrypoints().Create(entrypoint)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := MarshalBody(w, entrypoint)
	if err != nil {
		return
	}
	RenderWithStatusCreated(w, body)
}

func (c *EntrypointController) Index(w http.ResponseWriter, r *http.Request) {
	entrypoints, err := c.core.Entrypoints().List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := MarshalBody(w, entrypoints)
	if err != nil {
		return
	}
	RenderWithStatusOK(w, body)
}

func (c *EntrypointController) Show(w http.ResponseWriter, r *http.Request) {
	entrypoint, err := LoadEntrypoint(c.core, w, r)
	if err != nil {
		return
	}

	body, err := MarshalBody(w, entrypoint)
	if err != nil {
		return
	}
	RenderWithStatusOK(w, body)
}

func (c *EntrypointController) Delete(w http.ResponseWriter, r *http.Request) {
	entrypoint, err := LoadEntrypoint(c.core, w, r)
	if err != nil {
		return
	}
	if err = entrypoint.Delete(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO what do we return on immediate deletes like this? generic OK message?

	// body, err := MarshalBody(w, app)
	// if err != nil {
	// 	return
	// }
	// RenderWithStatusAccepted(w, body)
}
