package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"supergiant/core/model"

	"github.com/gorilla/mux"
)

type ComponentController struct {
	client *model.Client
}

func (c *ComponentController) Create(w http.ResponseWriter, r *http.Request) {
	appName := mux.Vars(r)["app_name"]
	app, err := c.client.Apps().Get(appName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	component := new(model.Component)
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(component)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create the initial Release
	// release, err := s.db.ReleaseStorage.Create(component.Name, app.Name, &model.Release{ID: 1})
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// component.CurrentReleaseID = release.ID
	//
	// // Create the first Deployment (active_deployment)
	// deployment, err := s.db.DeploymentStorage.Create(&model.Deployment{ID: uuid.NewV4().String()}) // ------------------------- should not be UUID -- maybe 4-char string
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// component.ActiveDeploymentID = deployment.ID

	component, err = app.Components().Create(component)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(component)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// // TODO
	// msg := &job.CreateComponentMessage{
	// 	AppName:       app.Name,
	// 	ComponentName: component.Name,
	// }
	// data, err := json.Marshal(msg)
	//
	// fmt.Println("CreateComponentMessage: ", data)
	//
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// job := &model.Job{
	// 	Type:   job.JobTypeCreateComponent,
	// 	Data:   data,
	// 	Status: "QUEUED", // TODO
	// }
	// job, err = s.db.JobStorage.Create(job)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(out))
}

func (c *ComponentController) Index(w http.ResponseWriter, r *http.Request) {
	appName := mux.Vars(r)["app_name"]
	app, err := c.client.Apps().Get(appName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	components, err := app.Components().List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(components)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(out))
}

func (c *ComponentController) Show(w http.ResponseWriter, r *http.Request) {
	urlVars := mux.Vars(r)
	appName := urlVars["app_name"]
	compName := urlVars["name"]
	app, err := c.client.Apps().Get(appName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	component, err := app.Components().Get(compName)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	out, err := json.Marshal(component)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(out))
}

func (c *ComponentController) Delete(w http.ResponseWriter, r *http.Request) {
	urlVars := mux.Vars(r)
	appName := urlVars["app_name"]
	compName := urlVars["name"]
	app, err := c.client.Apps().Get(appName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = app.Components().Delete(compName); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// fmt.Fprint(w, string(out))
}
