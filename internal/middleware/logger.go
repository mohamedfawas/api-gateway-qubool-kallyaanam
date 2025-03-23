package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func init() {
	// Initialize logger
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	logger, err = config.Build()
	if err != nil {
		panic(err)
	}
}

// Logger middleware logs each request
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Add trace ID to context if needed
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		c.Set("requestID", requestID)
		c.Header("X-Request-ID", requestID)

		// Process request
		c.Next()

		// Calculate request duration
		duration := time.Since(start)

		// Log request details
		logger.Info("Request",
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Int("status", c.Writer.Status()),
			zap.Int("size", c.Writer.Size()),
			zap.Duration("duration", duration),
			zap.String("request-id", requestID),
		)
	}
}

// generateRequestID creates a unique request ID
func generateRequestID() string {
	// Simple implementation - in production, use a more robust method
	return time.Now().Format("20060102150405") + "-" +
		RandStringBytes(8)
}

// RandStringBytes generates a random string of n bytes
func RandStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[time.Now().UnixNano()%int64(len(letterBytes))]
		time.Sleep(1 * time.Nanosecond) // Ensure unique values even in fast loops
	}
	return string(b)
}
