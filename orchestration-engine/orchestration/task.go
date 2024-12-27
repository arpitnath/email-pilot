package orchestration

import (
	"time"
)

type TaskState string

const (
	Pending    TaskState = "Pending"
	InProgress TaskState = "In Progress"
	Completed  TaskState = "Completed"
	Failed     TaskState = "Failed"
)

type Task struct {
	ID        string
	Type      string
	Payload   interface{}
	State     TaskState
	Retries   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewTask(id string, taskType string, payload interface{}) *Task {
	return &Task{
		ID:        id,
		Type:      taskType,
		Payload:   payload,
		State:     Pending,
		Retries:   0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (t *Task) MarkInProgress() {
	t.State = InProgress
	t.UpdatedAt = time.Now()
}

func (t *Task) MarkCompleted() {
	t.State = Completed
	t.UpdatedAt = time.Now()
}

func (t *Task) MarkFailed() {
	t.State = Failed
	t.UpdatedAt = time.Now()
}

func (t *Task) IncrementRetries() {
	t.Retries++
	t.UpdatedAt = time.Now()
}
