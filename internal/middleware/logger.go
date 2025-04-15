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
		// Record the time when the request starts
		start := time.Now()

		// Try to get the request ID from the incoming request headers
		requestID := c.GetHeader(constants.HeaderRequestID)

		// If no request ID is provided, generate a new one
		if requestID == "" {
			requestID = uuid.New().String()                            // Create a new unique request ID
			c.Request.Header.Set(constants.HeaderRequestID, requestID) // Set it in the request header
		}

		// Add the request ID to the response headers so clients can see it
		c.Writer.Header().Set(constants.HeaderRequestID, requestID)

		// Call the next middleware/handler in the chain to continue request processing
		c.Next()

		// After the request is finished, record the end time
		end := time.Now()

		// Calculate how long the request took
		latency := end.Sub(start)

		// Get the requested URL path (e.g., /api/user)
		path := c.Request.URL.Path

		// If there are any query parameters (e.g., ?name=joe), append them to the path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		// Log the details of the request using zap logger
		logger.Info("Request processed",
			zap.String("requestID", requestID),             // The unique ID for this request
			zap.String("method", c.Request.Method),         // HTTP method used (GET, POST, etc.)
			zap.String("path", path),                       // Full request path
			zap.Int("status", c.Writer.Status()),           // HTTP status code (200, 404, etc.)
			zap.Int("size", c.Writer.Size()),               // Size of the response in bytes
			zap.Duration("latency", latency),               // Time taken to process the request
			zap.String("clientIP", c.ClientIP()),           // IP address of the client
			zap.String("userAgent", c.Request.UserAgent()), // Browser or client making the request
		)
	}
}
