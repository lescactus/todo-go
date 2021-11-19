package databases

import (
	"fmt"
	"sync"
	"todo-go/models"
)

type InMemoryDatabase struct {
	Tasks []models.Task
	rwm   sync.RWMutex
}

func NewInMemoryDatabase() *InMemoryDatabase {
	return &InMemoryDatabase{
		Tasks: make([]models.Task, 0),
	}
}

func (db *InMemoryDatabase) GetTaskByID(id uint64) (*models.Task, error) {
	db.rwm.RLock()
	defer db.rwm.RUnlock()

	for k, v := range db.Tasks {
		if v.Id == id {
			return &db.Tasks[k], nil
		}
	}

	return &models.Task{}, fmt.Errorf("no task with id %v exists", id)
}

func (db *InMemoryDatabase) GetAllTasks() []models.Task {
	db.rwm.RLock()
	defer db.rwm.RUnlock()

	return db.Tasks
}

func (db *InMemoryDatabase) CreateTask(t models.Task) (uint64, error) {
	db.rwm.Lock()
	defer db.rwm.Unlock()

	var id uint64
	length := len(db.Tasks)
	// First id is 1
	if length == 0 {
		id = 0
	} else {
		id = uint64(length) + 1
	}

	t.Id = uint64(id)
	db.Tasks = append(db.Tasks, t)

	return uint64(id), nil
}

func (db *InMemoryDatabase) UpdateTask(t models.Task) error {
	task, err := db.GetTaskByID(t.Id)
	if err != nil {
		return fmt.Errorf("error updating task id %d: %s", t.Id, err.Error())
	}

	db.rwm.Lock()
	defer db.rwm.Unlock()

	task.Title = t.Title
	task.Body = t.Body
	task.Priority = t.Priority
	task.Status = t.Status

	return nil
}

func (db *InMemoryDatabase) DeleteTask(id uint64) error {
	_, err := db.GetTaskByID(id)
	if err != nil {
		return fmt.Errorf("error deleting task id %d: %s", id, err.Error())
	}

	for k, v := range db.Tasks {
		if v.Id == id {
			// Remove the element at index k
			// https://github.com/golang/go/wiki/SliceTricks#delete
			//s.Tasks = append(s.Tasks[:k], s.Tasks[k+1:]...)
			empty := new(models.Task)
			copy(db.Tasks[k:], db.Tasks[k+1:])
			db.Tasks[len(db.Tasks)-1] = *empty
			db.Tasks = db.Tasks[:len(db.Tasks)-1]
		}
	}

	return nil
}
