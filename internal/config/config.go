package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	APIKey      string
	DatabaseURL string
	BaseURL     string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	return &Config{
		Port:        getEnv("PORT", "8080"),
		APIKey:      getEnv("API_KEY", ""),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		BaseURL:     getEnv("PRICE_API_URL", "http://localhost:8081"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
