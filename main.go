package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/supergiant/supergiant/api"
	"github.com/supergiant/supergiant/api/task"
	"github.com/supergiant/supergiant/core"
)

func main() {
	core := core.New()

	// TODO should probably be able to say api.New(), because we shouldn't have to import task here
	// NOTE using pool size of 4
	go task.NewSupervisor(core, 20).Run()

	router := api.NewRouter(core)

	fmt.Println("Serving API on port :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
