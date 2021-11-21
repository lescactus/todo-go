package models

import "time"

type Priority int
type Status string

const (
	Lowest Priority = iota
	Low
	Medium
	High
	Highest

	StatusToDo       = "TODO"
	StatusInProgress = "INPROGRESS"
	StatusDone       = "DONE"
)

type Task struct {
	Id       uint64   `json:"id"`
	Title    string   `json:"title"`
	Body     string   `json:"body"`
	Priority Priority `json:"priority"`
	Status   Status   `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TaskRepository interface {
	GetTaskByID(id uint64) (*Task, error)
	GetAllTasks() []Task
	CreateTask(t Task) (uint64, error)
	UpdateTask(t Task) error
	DeleteTask(id uint64) error
}
