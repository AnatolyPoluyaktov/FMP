package main

import (
	"log"
	"os"

	"minapp-backend/internal/api"
	"minapp-backend/internal/config"
	"minapp-backend/internal/telegram"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Telegram bot
	bot, err := telegram.InitializeBot(cfg.TelegramBotToken)
	if err != nil {
		log.Fatal("Failed to initialize Telegram bot:", err)
	}

	// Initialize router
	router := gin.Default()

	// Setup API routes
	api.SetupRoutes(router, bot, cfg)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting minapp backend server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
