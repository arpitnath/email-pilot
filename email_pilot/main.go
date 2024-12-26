package main

import (
	"email_pilot/config"
	"email_pilot/controllers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load application configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := config.InitializeDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Run database migrations
	config.MigrateDatabase()

	// Extract sql.DB from GORM and defer its close
	psqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB from GORM: %v", err)
	}
	defer psqlDB.Close()

	// Initialize Gin router
	router := gin.Default()

	router.GET("/auth", controllers.AuthController)
	router.GET("/auth/callback", controllers.CallbackController)
	router.GET("/emails", controllers.GetEmailsController)

	// Start the server
	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
