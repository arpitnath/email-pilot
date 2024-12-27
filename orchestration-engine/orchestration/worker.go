package orchestration

import (
	"log"
	"sync"
)

// WorkerPool manages a pool of workers to process tasks concurrently
type WorkerPool struct {
	queue       *TaskQueue
	workerCount int
	wg          sync.WaitGroup
	stopCh      chan struct{}
}

func NewWorkerPool(queue *TaskQueue, workerCount int) *WorkerPool {
	return &WorkerPool{
		queue:       queue,
		workerCount: workerCount,
		stopCh:      make(chan struct{}),
	}
}

func (wp *WorkerPool) Start() {
	log.Printf("Starting worker pool with %d workers...\n", wp.workerCount)
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

func (wp *WorkerPool) Stop() {
	log.Println("Stopping worker pool...")
	close(wp.stopCh)
	wp.wg.Wait()
	log.Println("Worker pool stopped.")
}

// worker is a single worker goroutine that processes tasks
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	log.Printf("Worker %d started.\n", id)

	for {
		select {
		case <-wp.stopCh: // Stop signal received
			log.Printf("Worker %d stopping...\n", id)
			return
		default:
			task, err := wp.queue.Dequeue()
			if err != nil {
				// Queue is empty; let the worker wait for a task
				continue
			}

			wp.processTask(id, task)
		}
	}
}

func (wp *WorkerPool) processTask(workerID int, task *Task) {
	log.Printf("Worker %d processing task: %s\n", workerID, task.ID)

	//Todo: Simulate task processing (replace with actual task logic)
	task.MarkInProgress()

	//Todo: Simulate success or failure
	task.MarkCompleted()
	log.Printf("Worker %d completed task: %s\n", workerID, task.ID)
}
