package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/constants"
	"go.uber.org/zap"
)

// LoggerMiddleware logs request details using zap
func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Generate request ID if it doesn't exist
		requestID := c.GetHeader(constants.HeaderRequestID)
		if requestID == "" {
			requestID = uuid.New().String()
			c.Request.Header.Set(constants.HeaderRequestID, requestID)
		}

		// Set the request ID in response headers
		c.Writer.Header().Set(constants.HeaderRequestID, requestID)

		// Process request
		c.Next()

		// Log request details
		end := time.Now()
		latency := end.Sub(start)

		// Get path, status code, and method
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		logger.Info("Request processed",
			zap.String("requestID", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Int("status", c.Writer.Status()),
			zap.Int("size", c.Writer.Size()),
			zap.Duration("latency", latency),
			zap.String("clientIP", c.ClientIP()),
			zap.String("userAgent", c.Request.UserAgent()),
		)
	}
}
