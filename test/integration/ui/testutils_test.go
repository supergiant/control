package ui

import (
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/server"
)

func newTestServer() *server.Server {
	c := new(core.Core)
	c.PublishHost = "localhost"
	c.HTTPPort = "10000"
	c.UIEnabled = true

	c.Log = logrus.New()

	srv, err := server.New(c)
	if err != nil {
		c.Log.Warn(err, "waiting a second")
		time.Sleep(time.Second)
		srv, err = server.New(c)
		if err != nil {
			panic(err)
		}
	}
	return srv
}
