// +build dev

package controlplane

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func ServeUI(cfg *Config, router *mux.Router) error {
	if _, err := os.Stat(cfg.UIDir); err != nil {
		return errors.Wrap(err, "no ui directory found")
	}

	router.PathPrefix("/").Handler(trimPrefix(http.FileServer(http.Dir(cfg.UIDir))))
	return nil
}

