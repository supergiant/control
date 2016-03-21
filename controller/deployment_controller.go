package controller

// import (
// 	"net/http"
// 	"supergiant/core/storage"
//
// 	"github.com/gorilla/mux"
// )
//
// type DeploymentController struct {
// 	db *storage.Client
// }
//
// // type releaseListResponse struct {
// // 	Items []*model.Deployment `json:"items"`
// // }
//
// func NewDeploymentController(router *mux.Router, db *storage.Client) *DeploymentController {
// 	controller := DeploymentController{db: db}
// 	router.HandleFunc("/deployments/{id}", controller.Delete).Methods("DELETE")
// 	return &controller
// }
//
// func (s *DeploymentController) Delete(w http.ResponseWriter, r *http.Request) {
// 	id := mux.Vars(r)["id"]
// 	if err := s.db.DeploymentStorage.Delete(id); err != nil {
// 		http.Error(w, "Not Found", http.StatusNotFound)
// 		return
// 	}
//
// 	// fmt.Fprint(w, string(out))
// }
