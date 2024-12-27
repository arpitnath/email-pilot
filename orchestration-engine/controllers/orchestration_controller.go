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
	engine := orchestration.NewOrchestrationEngine(10, 3)
	engine.Start()

	var results []string // Collect results from processed tasks

	// Add tasks to the engine
	for _, email := range SampleEmails {
		// Summarization task
		taskID := fmt.Sprintf("summarize-%s", email.ID)
		task := orchestration.NewTask(taskID, "Summarization", email.Body)
		engine.AddTask(task)
		results = append(results, fmt.Sprintf("Summarization task created for email: %s", email.Subject))

		// Categorization task
		taskID = fmt.Sprintf("categorize-%s", email.ID)
		task = orchestration.NewTask(taskID, "Categorization", email.Body)
		engine.AddTask(task)
		results = append(results, fmt.Sprintf("Categorization task created for email: %s", email.Subject))

		// Sentiment analysis task
		taskID = fmt.Sprintf("sentiment-%s", email.ID)
		task = orchestration.NewTask(taskID, "Sentiment", email.Body)
		engine.AddTask(task)
		results = append(results, fmt.Sprintf("Sentiment task created for email: %s", email.Subject))
	}

	// Wait for all tasks to be processed
	for engine.QueueSize() > 0 {
		log.Printf("Waiting for queue to empty... Current size: %d\n", engine.QueueSize())
		time.Sleep(500 * time.Millisecond)
	}

	engine.Stop()

	// Return the results
	c.JSON(http.StatusOK, gin.H{
		"message": "Orchestration simulation completed successfully.",
		"results": results,
	})
}
