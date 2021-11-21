package databases

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
	"todo-go/models"

	"github.com/stretchr/testify/assert"
)

var (
	inMemoryDatabase = &InMemoryDatabase{
		Tasks: []models.Task{
			{
				Id:       0,
				Title:    "Test Title 0",
				Body:     "Test body 0",
				Priority: models.Highest,
				Status:   models.StatusInProgress,
			},
			{
				Id:       1,
				Title:    "Test Title 0",
				Body:     "Test body 0",
				Priority: models.Low,
				Status:   models.StatusDone,
			},
		},
	}
)

func TestNewInMemoryDatabase(t *testing.T) {
	t.Run("Create a InMemoryDatabase", func(t *testing.T) {
		db := NewInMemoryDatabase()
		assert.NotEmpty(t, db)
	})
}

func TestGetTaskByID(t *testing.T) {
	for k := range inMemoryDatabase.Tasks {
		t.Run(fmt.Sprintf("Id %d exists", k), func(t *testing.T) {
			id := uint64(k)
			task, err := inMemoryDatabase.GetTaskByID(id)
			assert.NoError(t, err)
			assert.NotEmpty(t, task)
			assert.Equal(t, inMemoryDatabase.Tasks[id].Id, id)
			assert.Equal(t, inMemoryDatabase.Tasks[id].Title, task.Title)
			assert.Equal(t, inMemoryDatabase.Tasks[id].Body, task.Body)
			assert.Equal(t, inMemoryDatabase.Tasks[id].Priority, task.Priority)
			assert.Equal(t, inMemoryDatabase.Tasks[id].Status, task.Status)
		})
	}

	t.Run("Id doesn't exist", func(t *testing.T) {
		id := uint64(99999)
		task, err := inMemoryDatabase.GetTaskByID(id)
		assert.Error(t, err)
		assert.Empty(t, task)
	})
}

func TestGetAllTasks(t *testing.T) {
	t.Run("Get all tasks", func(t *testing.T) {
		tasks := inMemoryDatabase.GetAllTasks()
		assert.NotEmpty(t, tasks)
		assert.Equal(t, len(tasks), len(inMemoryDatabase.Tasks))
		assert.Equal(t, tasks, inMemoryDatabase.Tasks)
	})
}

func TestStoreCreateTask(t *testing.T) {
	t.Run("Add a task without id", func(t *testing.T) {
		newId := uint64(2)
		id, err := inMemoryDatabase.CreateTask(models.Task{
			Title:    fmt.Sprintf("Test Title %d", newId),
			Body:     fmt.Sprintf("Test body %d", newId),
			Priority: models.Highest,
			Status:   models.StatusInProgress,
		})
		assert.NoError(t, err)
		assert.Equal(t, id, newId)
		assert.Equal(t, inMemoryDatabase.Tasks[id].CreatedAt, inMemoryDatabase.Tasks[id].UpdatedAt)
	})

	t.Run("Add a task with non existent id", func(t *testing.T) {
		newId := uint64(3)
		id, err := inMemoryDatabase.CreateTask(models.Task{
			Id:       newId,
			Title:    fmt.Sprintf("Test Title %d", newId),
			Body:     fmt.Sprintf("Test body %d", newId),
			Priority: models.Highest,
			Status:   models.StatusInProgress,
		})
		assert.NoError(t, err)
		assert.Equal(t, id, newId)
	})

	t.Run("Add a task with existent id", func(t *testing.T) {
		newId := uint64(3)
		id, err := inMemoryDatabase.CreateTask(models.Task{
			Id:       newId,
			Title:    fmt.Sprintf("Test Title %d", newId),
			Body:     fmt.Sprintf("Test body %d", newId),
			Priority: models.Highest,
			Status:   models.StatusInProgress,
		})
		assert.NoError(t, err)
		assert.Equal(t, id, newId+1)
	})

	t.Run("Add tasks with an empty InMemoryDatabase", func(t *testing.T) {
		db := &InMemoryDatabase{Tasks: []models.Task{}}
		id, err := db.CreateTask(models.Task{
			Title:    "Test Title 1",
			Body:     "Test body",
			Priority: models.Highest,
			Status:   models.StatusInProgress,
		})
		assert.NoError(t, err)
		assert.Equal(t, uint64(0), id)

		id, err = db.CreateTask(models.Task{
			Title:    "Test Title 2",
			Body:     "Test body",
			Priority: models.Highest,
			Status:   models.StatusInProgress,
		})
		assert.NoError(t, err)
		assert.Equal(t, uint64(1), id)
	})

	t.Run("Add a task with a field CreatedAt set", func(t *testing.T) {
		randomDate := randomDate()

		id, err := inMemoryDatabase.CreateTask(models.Task{
			CreatedAt: randomDate,
		})
		assert.NoError(t, err)
		assert.NotEqual(t, inMemoryDatabase.Tasks[id].CreatedAt, randomDate)
	})

	t.Run("Add a task with a field UpdatedAt set", func(t *testing.T) {
		randomDate := randomDate()

		id, err := inMemoryDatabase.CreateTask(models.Task{
			CreatedAt: randomDate,
		})
		assert.NoError(t, err)
		assert.NotEqual(t, inMemoryDatabase.Tasks[id].UpdatedAt, randomDate)
	})
}

func TestStoreUpdateTask(t *testing.T) {
	t.Run("Update task with existing id", func(t *testing.T) {
		id := uint64(0)
		task := models.Task{
			Id:       id,
			Title:    "Updated title",
			Body:     "Updated body",
			Priority: models.Highest,
			Status:   models.StatusInProgress,
		}
		err := inMemoryDatabase.UpdateTask(task)
		assert.NoError(t, err)
		assert.Equal(t, uint64(id), inMemoryDatabase.Tasks[id].Id)
		assert.Equal(t, inMemoryDatabase.Tasks[id].Title, task.Title)
		assert.Equal(t, inMemoryDatabase.Tasks[id].Body, task.Body)
		assert.Equal(t, inMemoryDatabase.Tasks[id].Priority, task.Priority)
		assert.Equal(t, inMemoryDatabase.Tasks[id].Status, task.Status)
		assert.NotEqual(t, inMemoryDatabase.Tasks[id].UpdatedAt, inMemoryDatabase.Tasks[id].CreatedAt)
	})

	t.Run("Update task with non existing id", func(t *testing.T) {
		id := uint64(999999999)
		task := models.Task{
			Id:       id,
			Title:    "Updated title",
			Body:     "Updated body",
			Priority: models.Highest,
			Status:   models.StatusInProgress,
		}
		err := inMemoryDatabase.UpdateTask(task)
		assert.Error(t, err)
	})

	t.Run("Update task with field UpdatedAt set", func(t *testing.T) {
		id := uint64(0)
		randomDate := randomDate()
		task := models.Task{
			Id:        id,
			Title:     "Updated title",
			Body:      "Updated body",
			Priority:  models.Highest,
			Status:    models.StatusInProgress,
			UpdatedAt: randomDate,
		}
		err := inMemoryDatabase.UpdateTask(task)
		assert.NoError(t, err)
		assert.Equal(t, uint64(id), inMemoryDatabase.Tasks[id].Id)
		assert.Equal(t, inMemoryDatabase.Tasks[id].Title, task.Title)
		assert.Equal(t, inMemoryDatabase.Tasks[id].Body, task.Body)
		assert.Equal(t, inMemoryDatabase.Tasks[id].Priority, task.Priority)
		assert.Equal(t, inMemoryDatabase.Tasks[id].Status, task.Status)
		assert.NotEqual(t, inMemoryDatabase.Tasks[id].UpdatedAt, inMemoryDatabase.Tasks[id].CreatedAt)
		assert.NotEqual(t, inMemoryDatabase.Tasks[id].UpdatedAt, randomDate)
	})
}

func TestStoreDeleteTask(t *testing.T) {
	t.Run("Delete task with existing id", func(t *testing.T) {
		id := uint64(1)
		err := inMemoryDatabase.DeleteTask(id)
		assert.NoError(t, err)
	})

	t.Run("Delete task with non existing id", func(t *testing.T) {
		id := uint64(999999)
		err := inMemoryDatabase.DeleteTask(id)
		assert.Error(t, err)
	})
}

func randomDate() time.Time {
	// Generate a random between min and max
	// ref: https://stackoverflow.com/questions/43495745/how-to-generate-random-date-in-go-lang/43497333
	min := time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2020, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min
	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

func benchmarkCreateTask(i int, b *testing.B) {
	db := NewInMemoryDatabase()

	for n := 0; n < b.N; n++ {
		db.CreateTask(models.Task{
			Title: fmt.Sprintf("Title %d", n),
		})
	}
}

func BenchmarkCreateTask1(b *testing.B)    { benchmarkCreateTask(1, b) }
func BenchmarkCreateTask2(b *testing.B)    { benchmarkCreateTask(1, b) }
func BenchmarkCreateTask3(b *testing.B)    { benchmarkCreateTask(1, b) }
func BenchmarkCreateTask10(b *testing.B)   { benchmarkCreateTask(1, b) }
func BenchmarkCreateTask100(b *testing.B)  { benchmarkCreateTask(1, b) }
func BenchmarkCreateTask1000(b *testing.B) { benchmarkCreateTask(1, b) }

func benchmarkGetTasks(i int, b *testing.B) {
	db := NewInMemoryDatabase()

	for n := 0; n < 100000; n++ {
		_, _ = db.CreateTask(models.Task{})
	}

	for n := 0; n < b.N; n++ {
		db.GetAllTasks()
	}
}

func BenchmarkGetTasks1(b *testing.B)    { benchmarkGetTasks(1, b) }
func BenchmarkGetTasks2(b *testing.B)    { benchmarkGetTasks(1, b) }
func BenchmarkGetTasks3(b *testing.B)    { benchmarkGetTasks(1, b) }
func BenchmarkGetTasks10(b *testing.B)   { benchmarkGetTasks(1, b) }
func BenchmarkGetTasks100(b *testing.B)  { benchmarkGetTasks(1, b) }
func BenchmarkGetTasks1000(b *testing.B) { benchmarkGetTasks(1, b) }