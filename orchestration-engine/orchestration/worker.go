package orchestration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	task.MarkInProgress()

	// Call the appropriate handler based on task type
	var err error
	switch task.Type {
	case "Summarization":
		err = handleSummarization(task)
	case "Categorization":
		err = handleCategorization(task)
	case "Sentiment Analysis":
		err = handleSentimentAnalysis(task)
	default:
		err = fmt.Errorf("unknown task type: %s", task.Type)
	}

	if err != nil {
		log.Printf("Worker %d failed to process task %s: %v\n", workerID, task.ID, err)
		task.MarkFailed()
		task.IncrementRetries()

		if task.Retries < 3 {
			log.Printf("Retrying task %s (retry #%d)\n", task.ID, task.Retries)
			wp.queue.Enqueue(task)
		} else {
			log.Printf("Task %s failed after maximum retries.\n", task.ID)
		}
		return
	}

	task.MarkCompleted()
	log.Printf("Worker %d completed task: %s\n", workerID, task.ID)
}

func handleSummarization(task *Task) error {
	payload := map[string]string{"prompt": task.Payload.(string)}
	return callLLMService("http://localhost:8000/api/summarize", payload)
}

func handleCategorization(task *Task) error {
	payload := map[string]string{"prompt": task.Payload.(string)}
	return callLLMService("http://localhost:8000/api/categorize", payload)
}

func handleSentimentAnalysis(task *Task) error {
	payload := map[string]string{"prompt": task.Payload.(string)}
	return callLLMService("http://localhost:8000/api/sentiment", payload)
}

// callLLMService sends a request to the LLM service and handles the response
func callLLMService(url string, payload map[string]string) error {
	// Marshal the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to call LLM service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("LLM service returned status: %d", resp.StatusCode)
	}

	log.Printf("LLM service call to %s succeeded.\n", url)
	return nil
}
