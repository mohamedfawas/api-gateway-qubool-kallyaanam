package routes

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/config"
)

// registerUserRoutes registers all routes for the user service
func registerUserRoutes(rg *gin.RouterGroup, cfg *config.Config, logger *zap.Logger) {
	user := rg.Group("/user")

	// Health check endpoint
	user.GET("/health", func(c *gin.Context) {
		logger.Debug("Forwarding request to user service health endpoint")
		forwardRequest(c, cfg.Services.UserServiceURL+"/health", logger)
	})

	// Add more user routes as needed
	// Examples:
	// user.GET("/profile", handleGetProfile(cfg, logger))
	// user.PUT("/profile", handleUpdateProfile(cfg, logger))
	// user.GET("/matches", handleGetMatches(cfg, logger))
}
