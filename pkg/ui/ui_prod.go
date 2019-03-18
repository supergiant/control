// +build prod

package ui

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"

	"github.com/supergiant/control/pkg/controlplane/config"
	_ "github.com/supergiant/control/statik"
)

func ServeUI(cfg *config.Config, router *mux.Router) error {
	statikFS, err := fs.New()
	if err != nil {
		return err
	}

	router.PathPrefix("/").Handler(trimPrefix(http.FileServer(statikFS)))
	return nil
}