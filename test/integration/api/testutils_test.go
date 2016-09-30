package api

import (
	"encoding/json"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/server"
)

func newTestServer() *server.Server {
	c := new(core.Core)
	c.LogLevel = "fatal"
	c.PublishHost = "localhost"
	c.HTTPPort = "9999"
	c.SQLiteFile = "../../../tmp/test.db"

	wipeAndInitialize(c)

	srv, err := server.New(c)
	if err != nil {
		panic(err)
	}
	return srv
}

func wipeDatabase(c *core.Core) {
	// // NOTE some sporadic IO errors thrown here
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println("Recovered in wipeDatabase: ", err, "waiting a second")
	// 		time.Sleep(time.Second)
	// 	}
	// }()
	// os.Remove(c.SQLiteFile)
	c.DB.Delete(&model.User{})
	c.DB.Delete(&model.Kube{})
	c.DB.Delete(&model.KubeResource{})
	c.DB.Delete(&model.CloudAccount{})
	c.DB.Delete(&model.Volume{})
	c.DB.Delete(&model.Entrypoint{})
	c.DB.Delete(&model.EntrypointListener{})
	c.DB.Delete(&model.Node{})
}

func wipeAndInitialize(c *core.Core) {
	// wipeDatabase(c)
	if err := c.InitializeForeground(); err != nil {
		panic(err)
	}
	wipeDatabase(c)
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

func newRawMessage(str string) *json.RawMessage {
	rawmsg := json.RawMessage([]byte(str))
	return &rawmsg
}
