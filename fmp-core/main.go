package main

import (
	"log"
	"os"

	"fmp-core/internal/api"
	"fmp-core/internal/config"
	"fmp-core/internal/database"
	"fmp-core/internal/migrations"
	"fmp-core/internal/services"

	"github.com/gin-gonic/gin"
)

// @title FMP Core API
// @version 1.0
// @description Financial Manager Platform Core API
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load configuration
	cfg := config.Load()
	log.Printf("Configuration loaded: Environment=%s, Port=%s", cfg.Environment, cfg.Port)
	log.Printf("Database URL: %s", cfg.DatabaseURL)

	// Initialize database
	log.Println("Initializing database connection...")
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()
	log.Println("Database connection established successfully")

	// Run migrations
	log.Println("Running database migrations...")
	if err := migrations.Run(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
	log.Println("Database migrations completed successfully")

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
		log.Println("Gin mode set to production")
	} else {
		log.Println("Gin mode set to development")
	}

	// Initialize router
	router := gin.Default()

	// Initialize services with database
	services.SetDB(db)
	log.Println("Services initialized with database")

	// Setup API routes
	api.SetupRoutes(router, db)
	log.Println("API routes configured")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
