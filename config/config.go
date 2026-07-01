package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	MangaDexBaseURL   string
	MangaDexRateLimit int
}

var AppConfig Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	AppConfig.Port = getEnv("PORT", "8080")
	AppConfig.MangaDexBaseURL = getEnv("MANGADEX_BASE_URL", "https://api.mangadex.org")
	AppConfig.MangaDexRateLimit = getEnvAsInt("MANGADEX_RATE_LIMIT", 5)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}