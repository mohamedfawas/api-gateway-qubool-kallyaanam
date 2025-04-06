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

// registerAuthRoutes registers all routes for the auth service
func registerAuthRoutes(rg *gin.RouterGroup, cfg *config.Config, logger *zap.Logger) {
	auth := rg.Group("/auth")

	// Health check passthrough
	auth.GET("/health", func(c *gin.Context) {
		forwardRequest(c, cfg.Services.AuthServiceURL+"/health", logger)
	})

	// Add more auth routes here
}

// registerUserRoutes registers all routes for the user service
func registerUserRoutes(rg *gin.RouterGroup, cfg *config.Config, logger *zap.Logger) {
	// Implementation needed
	logger.Info("User routes registered")
}

// registerAdminRoutes registers all routes for the admin service
func registerAdminRoutes(rg *gin.RouterGroup, cfg *config.Config, logger *zap.Logger) {
	// Implementation needed
	logger.Info("Admin routes registered")
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

// forwardRequest forwards requests to the appropriate service with proper error handling
func forwardRequest(c *gin.Context, serviceURL string, logger *zap.Logger) {
	// Create a client with timeout
	client := &http.Client{
		Timeout: defaultTimeout,
	}

	// Create a new request
	req, err := http.NewRequestWithContext(c.Request.Context(), c.Request.Method, serviceURL, c.Request.Body)
	if err != nil {
		logger.Error("Failed to create request",
			zap.String("url", serviceURL),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to create request",
			"error":   err.Error(),
		})
		return
	}

	// Copy headers from original request
	for k, v := range c.Request.Header {
		for _, h := range v {
			req.Header.Add(k, h)
		}
	}

	// Send the request
	resp, err := client.Do(req)
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

	// Forward the status code
	c.Status(resp.StatusCode)

	// Copy headers
	for k, v := range resp.Header {
		for _, h := range v {
			c.Writer.Header().Add(k, h)
		}
	}

	// Copy the body
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		logger.Error("Failed to copy response body",
			zap.String("url", serviceURL),
			zap.Error(err),
		)
	}
}
