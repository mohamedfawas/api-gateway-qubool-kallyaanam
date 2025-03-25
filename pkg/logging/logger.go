package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// Initialize sets up the logger
func Initialize(isProduction bool) (*zap.Logger, error) {
	var config zap.Config

	if isProduction {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	var err error
	log, err = config.Build()
	if err != nil {
		return nil, err
	}

	return log, nil
}

// Logger returns the global logger
func Logger() *zap.Logger {
	if log == nil {
		log, _ = zap.NewProduction()
	}
	return log
}
