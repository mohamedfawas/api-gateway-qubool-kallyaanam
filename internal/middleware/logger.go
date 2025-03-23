package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger middleware logs request details
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(startTime)
		statusCode := c.Writer.Status()

		// Get user ID from context if set by auth middleware
		userID, _ := c.Get("userID")

		// Log request details
		if userID != nil {
			log.Printf("[%d] %s %s - %v | IP: %s | User: %v",
				statusCode, method, path, latency, clientIP, userID)
		} else {
			log.Printf("[%d] %s %s - %v | IP: %s",
				statusCode, method, path, latency, clientIP)
		}
	}
}
