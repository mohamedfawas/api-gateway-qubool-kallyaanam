package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// Initialize sets up the logger
func Initialize(isProduction bool) (*zap.Logger, error) {
	var config zap.Config // holds the configuration settings for the logger.

	if isProduction {
		config = zap.NewProductionConfig() // sets up a structured, high-performance logger with JSON formatting
	} else {
		config = zap.NewDevelopmentConfig() // provides a more human-friendly output, usually formatted for console readability

		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // ensures log levels (INFO, DEBUG, etc.) are printed in capital letters with color coding
	}

	var err error
	log, err = config.Build() // creates the actual logger instance based on the configuration
	if err != nil {
		return nil, err
	}

	return log, nil // If everything is successful, the logger instance (log) is returned
}

// Logger returns the global logger
func Logger() *zap.Logger {
	if log == nil { // checks if log is nil, meaning the logger has not been initialized yet
		log, _ = zap.NewProduction() // If log is nil, it initializes a production logger
	}
	return log
}
