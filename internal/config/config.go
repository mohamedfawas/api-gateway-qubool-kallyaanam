package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds the application configuration
type Config struct {
	Server   ServerConfig
	Services ServicesConfig
	Logging  LoggingConfig
	JWT      JWTConfig
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	Secret           string
	ExpirationHours  int
	RefreshSecret    string
	RefreshExpHours  int
	SigningAlgorithm string
	Issuer           string
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// ServicesConfig holds the URLs for downstream services
type ServicesConfig struct {
	AuthServiceURL  string
	UserServiceURL  string
	AdminServiceURL string
}

// LoggingConfig holds logging-related configuration
type LoggingConfig struct {
	Level       string
	Development bool
}

// NewConfig creates and initializes a new Config instance
func NewConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			ReadTimeout:  getDurationEnv("READ_TIMEOUT", 5*time.Second),
			WriteTimeout: getDurationEnv("WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:  getDurationEnv("IDLE_TIMEOUT", 120*time.Second),
		},
		Services: ServicesConfig{
			AuthServiceURL:  getEnv("AUTH_SERVICE_URL", "http://auth-service:8081"),
			UserServiceURL:  getEnv("USER_SERVICE_URL", "http://user-service:8082"),
			AdminServiceURL: getEnv("ADMIN_SERVICE_URL", "http://admin-service:8083"),
		},
		Logging: LoggingConfig{
			Level:       getEnv("LOG_LEVEL", "info"),
			Development: getBoolEnv("DEVELOPMENT", true),
		},
		JWT: JWTConfig{
			Secret:           getEnv("JWT_SECRET", "your-secret-key"),
			ExpirationHours:  getIntEnv("JWT_EXPIRATION_HOURS", 24),
			RefreshSecret:    getEnv("JWT_REFRESH_SECRET", "your-refresh-secret-key"),
			RefreshExpHours:  getIntEnv("JWT_REFRESH_EXPIRATION_HOURS", 168), // 7 days
			SigningAlgorithm: getEnv("JWT_SIGNING_ALGORITHM", "HS256"),
			Issuer:           getEnv("JWT_ISSUER", "qubool-kallyaanam-api"),
		},
	}
}

// Helper function to get an environment variable with a fallback value
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(strings.TrimSpace(value)) == 0 {
		return fallback
	}
	return value
}

func getIntEnv(key string, fallback int) int {
	value := os.Getenv(key)
	if len(strings.TrimSpace(value)) == 0 {
		return fallback
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return intValue
}

// Helper function to get a duration environment variable with a fallback value
func getDurationEnv(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if len(strings.TrimSpace(value)) == 0 {
		return fallback
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return duration
}

// Helper function to get a boolean environment variable with a fallback value
func getBoolEnv(key string, fallback bool) bool {
	value := os.Getenv(key)
	if len(strings.TrimSpace(value)) == 0 {
		return fallback
	}

	return strings.ToLower(value) == "true"
}
