package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds the application configuration
type Config struct {
	Server         ServerConfig
	Services       ServicesConfig
	Logging        LoggingConfig
	JWT            JWTConfig
	CORS           CORSConfig
	RateLimiting   RateLimitingConfig
	CircuitBreaker CircuitBreakerConfig
}

// RateLimitingConfig holds rate limiting configuration
type RateLimitingConfig struct {
	Enabled      bool
	Limit        int           // Requests per time window
	Burst        int           // Maximum burst size
	Window       time.Duration // Time window for rate limiting
	StoreType    string        // "memory" or "redis"
	RedisAddress string        // Redis address if using Redis
}

// CircuitBreakerConfig holds circuit breaker configuration
type CircuitBreakerConfig struct {
	Enabled                  bool
	Timeout                  time.Duration // How long to wait before trying again
	MaxRequests              uint32        // Max number of requests allowed to half-open state
	RequestVolumeThreshold   uint32        // Minimum requests needed before tripping
	ErrorThresholdPercentage int           // Error percentage to trip circuit
	SleepWindow              time.Duration // How long to wait before testing the service again
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

// CORSConfig holds CORS-related configuration
type CORSConfig struct {
	Enabled          bool
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           time.Duration
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
		CORS: CORSConfig{
			Enabled:      getBoolEnv("CORS_ENABLED", true),
			AllowOrigins: getStringSliceEnv("CORS_ALLOW_ORIGINS", []string{"*"}),
			AllowMethods: getStringSliceEnv("CORS_ALLOW_METHODS", []string{
				"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
			}),
			AllowHeaders: getStringSliceEnv("CORS_ALLOW_HEADERS", []string{
				"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID",
			}),
			ExposeHeaders: getStringSliceEnv("CORS_EXPOSE_HEADERS", []string{
				"Content-Length", "X-Request-ID",
			}),
			AllowCredentials: getBoolEnv("CORS_ALLOW_CREDENTIALS", true),
			MaxAge:           getDurationEnv("CORS_MAX_AGE", 12*time.Hour),
		},
		RateLimiting: RateLimitingConfig{
			Enabled:      getBoolEnv("RATE_LIMIT_ENABLED", true),
			Limit:        getIntEnv("RATE_LIMIT", 100),
			Burst:        getIntEnv("RATE_LIMIT_BURST", 150),
			Window:       getDurationEnv("RATE_LIMIT_WINDOW", time.Minute),
			StoreType:    getEnv("RATE_LIMIT_STORE", "memory"),
			RedisAddress: getEnv("REDIS_ADDRESS", "redis:6379"),
		},

		CircuitBreaker: CircuitBreakerConfig{
			Enabled:                  getBoolEnv("CIRCUIT_BREAKER_ENABLED", true),
			Timeout:                  getDurationEnv("CIRCUIT_BREAKER_TIMEOUT", 30*time.Second),
			MaxRequests:              uint32(getIntEnv("CIRCUIT_BREAKER_MAX_REQUESTS", 5)),
			RequestVolumeThreshold:   uint32(getIntEnv("CIRCUIT_BREAKER_REQUEST_VOLUME", 10)),
			ErrorThresholdPercentage: getIntEnv("CIRCUIT_BREAKER_ERROR_THRESHOLD", 50),
			SleepWindow:              getDurationEnv("CIRCUIT_BREAKER_SLEEP_WINDOW", 10*time.Second),
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

// Add a helper function for string slice environment variables
func getStringSliceEnv(key string, fallback []string) []string {
	value := os.Getenv(key)
	if len(strings.TrimSpace(value)) == 0 {
		return fallback
	}

	return strings.Split(value, ",")
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
