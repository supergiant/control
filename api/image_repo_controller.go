package api

import (
	"net/http"

	"github.com/supergiant/supergiant/core"
)

type ImageRepoController struct {
	core *core.Core
}

func (c *ImageRepoController) Create(w http.ResponseWriter, r *http.Request) {
	repo := c.core.ImageRepos().New()

	if err := UnmarshalBodyInto(w, r, repo); err != nil {
		return
	}

	repo, err := c.core.ImageRepos().Create(repo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := MarshalBody(w, repo)
	if err != nil {
		return
	}
	RenderWithStatusCreated(w, body)
}

func (c *ImageRepoController) Delete(w http.ResponseWriter, r *http.Request) {
	repo, err := LoadImageRepo(c.core, w, r)
	if err != nil {
		return
	}
	if err = repo.Delete(); err != nil {
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
