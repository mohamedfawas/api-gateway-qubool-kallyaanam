// Package routes defines API routing for the gateway
package routes

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/config"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/constants"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/errors"
)

// Constants for configuration
const (
	defaultTimeout = 10 * time.Second
)

// RegisterRoutes registers all API routes with the router
func RegisterRoutes(router *gin.Engine, cfg *config.Config, logger *zap.Logger) {
	// Register health check routes first
	registerHealthRoutes(router, cfg, logger)

	// API version group
	apiV1 := router.Group("/api/v1")

	// Health check endpoint
	router.GET("/health", healthCheck())

	// Create public groups for each service
	authPublic := NewPublicRouteGroup(apiV1.Group("/auth"), cfg, logger)
	usersPublic := NewPublicRouteGroup(apiV1.Group("/users"), cfg, logger)
	adminPublic := NewPublicRouteGroup(apiV1.Group("/admin"), cfg, logger)

	// Create protected groups for each service
	usersProtected := NewProtectedRouteGroup(apiV1.Group("/users"), cfg, logger, []string{constants.RoleUser})
	adminProtected := NewProtectedRouteGroup(apiV1.Group("/admin"), cfg, logger, []string{constants.RoleAdmin})

	// Register routes
	registerAuthRoutes(authPublic, cfg, logger)
	registerUserPublicRoutes(usersPublic, cfg, logger)
	registerUserProtectedRoutes(usersProtected, cfg, logger)
	registerAdminPublicRoutes(adminPublic, cfg, logger)
	registerAdminProtectedRoutes(adminProtected, cfg, logger)
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

// forwardRequest forwards a request to a service and returns the response
func forwardRequest(c *gin.Context, serviceURL string, method string, logger *zap.Logger) {
	// Create timeout context
	ctx, cancel := context.WithTimeout(c.Request.Context(), defaultTimeout)
	defer cancel()

	// Create the request to the service
	var req *http.Request
	var err error

	// Handle request based on HTTP method
	switch method {
	case http.MethodGet, http.MethodDelete:
		req, err = http.NewRequestWithContext(ctx, method, serviceURL, nil)
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		// Read the request body
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error("Failed to read request body", zap.Error(err))
			c.Error(errors.BadRequestError("Invalid request body", err))
			return
		}
		// Create new request with the body
		req, err = http.NewRequestWithContext(ctx, method, serviceURL, bytes.NewBuffer(bodyBytes))
	default:
		logger.Error("Unsupported HTTP method", zap.String("method", method))
		c.Error(errors.BadRequestError("Unsupported HTTP method", nil))
		return
	}

	if err != nil {
		logger.Error("Failed to create request", zap.Error(err))
		c.Error(errors.InternalError("Failed to create request", err))
		return
	}

	// Copy headers from original request
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Set content type if it's not set
	if req.Header.Get(constants.HeaderContentType) == "" {
		req.Header.Set(constants.HeaderContentType, constants.HeaderApplicationJSON)
	}

	// Forward the request ID if available
	if requestID := c.GetHeader(constants.HeaderRequestID); requestID != "" {
		req.Header.Set(constants.HeaderRequestID, requestID)
	}

	// Send request to service
	client := &http.Client{Timeout: defaultTimeout}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Service request failed",
			zap.String("url", serviceURL),
			zap.Error(err),
		)
		c.Error(errors.ServiceUnavailableError("Service unavailable", err))
		return
	}
	defer resp.Body.Close()

	// Copy the status code
	c.Status(resp.StatusCode)

	// Copy headers from service response
	for key, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}

	// Copy the response body
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		logger.Error("Failed to copy response body", zap.Error(err))
		c.Error(errors.InternalError("Failed to process service response", err))
		return
	}
}
