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

	var results []map[string]interface{}
	for _, email := range SampleEmails {
		taskID := fmt.Sprintf("summarize-%s", email.ID)
		task := orchestration.NewTask(taskID, "Summarization", email.Body)
		engine.AddTask(task)

		taskID = fmt.Sprintf("categorize-%s", email.ID)
		task = orchestration.NewTask(taskID, "Categorization", email.Body)
		engine.AddTask(task)

		taskID = fmt.Sprintf("sentiment-%s", email.ID)
		task = orchestration.NewTask(taskID, "Sentiment", email.Body)
		engine.AddTask(task)

		results = append(results, map[string]interface{}{
			"email_subject": email.Subject,
			"tasks": []map[string]interface{}{
				{"task_id": taskID, "reasoning_steps": task.ReasoningSteps},
			},
		})
	}

	for engine.QueueSize() > 0 {
		log.Printf("Waiting for queue to empty... Current size: %d\n", engine.QueueSize())
		time.Sleep(500 * time.Millisecond)
	}

	engine.Stop()

	c.JSON(http.StatusOK, gin.H{
		"message": "Orchestration simulation completed successfully.",
		"results": results,
	})
}
