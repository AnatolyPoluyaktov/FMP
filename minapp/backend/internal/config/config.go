package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment      string
	Port             string
	TelegramBotToken string
	FMPCoreAPIURL    string
	FrontendURL      string
}

func Load() *Config {
	// Load .env file if it exists
	godotenv.Load()

	return &Config{
		Environment:      getEnv("ENVIRONMENT", "development"),
		Port:             getEnv("PORT", "8080"),
		TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		FMPCoreAPIURL:    getEnv("FMP_CORE_API_URL", "http://localhost:8080/api/v1"),
		FrontendURL:      getEnv("FRONTEND_URL", "http://localhost:3000"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
