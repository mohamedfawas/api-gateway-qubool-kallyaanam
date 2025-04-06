// Package utils provides utility functions for the API Gateway
package utils

import (
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates and configures a new zap logger
func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	var loggerConfig zap.Config

	if cfg.Logging.Development {
		loggerConfig = zap.NewDevelopmentConfig()
	} else {
		loggerConfig = zap.NewProductionConfig()
	}

	// Set log level based on configuration
	level := zapcore.InfoLevel
	if err := level.UnmarshalText([]byte(cfg.Logging.Level)); err == nil {
		loggerConfig.Level.SetLevel(level)
	}

	// Build the logger
	return loggerConfig.Build()
}
