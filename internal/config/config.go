package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	Server       ServerConfig
	Services     ServicesConfig
	Logging      LoggingConfig
	JWT          JWTConfig
	CORS         CORSConfig
	RateLimiting RateLimitingConfig
}

// RateLimitingConfig holds rate limiting configuration
type RateLimitingConfig struct {
	Enabled      bool
	Limit        int           // Requests per time window
	Burst        int           // Maximum burst size
	Window       time.Duration // Time window for rate limiting
	RedisAddress string        // Redis address for rate limiting
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
	setDefaults()

	// Auto-read from environment variables
	viper.AutomaticEnv()

	return &Config{
		Server: ServerConfig{
			Port:         viper.GetString("PORT"),
			ReadTimeout:  viper.GetDuration("READ_TIMEOUT"),
			WriteTimeout: viper.GetDuration("WRITE_TIMEOUT"),
			IdleTimeout:  viper.GetDuration("IDLE_TIMEOUT"),
		},
		Services: ServicesConfig{
			AuthServiceURL:  viper.GetString("AUTH_SERVICE_URL"),
			UserServiceURL:  viper.GetString("USER_SERVICE_URL"),
			AdminServiceURL: viper.GetString("ADMIN_SERVICE_URL"),
		},
		Logging: LoggingConfig{
			Level:       viper.GetString("LOG_LEVEL"),
			Development: viper.GetBool("DEVELOPMENT"),
		},
		JWT: JWTConfig{
			Secret:           viper.GetString("JWT_SECRET"),
			ExpirationHours:  viper.GetInt("JWT_EXPIRATION_HOURS"),
			RefreshSecret:    viper.GetString("JWT_REFRESH_SECRET"),
			RefreshExpHours:  viper.GetInt("JWT_REFRESH_EXPIRATION_HOURS"),
			SigningAlgorithm: viper.GetString("JWT_SIGNING_ALGORITHM"),
			Issuer:           viper.GetString("JWT_ISSUER"),
		},
		CORS: CORSConfig{
			Enabled:          viper.GetBool("CORS_ENABLED"),
			AllowOrigins:     viper.GetStringSlice("CORS_ALLOW_ORIGINS"),
			AllowMethods:     viper.GetStringSlice("CORS_ALLOW_METHODS"),
			AllowHeaders:     viper.GetStringSlice("CORS_ALLOW_HEADERS"),
			ExposeHeaders:    viper.GetStringSlice("CORS_EXPOSE_HEADERS"),
			AllowCredentials: viper.GetBool("CORS_ALLOW_CREDENTIALS"),
			MaxAge:           viper.GetDuration("CORS_MAX_AGE"),
		},
		RateLimiting: RateLimitingConfig{
			Enabled:      viper.GetBool("RATE_LIMIT_ENABLED"),
			Limit:        viper.GetInt("RATE_LIMIT"),
			Burst:        viper.GetInt("RATE_LIMIT_BURST"),
			Window:       viper.GetDuration("RATE_LIMIT_WINDOW"),
			RedisAddress: viper.GetString("REDIS_ADDRESS"),
		},
	}
}

// setDefaults configures all the default values
func setDefaults() {
	// Server defaults
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("READ_TIMEOUT", 5*time.Second)
	viper.SetDefault("WRITE_TIMEOUT", 10*time.Second)
	viper.SetDefault("IDLE_TIMEOUT", 120*time.Second)

	// Services defaults
	viper.SetDefault("AUTH_SERVICE_URL", "http://auth-service:8081")
	viper.SetDefault("USER_SERVICE_URL", "http://user-service:8082")
	viper.SetDefault("ADMIN_SERVICE_URL", "http://admin-service:8083")

	// Logging defaults
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("DEVELOPMENT", true)

	// JWT defaults - using the same secret across all services for consistent authentication
	viper.SetDefault("JWT_SECRET", "your-strong-secret-key-change-in-production")
	viper.SetDefault("JWT_EXPIRATION_HOURS", 24)
	viper.SetDefault("JWT_REFRESH_SECRET", "your-refresh-secret-key-change-in-production")
	viper.SetDefault("JWT_REFRESH_EXPIRATION_HOURS", 168) // 7 days
	viper.SetDefault("JWT_SIGNING_ALGORITHM", "HS256")
	viper.SetDefault("JWT_ISSUER", "qubool-kallyaanam-api")

	// CORS defaults
	viper.SetDefault("CORS_ENABLED", true)
	viper.SetDefault("CORS_ALLOW_ORIGINS", []string{"*"})
	viper.SetDefault("CORS_ALLOW_METHODS", []string{
		"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
	})
	viper.SetDefault("CORS_ALLOW_HEADERS", []string{
		"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID",
		"X-User-ID", "X-User-Role", "X-Username", // Adding our identity propagation headers
	})
	viper.SetDefault("CORS_EXPOSE_HEADERS", []string{
		"Content-Length", "X-Request-ID",
	})
	viper.SetDefault("CORS_ALLOW_CREDENTIALS", true)
	viper.SetDefault("CORS_MAX_AGE", 12*time.Hour)

	// Rate limiting defaults - simplified to only use Redis
	viper.SetDefault("RATE_LIMIT_ENABLED", true)
	viper.SetDefault("RATE_LIMIT", 100)
	viper.SetDefault("RATE_LIMIT_BURST", 150)
	viper.SetDefault("RATE_LIMIT_WINDOW", time.Minute)
	viper.SetDefault("REDIS_ADDRESS", "redis:6379")
}
