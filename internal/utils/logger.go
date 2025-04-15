// Package utils provides utility functions for the API Gateway
package utils

import (
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates and configures a new zap logger
func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	var loggerConfig zap.Config //  This will hold the configuration for the new logger

	if cfg.Logging.Development {
		// set up a logger configuration designed for development.
		// means more human-readable output and additional debugging features.
		loggerConfig = zap.NewDevelopmentConfig()
	} else {
		// set up a logger configuration designed for production.
		// means more concise output and better performance.
		loggerConfig = zap.NewProductionConfig()
	}

	// Set the default log level to "info"
	// Log levels determine what type of messages are shown: debug < info < warn < error
	level := zapcore.InfoLevel

	// Try to convert the configured log level (from string) into a zap log level
	// Example: if cfg.Logging.Level = "debug", it'll become zapcore.DebugLevel
	if err := level.UnmarshalText([]byte(cfg.Logging.Level)); err == nil {
		// If conversion is successful, update the logger config with this level
		loggerConfig.Level.SetLevel(level)
	}

	// Build the logger using the configured settings and return it
	return loggerConfig.Build()
}
