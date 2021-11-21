package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"todo-go/databases"
	"todo-go/models"

	"github.com/stretchr/testify/assert"
)

type MockTaskRepository struct {
	taskRepo databases.InMemoryDatabase
}

func (m MockTaskRepository) GetAllTasks() []models.Task {
	return []models.Task{
		{
			Id: 0,
		}, {
			Id: 1,
		},
	}
}

func (m MockTaskRepository) CreateTask(t models.Task) (uint64, error) {
	return uint64(1), nil
}

func (m MockTaskRepository) UpdateTask(t models.Task) error {
	return nil
}

func (m MockTaskRepository) DeleteTask(id uint64) error {
	return nil
}

func (m MockTaskRepository) GetTaskByID(id uint64) (*models.Task, error) {
	if id != 99 {
		return &models.Task{Id: id}, nil
	}
	return nil, fmt.Errorf("error: task with id %d doesn't exist", id)
}

var m MockTaskRepository

func TestNewBaseHandler(t *testing.T) {
	t.Run("Create a BaseHandler", func(t *testing.T) {
		var db = databases.NewInMemoryDatabase()
		h := NewBaseHandler(db)
		assert.NotEmpty(t, h)
	})
}

func TestRootHandler(t *testing.T) {
	t.Run("Get a response", func(t *testing.T) {
		h := &BaseHandler{&m}
		req, _ := http.NewRequest("GET", "/", nil)
		res := httptest.NewRecorder()

		h.RootHandler(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.GreaterOrEqual(t, res.Body.Len(), 1)
		assert.Equal(t, []byte(`[{"id":0,"title":"","body":"","priority":0,"status":""},{"id":1,"title":"","body":"","priority":0,"status":""}]`), res.Body.Bytes())
		assert.Equal(t, "application/json", res.Result().Header["Content-Type"][0])
	})
}

func TestGetTasks(t *testing.T) {
	t.Run("Get all tasks", func(t *testing.T) {
		h := &BaseHandler{&m}
		req, _ := http.NewRequest("GET", "/", nil)
		res := httptest.NewRecorder()

		h.GetTasks(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.GreaterOrEqual(t, res.Body.Len(), 1)
		assert.Equal(t, []byte(`[{"id":0,"title":"","body":"","priority":0,"status":""},{"id":1,"title":"","body":"","priority":0,"status":""}]`), res.Body.Bytes())
		assert.Equal(t, "application/json", res.Result().Header["Content-Type"][0])
	})
}

func TestCreateTask(t *testing.T) {
	t.Run("Create task without id", func(t *testing.T) {
		h := &BaseHandler{&m}
		req, _ := http.NewRequest("POST", "/task", strings.NewReader(`{"title":"new title","body":"new body","priority":0,"status":"TODO"}`))
		res := httptest.NewRecorder()

		h.CreateTask(res, req)
		var task models.Task
		err := json.Unmarshal(res.Body.Bytes(), &task)

		assert.Equal(t, http.StatusCreated, res.Code)
		assert.GreaterOrEqual(t, res.Body.Len(), 1)
		assert.Equal(t, "application/json", res.Result().Header["Content-Type"][0])
		assert.NoError(t, err)
	})

	t.Run("Create task with id", func(t *testing.T) {
		h := &BaseHandler{&m}
		req, _ := http.NewRequest("POST", "/task", strings.NewReader(`{"id": 100, "title":"new title","body":"new body","priority":0,"status":"TODO"}`))
		res := httptest.NewRecorder()

		h.CreateTask(res, req)
		var task models.Task
		err := json.Unmarshal(res.Body.Bytes(), &task)

		assert.Equal(t, http.StatusCreated, res.Code)
		assert.GreaterOrEqual(t, res.Body.Len(), 1)
		assert.Equal(t, "application/json", res.Result().Header["Content-Type"][0])
		assert.NoError(t, err)
	})

	t.Run("Create invalid task", func(t *testing.T) {
		h := &BaseHandler{&m}
		req, _ := http.NewRequest("POST", "/task", strings.NewReader(`{"invalid": 121, "eded": eded}`))
		res := httptest.NewRecorder()

		h.CreateTask(res, req)

		assert.NotEqual(t, http.StatusCreated, res.Code)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.GreaterOrEqual(t, res.Body.Len(), 1)
	})
}

/*func TestGetTaskByID(t *testing.T) {
	t.Run("Get a task with existing id", func(t *testing.T) {
		h := &BaseHandler{&m}

		id := uint64(1)
		req, _ := http.NewRequest("GET", fmt.Sprintf("/task/%d", 1), nil)
		res := httptest.NewRecorder()

		h.GetTaskByID(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.GreaterOrEqual(t, res.Body.Len(), 1)
		assert.Equal(t, []byte(fmt.Sprintf(`[{"id":%d,"title":"","body":"","priority":0,"status":""}]`, id)), res.Body.Bytes())
		assert.Equal(t, "application/json", res.Result().Header["Content-Type"][0])
	})
}*/
