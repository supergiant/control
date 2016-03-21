package main

import (
	"fmt"
	"log"
	"net/http"
	"supergiant/core/controller"
	"supergiant/core/model"
)

func main() {

	client := model.NewClient()

	// TODO
	// go job.NewWorker(db, kube).Work()

	router := controller.NewRouter(client)

	fmt.Println("Serving API on port :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
