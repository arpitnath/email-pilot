package orchestration

import (
	"log"
)

type OrchestrationEngine struct {
	queue      *TaskQueue
	workerPool *WorkerPool
	isRunning  bool
}

func NewOrchestrationEngine(queueCapacity, workerCount int) *OrchestrationEngine {
	// Initialize the task queue and worker pool
	queue := NewTaskQueue(queueCapacity)
	workerPool := NewWorkerPool(queue, workerCount)

	return &OrchestrationEngine{
		queue:      queue,
		workerPool: workerPool,
		isRunning:  false,
	}
}

func (oe *OrchestrationEngine) Start() {
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
	if !oe.isRunning {
		log.Println("Orchestration engine is not running.")
		return
	}

	log.Println("Stopping orchestration engine...")
	oe.workerPool.Stop()
	oe.isRunning = false
	log.Println("Orchestration engine stopped.")
}

// AddTask adds a new task to the queue
func (oe *OrchestrationEngine) AddTask(task *Task) error {
	if !oe.isRunning {
		return &EngineError{Message: "Cannot add tasks. Orchestration engine is not running."}
	}

	err := oe.queue.Enqueue(task)
	if err != nil {
		log.Printf("Failed to enqueue task %s: %v\n", task.ID, err)
		return err
	}

	log.Printf("Task %s added to the queue.\n", task.ID)
	return nil
}

func (oe *OrchestrationEngine) Monitor() {
	//Todo: Implement monitoring logic
	log.Printf("Queue size: %d\n", oe.queue.Size())
}
