package controllers

import (
	"email_pilot/orchestration"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SampleEmails defines mock data for email simulation
var SampleEmails = []struct {
	ID      string
	Subject string
	Body    string
}{
	{"1", "Your invoice for October", "Dear user, your invoice is $45. Please pay by the due date."},
	{"2", "Congratulations!", "We are thrilled to announce you got the job!"},
	{"3", "Limited Offer!", "Get 50% off on all electronics until midnight."},
}

// SimulateOrchestrationController simulates orchestration using sample email data
func SimulateOrchestrationController(c *gin.Context) {
	engine := orchestration.NewOrchestrationEngine(10, 3) // New engine instance

	// Start the orchestration engine
	engine.Start()

	// Add tasks to the engine
	for _, email := range SampleEmails {
		taskID := fmt.Sprintf("task-%s", email.ID)
		task := orchestration.NewTask(taskID, "Summarization", email.Body)
		if err := engine.AddTask(task); err != nil {
			log.Printf("Failed to add task %s: %v\n", taskID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add task"})
			return
		}
	}

	// Wait for the queue to be processed
	for {
		queueSize := engine.QueueSize()
		log.Printf("Waiting for queue to empty... Current size: %d\n", queueSize)
		if queueSize == 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}

	// Stop the engine after processing
	engine.Stop()

	// Collect results
	results := []string{}
	for _, email := range SampleEmails {
		results = append(results, fmt.Sprintf("Task %s for email %s completed successfully.", email.ID, email.Subject))
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Orchestration simulation completed successfully.",
		"results": results,
	})
}
