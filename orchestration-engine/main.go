package main

import (
	"email_pilot/config"
	"email_pilot/controllers"
	"email_pilot/orchestration"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := config.InitializeDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	config.MigrateDatabase()

	psqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB from GORM: %v", err)
	}
	defer psqlDB.Close()

	// Initialize the orchestration engine
	engine := orchestration.NewOrchestrationEngine(10, 3)

	router := gin.Default()

	router.GET("/auth", controllers.AuthController)
	router.GET("/auth/callback", controllers.CallbackController)
	router.GET("/emails", controllers.GetEmailsController)

	// Simulate orchestration using sample data
	router.POST("/simulate", func(c *gin.Context) {
		controllers.SimulateOrchestrationController(c, engine)
	})

	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
