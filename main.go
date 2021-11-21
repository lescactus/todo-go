package main

import (
	"log"
	"net/http"
	"os"
	"todo-go/controllers"
	"todo-go/databases"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	h := controllers.NewBaseHandler(databases.NewInMemoryDatabase())

	r.Handle("/", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(h.RootHandler))).Methods("GET")
	r.Handle("/tasks", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(h.GetTasks))).Methods("GET")
	r.Handle("/task", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(h.CreateTask))).Methods("POST")
	r.Handle("/task/{id:[0-9]+}", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(h.GetTaskByID))).Methods("GET")
	r.Handle("/task/{id:[0-9]+}", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(h.UpdateTask))).Methods("PUT")
	r.Handle("/task/{id:[0-9]+}", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(h.DeleteTask))).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", r))
}
