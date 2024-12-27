package orchestration

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// OrchestrationEngine represents the core of the orchestration system
type OrchestrationEngine struct {
	queue      *TaskQueue  // Task queue
	workerPool *WorkerPool // Worker pool
	isRunning  bool        // Indicates if the engine is running
	mu         sync.Mutex  // Mutex for thread-safe access
}

// NewOrchestrationEngine creates a new orchestration engine
func NewOrchestrationEngine(queueCapacity, workerCount int) *OrchestrationEngine {
	queue := NewTaskQueue(queueCapacity)
	workerPool := NewWorkerPool(queue, workerCount)

	return &OrchestrationEngine{
		queue:      queue,
		workerPool: workerPool,
		isRunning:  false,
	}
}

// Start starts the orchestration engine
func (oe *OrchestrationEngine) Start() {
	oe.mu.Lock()
	defer oe.mu.Unlock()

	if oe.isRunning {
		log.Println("Orchestration engine is already running.")
		return
	}

	log.Println("Starting orchestration engine...")
	oe.isRunning = true
	oe.workerPool.Start()
	log.Println("Orchestration engine started.")
}

// Stop stops the orchestration engine
func (oe *OrchestrationEngine) Stop() {
	oe.mu.Lock()
	defer oe.mu.Unlock()

	if !oe.isRunning {
		log.Println("Orchestration engine is not running.")
		return
	}

	log.Println("Stopping orchestration engine...")
	oe.workerPool.Stop()
	oe.queue.Clear()
	oe.isRunning = false
	log.Println("Orchestration engine stopped.")
}

// AddTask adds a new task to the queue
func (oe *OrchestrationEngine) AddTask(task *Task) error {
	oe.mu.Lock()
	defer oe.mu.Unlock()

	if !oe.isRunning {
		log.Println("Cannot add task. Orchestration engine is not running.")
		return NewEngineError("Orchestration engine is not running")
	}

	log.Printf("Adding task %s to the queue. Current queue size: %d\n", task.ID, oe.queue.Size())
	err := oe.queue.Enqueue(task)
	if err != nil {
		log.Printf("Failed to enqueue task %s: %v\n", task.ID, err)
		return err
	}
	log.Printf("Task %s added to the queue. New queue size: %d\n", task.ID, oe.queue.Size())
	return nil
}

// QueueSize returns the current size of the task queue
func (oe *OrchestrationEngine) QueueSize() int {
	oe.mu.Lock()
	defer oe.mu.Unlock()
	return oe.queue.Size()
}

// Run executes a batch of tasks and waits for them to complete
func (oe *OrchestrationEngine) Run(tasks []*Task) ([]string, error) {
	oe.Start()

	for _, task := range tasks {
		if err := oe.AddTask(task); err != nil {
			return nil, fmt.Errorf("failed to add task %s: %w", task.ID, err)
		}
	}

	// Wait for all tasks to complete
	for oe.QueueSize() > 0 {
		log.Printf("Waiting for queue to empty... Current size: %d\n", oe.QueueSize())
		time.Sleep(100 * time.Millisecond)
	}

	oe.Stop()

	// Collect results
	var results []string
	for _, task := range tasks {
		if task.State == Completed {
			results = append(results, task.Result)
		} else {
			log.Printf("Task %s did not complete successfully. State: %s", task.ID, task.State)
		}
	}

	return results, nil
}
