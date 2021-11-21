package main

import (
	"log"
	"net/http"
	"todo-go/controllers"
	"todo-go/databases"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	h := controllers.NewBaseHandler(databases.NewInMemoryDatabase())

	r.HandleFunc("/", h.RootHandler).Methods("GET")
	r.HandleFunc("/tasks", h.GetTasks).Methods("GET")
	r.HandleFunc("/task", h.CreateTask).Methods("POST")
	r.HandleFunc("/task/{id:[0-9]+}", h.GetTaskByID).Methods("GET")
	r.HandleFunc("/task/{id:[0-9]+}", h.UpdateTask).Methods("PUT")
	r.HandleFunc("/task/{id:[0-9]+}", h.DeleteTask).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", r))
}
