package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	DatabaseURL string
	Port        string
}

func Load() *Config {
	// Load .env file if it exists
	godotenv.Load()

	// Build database URL from individual components
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "fmp_db")
	dbUser := getEnv("DB_USER", "user")
	dbPassword := getEnv("DB_PASSWORD", "password")

	databaseURL := getEnv("DATABASE_URL", "")
	if databaseURL == "" {
		databaseURL = "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=disable"
	}

	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		DatabaseURL: databaseURL,
		Port:        getEnv("API_PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
