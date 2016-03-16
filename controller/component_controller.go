package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"supergiant/core/job"
	"supergiant/core/model"
	"supergiant/core/storage"

	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
)

type ComponentController struct {
	db *storage.Client
}

type componentListResponse struct {
	Items []*model.Component `json:"items"`
}

func NewComponentController(router *mux.Router, db *storage.Client) *ComponentController {
	controller := ComponentController{db: db}
	router.HandleFunc("/apps/{app_name}/components", controller.Create).Methods("POST")
	router.HandleFunc("/apps/{app_name}/components", controller.Index).Methods("GET")
	router.HandleFunc("/apps/{app_name}/components/{name}", controller.Show).Methods("GET")
	router.HandleFunc("/apps/{app_name}/components/{name}", controller.Delete).Methods("DELETE")
	return &controller
}

func (s *ComponentController) loadApp(appName string) (*model.App, error) {
	return s.db.AppStorage.Get(appName)
}

func (s *ComponentController) Create(w http.ResponseWriter, r *http.Request) {
	appName := mux.Vars(r)["app_name"]
	app, err := s.loadApp(appName)
	if err != nil {
		http.Error(w, "App does not exist", http.StatusBadRequest)
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
	release, err := s.db.ReleaseStorage.Create(component.Name, app.Name, &model.Release{ID: 1})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	component.CurrentReleaseID = release.ID

	// Create the first Deployment (active_deployment)
	deployment, err := s.db.DeploymentStorage.Create(&model.Deployment{ID: uuid.NewV4().String()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	component.ActiveDeploymentID = deployment.ID

	component, err = s.db.ComponentStorage.Create(app.Name, component)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(component)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO
	msg := &job.DeployComponentMessage{
		AppName:       app.Name,
		ComponentName: component.Name,
	}
	data, err := json.Marshal(msg)

	fmt.Println("DeployComponentMessage: ", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	job := &model.Job{
		Type:   job.JobTypeDeployComponent,
		Data:   data,
		Status: "QUEUED", // TODO
	}
	job, err = s.db.JobStorage.Create(job)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(out))
}

func (s *ComponentController) Index(w http.ResponseWriter, r *http.Request) {
	appName := mux.Vars(r)["app_name"]
	app, err := s.loadApp(appName)
	if err != nil {
		http.Error(w, "App does not exist", http.StatusBadRequest)
		return
	}

	components, err := s.db.ComponentStorage.List(app.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(componentListResponse{Items: components})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(out))
}

func (s *ComponentController) Show(w http.ResponseWriter, r *http.Request) {
	urlVars := mux.Vars(r)
	appName := urlVars["app_name"]
	compName := urlVars["name"]
	app, err := s.loadApp(appName)
	if err != nil {
		http.Error(w, "App does not exist", http.StatusBadRequest)
		return
	}

	component, err := s.db.ComponentStorage.Get(app.Name, compName)
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

func (s *ComponentController) Delete(w http.ResponseWriter, r *http.Request) {
	urlVars := mux.Vars(r)
	appName := urlVars["app_name"]
	compName := urlVars["name"]
	_, err := s.loadApp(appName)
	if err != nil {
		http.Error(w, "App does not exist", http.StatusBadRequest)
		return
	}

	if err = s.db.ComponentStorage.Delete(appName, compName); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// fmt.Fprint(w, string(out))
}
