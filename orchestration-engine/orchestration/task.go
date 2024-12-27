package orchestration

import (
	"log"
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
	ID             string
	Type           string
	Payload        interface{}
	Result         string
	State          TaskState
	Retries        int
	ReasoningSteps []string // For CoT
	DynamicActions []string // For ReAct
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Update Result
func (t *Task) UpdateResult(result string) {
	t.Result = result
	t.UpdatedAt = time.Now()
}

func (t *Task) AddReasoningStep(step string) {
	t.ReasoningSteps = append(t.ReasoningSteps, step)
	t.UpdatedAt = time.Now()
}

func (t *Task) AddDynamicAction(action string) {
	t.DynamicActions = append(t.DynamicActions, action)
	t.UpdatedAt = time.Now()
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

func ProcessTask(task *Task) {
	task.MarkInProgress()
	log.Printf("Processing task: %s, Type: %s\n", task.ID, task.Type)

	// Simulate processing delay
	time.Sleep(500 * time.Millisecond)

	// Simulate task completion or failure based on type (for demo purposes)
	switch task.Type {
	case "Summarization":
		task.Payload = "Summarized content for: " + task.Payload.(string)
	case "Categorization":
		task.Payload = "Categorized content for: " + task.Payload.(string)
	case "Sentiment":
		task.Payload = "Sentiment analysis result for: " + task.Payload.(string)
	default:
		task.MarkFailed()
		log.Printf("Task %s failed: Unknown type\n", task.ID)
		return
	}

	task.MarkCompleted()
	log.Printf("Task %s completed. Result: %v\n", task.ID, task.Payload)
}
