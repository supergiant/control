package controller

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"supergiant/model"
// 	"supergiant/storage"
//
// 	"github.com/gorilla/mux"
// )
//
// type ReleaseController struct {
// 	db *storage.Client
// }
//
// type releaseListResponse struct {
// 	Items []*core.Release `json:"items"`
// }
//
// func NewReleaseController(router *mux.Router, db *storage.Client) *ReleaseController {
// 	controller := ReleaseController{db: db}
// 	router.HandleFunc("/apps/{app_name}/components/{comp_name}/releases", controller.Create).Methods("POST")
// 	router.HandleFunc("/apps/{app_name}/components/{comp_name}/releases", controller.Index).Methods("GET")
// 	router.HandleFunc("/apps/{app_name}/components/{comp_name}/releases/{id}", controller.Show).Methods("GET")
// 	return &controller
// }
//
// func (s *ReleaseController) loadAppComponent(appName string, name string) (*core.App, *core.Component, error) {
// 	app, err := s.db.AppStorage.Get(appName)
// 	comp, err := s.db.ComponentStorage.Get(appName, name)
// 	return app, comp, err
// }
//
// func (s *ReleaseController) Create(w http.ResponseWriter, r *http.Request) {
// 	urlVars := mux.Vars(r)
// 	appName := urlVars["app_name"]
// 	compName := urlVars["comp_name"]
// 	app, comp, err := s.loadAppComponent(appName, compName)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
//
// 	release := new(core.Release)
// 	decoder := json.NewDecoder(r.Body)
// 	err = decoder.Decode(release)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
//
// 	release, err = s.db.ReleaseStorage.Create(app.Name, comp.Name, release)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
//
// 	out, err := json.Marshal(release)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
//
// 	w.WriteHeader(http.StatusCreated)
// 	fmt.Fprint(w, string(out))
// }
//
// func (s *ReleaseController) Index(w http.ResponseWriter, r *http.Request) {
// 	urlVars := mux.Vars(r)
// 	appName := urlVars["app_name"]
// 	compName := urlVars["comp_name"]
// 	app, comp, err := s.loadAppComponent(appName, compName)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
//
// 	releases, err := s.db.ReleaseStorage.List(app.Name, comp.Name)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
//
// 	out, err := json.Marshal(releaseListResponse{Items: releases})
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
//
// 	fmt.Fprint(w, string(out))
// }
//
// func (s *ReleaseController) Show(w http.ResponseWriter, r *http.Request) {
// 	urlVars := mux.Vars(r)
// 	appName := urlVars["app_name"]
// 	compName := urlVars["comp_name"]
// 	releaseID := urlVars["id"]
// 	app, comp, err := s.loadAppComponent(appName, compName)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
//
// 	release, err := s.db.ReleaseStorage.Get(app.Name, comp.Name, releaseID)
// 	if err != nil {
// 		http.Error(w, "Not Found", http.StatusNotFound)
// 		return
// 	}
//
// 	out, err := json.Marshal(release)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
//
// 	fmt.Fprint(w, string(out))
// }
