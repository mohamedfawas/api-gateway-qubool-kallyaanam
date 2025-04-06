// Package routes defines API routing for the gateway
package routes

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/config"
)

// Constants for configuration
const (
	defaultTimeout = 10 * time.Second
)

// RegisterRoutes registers all API routes with the router
func RegisterRoutes(router *gin.Engine, cfg *config.Config, logger *zap.Logger) {
	// API version group
	apiV1 := router.Group("/api/v1")

	// Health check endpoint
	router.GET("/health", healthCheck())

	// Register service routes
	registerAuthRoutes(apiV1, cfg, logger)
	registerUserRoutes(apiV1, cfg, logger)
	registerAdminRoutes(apiV1, cfg, logger)
}

// healthCheck provides a health check endpoint for the API Gateway
func healthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "api-gateway",
			"status":  "UP",
			"version": "0.1.0",
		})
	}
}

// Example of what the forwardRequest function might look like in your routes.go
func forwardRequest(c *gin.Context, serviceURL string, logger *zap.Logger) {
	client := &http.Client{
		Timeout: defaultTimeout,
	}

	resp, err := client.Get(serviceURL)
	if err != nil {
		logger.Error("Service request failed",
			zap.String("url", serviceURL),
			zap.Error(err),
		)

		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  false,
			"message": "Service unavailable",
			"error":   err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	// Copy the status code
	c.Status(resp.StatusCode)

	// Copy headers
	for k, v := range resp.Header {
		for _, h := range v {
			c.Writer.Header().Add(k, h)
		}
	}

	// Copy the body
	io.Copy(c.Writer, resp.Body)
}
