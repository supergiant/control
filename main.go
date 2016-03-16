package main

import (
	"fmt"
	"guber"
	"log"
	"net/http"
	"os"
	"supergiant/core/controller"
	"supergiant/core/job"
	"supergiant/core/storage"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	// StrictSlash will redirect /apps to /apps/
	// otherwise mux will simply not match /apps/
	router.StrictSlash(true)

	db := storage.NewClient([]string{"http://localhost:2379"})

	var (
		kHost = os.Getenv("K_HOST")
		kUser = os.Getenv("K_USER")
		kPass = os.Getenv("K_PASS")
	)
	// TODO
	if kHost == "" {
		panic("K_HOST required")
	}
	if kUser == "" {
		panic("K_USER required")
	}
	if kPass == "" {
		panic("K_PASS required")
	}
	kube := guber.NewClient(kHost, kUser, kPass)

	// TODO
	go job.NewWorker(db, kube).Work()

	controller.NewAppController(router, db)
	controller.NewImageRepoController(router, db)
	controller.NewComponentController(router, db)
	controller.NewDeploymentController(router, db)
	controller.NewInstanceController(router, db)
	controller.NewReleaseController(router, db)

	fmt.Println("Serving API on port :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
