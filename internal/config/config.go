package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Port string
	}
	Services map[string]string
	Auth     struct {
		JWTSecret string
	}
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		// Just log the error but continue - it's not fatal if .env doesn't exist
		log.Println("Warning: .env file not found. Using environment variables only.")
	}

	config := &Config{
		Services: make(map[string]string),
	}

	// Server config
	config.Server.Port = getEnv("PORT", "8080")

	// Service endpoints - dynamically load any service with SERVICE_ prefix
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "SERVICE_") {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				serviceName := strings.ToLower(strings.TrimPrefix(parts[0], "SERVICE_"))
				config.Services[serviceName] = parts[1]
			}
		}
	}

	// Auth config
	config.Auth.JWTSecret = getEnv("JWT_SECRET", "")

	return config
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
