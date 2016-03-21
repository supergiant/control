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
// type ImageRepoController struct {
// 	db *storage.Client
// }
//
// type repoListResponse struct {
// 	Items []*model.ImageRepo `json:"items"`
// }
//
// func NewImageRepoController(router *mux.Router, db *storage.Client) *ImageRepoController {
// 	controller := ImageRepoController{db: db}
// 	router.HandleFunc("/registries/dockerhub/repos", controller.Create).Methods("POST")
// 	router.HandleFunc("/registries/dockerhub/repos/{name}", controller.Delete).Methods("DELETE")
// 	return &controller
// }
//
// func (e *ImageRepoController) Create(w http.ResponseWriter, r *http.Request) {
// 	repo := new(model.ImageRepo)
// 	decoder := json.NewDecoder(r.Body)
// 	err := decoder.Decode(repo)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
//
// 	repo, err = e.db.ImageRepoStorage.Create(repo)
//
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
//
// 	out, err := json.Marshal(repo)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
//
// 	w.WriteHeader(http.StatusCreated)
// 	fmt.Fprint(w, string(out))
// }
//
// func (e *ImageRepoController) Delete(w http.ResponseWriter, r *http.Request) {
// 	repoName := mux.Vars(r)["name"]
// 	if err := e.db.ImageRepoStorage.Delete(repoName); err != nil {
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}
//
// 	// fmt.Fprint(w, string(out))
// }
