package main

import (
	_ "Task/constants"
	_ "Task/database"
	"Task/handlers"
	repo "Task/repository"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	baseURL := "/api/v1.0.0"
	getRepoFuncs := &repo.Repository{}
	t := &handlers.TaskHandler{Repo: getRepoFuncs}
	r := mux.NewRouter()
	r.HandleFunc(baseURL+"/students", t.FetchAll).Methods("GET")
	r.HandleFunc(baseURL+"/students/{id}", t.FetchByID).Methods("GET")
	r.HandleFunc(baseURL+"/students/{id}", t.UpdateByID).Methods("PUT")
	r.HandleFunc(baseURL+"/students/{id}", t.DeleteByID).Methods("DELETE")
	r.HandleFunc(baseURL+"/students", t.Create).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", r))
}
