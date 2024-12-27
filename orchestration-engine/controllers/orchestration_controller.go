package controllers

import (
	"email_pilot/orchestration"
	"fmt"
	"log"
	"net/http"

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
func SimulateOrchestrationController(c *gin.Context, engine *orchestration.OrchestrationEngine) {
	// Start the orchestration engine
	engine.Start()
	defer engine.Stop()

	// Add tasks to the engine
	for _, email := range SampleEmails {
		taskID := fmt.Sprintf("task-%s", email.ID)
		task := orchestration.NewTask(taskID, "EmailProcessing", email)
		if err := engine.AddTask(task); err != nil {
			log.Printf("Failed to add task %s: %v\n", taskID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add task"})
			return
		}
	}

	// Monitor the engine's queue
	engine.Monitor()

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Orchestration simulation started"})
}
