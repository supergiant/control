package controller

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"supergiant/core/model"
// 	"supergiant/core/storage"
//
// 	"github.com/gorilla/mux"
// )
//
// type InstanceController struct {
// 	db *storage.Client
// }
//
// type instanceListResponse struct {
// 	Items []*model.Instance `json:"items"`
// }
//
// func NewInstanceController(router *mux.Router, db *storage.Client) *InstanceController {
// 	controller := InstanceController{db: db}
// 	router.HandleFunc("/deployments/{deployment_id}/instances", controller.Index).Methods("GET")
// 	router.HandleFunc("/deployments/{deployment_id}/instances/{id}", controller.Show).Methods("GET")
// 	router.HandleFunc("/deployments/{deployment_id}/instances/{id}", controller.Delete).Methods("DELETE")
// 	return &controller
// }
//
// func (s *InstanceController) loadDeployment(id string) (*model.Deployment, error) {
// 	return s.db.DeploymentStorage.Get(id)
// }
//
// func (s *InstanceController) Index(w http.ResponseWriter, r *http.Request) {
// 	deploymentID := mux.Vars(r)["deployment_id"]
// 	_, err := s.loadDeployment(deploymentID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
//
// 	instances, err := s.db.InstanceStorage.List(deploymentID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
//
// 	out, err := json.Marshal(instanceListResponse{Items: instances})
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
//
// 	fmt.Fprint(w, string(out))
// }
//
// func (s *InstanceController) Show(w http.ResponseWriter, r *http.Request) {
// 	urlVars := mux.Vars(r)
// 	deploymentID := urlVars["deployment_id"]
// 	id := urlVars["id"]
// 	_, err := s.loadDeployment(deploymentID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
//
// 	component, err := s.db.InstanceStorage.Get(deploymentID, id)
// 	if err != nil {
// 		http.Error(w, "Not Found", http.StatusNotFound)
// 		return
// 	}
//
// 	out, err := json.Marshal(component)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
//
// 	fmt.Fprint(w, string(out))
// }
//
// func (s *InstanceController) Delete(w http.ResponseWriter, r *http.Request) {
// 	urlVars := mux.Vars(r)
// 	deploymentID := urlVars["deployment_id"]
// 	id := urlVars["id"]
// 	_, err := s.loadDeployment(deploymentID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	if err := s.db.InstanceStorage.Delete(deploymentID, id); err != nil {
// 		http.Error(w, "Not Found", http.StatusNotFound)
// 		return
// 	}
//
// 	// fmt.Fprint(w, string(out))
// }
