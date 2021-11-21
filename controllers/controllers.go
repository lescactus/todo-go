package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"todo-go/models"

	"github.com/gorilla/mux"
)

// BaseHandler will hold everything that controller needs
type BaseHandler struct {
	taskRepo models.TaskRepository
}

// NewBaseHandler returns a new BaseHandler
func NewBaseHandler(taskRepo models.TaskRepository) *BaseHandler {
	return &BaseHandler{
		taskRepo: taskRepo,
	}
}

func (h *BaseHandler) RootHandler(w http.ResponseWriter, r *http.Request) {
	t := h.taskRepo.GetAllTasks()
	resp, err := json.Marshal(t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (h *BaseHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	t := h.taskRepo.GetAllTasks()
	resp, err := json.Marshal(t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (h *BaseHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	rBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var t models.Task
	if err := json.Unmarshal(rBody, &t); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	id, err := h.taskRepo.CreateTask(t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	n, err := h.taskRepo.GetTaskByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	resp, err := json.Marshal(n)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (h *BaseHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	// Get id from URL and parse it to uint64
	muxID := mux.Vars(r)["id"]
	id, err := strconv.ParseUint(muxID, 0, 0)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Lookup for the task with the requested id
	t, err := h.taskRepo.GetTaskByID(id)
	if err != nil {
		// If no task with the given id exists, respond 404
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	}

	resp, err := json.Marshal(t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application")
	w.Write(resp)
}

func (h *BaseHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	rBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var t models.Task
	if err := json.Unmarshal(rBody, &t); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if err := h.taskRepo.UpdateTask(t); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("resource updated successfully"))
}

func (h *BaseHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Get id from URL and parse it to uint64
	muxID := mux.Vars(r)["id"]
	id, err := strconv.ParseUint(muxID, 0, 0)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if err := h.taskRepo.DeleteTask(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}
