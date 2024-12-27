package orchestration

import (
	"log"
	"sync"
)

// OrchestrationEngine represents the core of the orchestration system
type OrchestrationEngine struct {
	queue      *TaskQueue     // Task queue
	workerPool *WorkerPool    // Worker pool
	isRunning  bool           // Indicates if the engine is running
	mu         sync.Mutex     // Mutex for thread-safe access
	wg         sync.WaitGroup // WaitGroup for worker completion
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
