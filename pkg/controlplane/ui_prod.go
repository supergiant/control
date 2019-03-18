// +build prod

package controlplane

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"

	_ "github.com/supergiant/control/statik"
)

func ServeUI(cfg *Config, router *mux.Router) error {
	statikFS, err := fs.New()
	if err != nil {
		return err
	}

	router.PathPrefix("/").Handler(trimPrefix(http.FileServer(statikFS)))
	return nil
}