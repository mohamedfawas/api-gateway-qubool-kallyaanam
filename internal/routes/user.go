package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/config"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/constants"
)

// registerUserRoutes registers all routes for the user service
func registerUserPublicRoutes(rg *RouteGroup, cfg *config.Config, logger *zap.Logger) {
	// Health check endpoint
	rg.Router.GET("/health", func(c *gin.Context) {
		serviceURL := cfg.Services.UserServiceURL + "/health"
		logger.Debug("Forwarding request to user service health endpoint",
			zap.String("url", serviceURL))
		forwardRequest(c, serviceURL, http.MethodGet, logger)
	})

	// Add more user routes as needed
	// Examples:
	// user.GET("/profile", handleGetProfile(cfg, logger))
	// user.PUT("/profile", handleUpdateProfile(cfg, logger))
	// user.GET("/matches", handleGetMatches(cfg, logger))
}

func registerUserProtectedRoutes(rg *RouteGroup, cfg *config.Config, logger *zap.Logger) {

	// Forward profile-related requests to the user service
	rg.Router.POST("/profile", func(c *gin.Context) {
		// Pass the request to the user service with the user ID from auth context
		userId, exists := c.Get(constants.ContextKeyUser)
		if !exists {
			logger.Error("User ID not found in context", zap.String("path", c.Request.URL.Path))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Internal server error",
			})
			return
		}

		// Extract user ID and add it to the headers
		userIdStr, ok := userId.(string)
		if !ok {
			logger.Error("User ID has invalid type", zap.String("path", c.Request.URL.Path))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Internal server error",
			})
			return
		}

		// Forward the user ID in the headers
		c.Request.Header.Set("X-User-ID", userIdStr)

		// Forward the request to the user service
		forwardRequest(c, cfg.Services.UserServiceURL+"/api/v1/user/profile", http.MethodPost, logger)
	})

	// Forward GET /profile request to retrieve user profile
	rg.Router.GET("/profile", func(c *gin.Context) {
		userId, exists := c.Get(constants.ContextKeyUser)
		if !exists {
			logger.Error("User ID not found in context", zap.String("path", c.Request.URL.Path))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Internal server error",
			})
			return
		}

		userIdStr, ok := userId.(string)
		if !ok {
			logger.Error("User ID has invalid type", zap.String("path", c.Request.URL.Path))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Internal server error",
			})
			return
		}

		c.Request.Header.Set("X-User-ID", userIdStr)
		forwardRequest(c, cfg.Services.UserServiceURL+"/api/v1/user/profile", http.MethodGet, logger)
	})
	// Protected endpoints...
	// For now, we won't add any since we're focusing on health checks
}
