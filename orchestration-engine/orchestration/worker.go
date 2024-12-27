package orchestration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// WorkerPool manages a pool of workers to process tasks concurrently
type WorkerPool struct {
	queue       *TaskQueue
	workerCount int
	wg          sync.WaitGroup
	stopChan    chan struct{}
}

func NewWorkerPool(queue *TaskQueue, workerCount int) *WorkerPool {
	return &WorkerPool{
		queue:       queue,
		workerCount: workerCount,
		stopChan:    make(chan struct{}),
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
	close(wp.stopChan)
	wp.wg.Wait()
	log.Println("Worker pool stopped.")
}

// worker is a single worker goroutine that processes tasks
func (wp *WorkerPool) worker(workerID int) {
	log.Printf("Worker %d started.\n", workerID)

	for {
		select {
		case <-wp.stopChan:
			log.Printf("Worker %d stopping...\n", workerID)
			wp.wg.Done() // Signal the WaitGroup
			return
		default:
			task, err := wp.queue.Dequeue()
			if err != nil {
				// Queue is empty; wait before retrying
				time.Sleep(100 * time.Millisecond)
				continue
			}

			wp.processTask(workerID, task) // Call processTask
		}
	}
}

func (wp *WorkerPool) processTask(workerID int, task *Task) {
	log.Printf("Worker %d processing task: %s\n", workerID, task.ID)
	task.MarkInProgress()

	var err error
	switch task.Type {
	case "Summarization":
		err = handleSummarization(task)
	case "Categorization":
		err = handleCategorization(task)
	case "Sentiment":
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

	// Add dynamic tasks for ReAct
	for _, action := range task.DynamicActions {
		log.Printf("Worker %d handling dynamic action: %s\n", workerID, action)
		// Example: Enqueue a new task
		if action == "Enqueue Notification Task" {
			newTask := NewTask(fmt.Sprintf("notification-%s", task.ID), "Notification", task.Payload)
			wp.queue.Enqueue(newTask)
		}
	}

	task.MarkCompleted()
	log.Printf("Worker %d completed task: %s\n", workerID, task.ID)
}

func handleSummarization(task *Task) error {
	task.AddReasoningStep("Started summarization.")
	payload := map[string]string{"prompt": task.Payload.(string)}

	if err := callLLMService("http://localhost:8000/api/summarize", payload); err != nil {
		task.AddReasoningStep("Failed summarization.")
		return err
	}

	task.AddReasoningStep("Completed summarization.")
	return nil
}

func handleCategorization(task *Task) error {
	task.AddReasoningStep("Started categorization.")
	payload := map[string]string{"prompt": task.Payload.(string)}

	if err := callLLMService("http://localhost:8000/api/categorize", payload); err != nil {
		task.AddReasoningStep("Failed categorization.")
		return err
	}

	task.AddReasoningStep("Completed categorization.")
	// Dynamically enqueue a notification if categorized as "High Priority"
	if isHighPriority(task) {
		task.AddDynamicAction("Enqueue Notification Task")
	}
	return nil
}

func handleSentimentAnalysis(task *Task) error {
	task.AddReasoningStep("Started sentiment analysis.")
	payload := map[string]string{"prompt": task.Payload.(string)}

	if err := callLLMService("http://localhost:8000/api/sentiment", payload); err != nil {
		task.AddReasoningStep("Failed sentiment analysis.")
		return err
	}

	task.AddReasoningStep("Completed sentiment analysis.")
	return nil
}

func isHighPriority(task *Task) bool {
	// Example logic to determine high priority
	return strings.Contains(task.Payload.(string), "urgent")
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
