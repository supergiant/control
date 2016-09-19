package api

import (
	"os"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/server"
)

func newTestServer() *server.Server {
	c := new(core.Core)
	c.PublishHost = "localhost"
	c.HTTPPort = "9999"
	c.SQLiteFile = "../../../tmp/test.db"

	// Wipe database
	os.Remove(c.SQLiteFile)

	if err := c.InitializeForeground(); err != nil {
		panic(err)
	}

	srv, err := server.New(c)
	if err != nil {
		panic(err)
	}
	return srv
}

func createUser(c *core.Core) *model.User {
	user := &model.User{
		Username: "user",
		Password: "password",
	}
	c.Users.Create(user)
	return user
}

func createAdmin(c *core.Core) *model.User {
	admin := &model.User{
		Username: "bossman",
		Password: "password",
		Role:     "admin",
	}
	c.Users.Create(admin)
	return admin
}

func createUserAndAdmin(c *core.Core) (*model.User, *model.User) {
	return createUser(c), createAdmin(c)
}
