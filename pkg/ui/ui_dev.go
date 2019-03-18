// +build dev

package ui

import (
	"net/http"
	"github.com/gorilla/mux"
	"os"
	"github.com/pkg/errors"
	"github.com/supergiant/control/pkg/controlplane/config"
)


func ServeUI(cfg *config.Config, router *mux.Router) error {
	if _, err := os.Stat(cfg.UIDir); err != nil {
		return errors.Wrap(err, "no ui directory found")
	}

	router.PathPrefix("/").Handler(trimPrefix(http.FileServer(http.Dir(cfg.UIDir))))
	return nil
}

