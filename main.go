package main

import (
	"log"
	"net/http"
	"os"
	"todo-go/config"
	"todo-go/controllers"
	"todo-go/databases"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	cfg := config.New()
	r := mux.NewRouter()

	var h *controllers.BaseHandler
	if cfg.GetString("DATABASE_TYPE") == "memory" {
		h = controllers.NewBaseHandler(databases.NewInMemoryDatabase())
	} else {
		log.Fatal("Invalid DATABASE_TYPE. Must be one of 'memory'")
	}

	r.Handle("/", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(h.RootHandler))).Methods("GET")
	r.Handle("/tasks", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(h.GetTasks))).Methods("GET")
	r.Handle("/task", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(h.CreateTask))).Methods("POST")
	r.Handle("/task/{id:[0-9]+}", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(h.GetTaskByID))).Methods("GET")
	r.Handle("/task/{id:[0-9]+}", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(h.UpdateTask))).Methods("PUT")
	r.Handle("/task/{id:[0-9]+}", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(h.DeleteTask))).Methods("DELETE")
	log.Fatal(http.ListenAndServe(cfg.GetString("APP_ADDR"), r))
}
