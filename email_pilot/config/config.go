package config

import (
	"email_pilot/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// LoadConfig loads environment variables from a .env file
func LoadConfig() error {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Using environment variables.")
	}
	return nil
}

// InitializeDatabase sets up the database connection
func InitializeDatabase() (*gorm.DB, error) {
	// Read DB credentials from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	// PostgreSQL DSN
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	// Initialize GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Assign the database instance to the global variable
	DB = db

	return db, nil
}

// MigrateDatabase ensures the database schema is up-to-date
func MigrateDatabase() {
	err := DB.AutoMigrate(&models.OAuthToken{})
	if err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}
	log.Println("Database migration completed successfully")
}
