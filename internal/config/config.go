package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Server struct {
		Port        string
		Environment string
	}
	Services map[string]string
	JWT      struct {
		Secret     string
		Expiration int // in minutes
	}
	CORS struct {
		AllowedOrigins []string
		AllowedMethods []string
		AllowedHeaders []string
	}
}

// Load loads configuration from environment variables and config files
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Set default values
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.environment", "development")

	// Service defaults (should be overridden in production)
	viper.SetDefault("services.auth-service", "http://auth-service:8080")
	viper.SetDefault("services.user-service", "http://user-service:8080")
	viper.SetDefault("services.admin-service", "http://admin-service:8080")

	// JWT defaults
	viper.SetDefault("jwt.expiration", 60) // 60 minutes

	// CORS defaults
	viper.SetDefault("cors.allowed_origins", []string{"*"})
	viper.SetDefault("cors.allowed_methods", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"})
	viper.SetDefault("cors.allowed_headers", []string{"Origin", "Content-Type", "Accept", "Authorization"})

	// Environment variables take precedence over config files
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Try to read config file, but continue if not found
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	// For security, JWT secret must be set via environment or config
	if viper.GetString("jwt.secret") == "" {
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			// In development, we can use a default secret, but warn about it
			if viper.GetString("server.environment") != "production" {
				secret = "dev-secret-do-not-use-in-production"
				println("WARNING: Using default JWT secret. This is insecure for production.")
			} else {
				return nil, ErrMissingJWTSecret
			}
		}
		viper.Set("jwt.secret", secret)
	}

	// Create config struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Custom errors
var (
	ErrMissingJWTSecret = &configError{"JWT secret is required"}
)

type configError struct {
	message string
}

func (e *configError) Error() string {
	return e.message
}
