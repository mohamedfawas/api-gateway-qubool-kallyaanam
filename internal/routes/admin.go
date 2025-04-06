package routes

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/config"
)

// registerAdminRoutes registers all routes for the admin service
func registerAdminRoutes(rg *gin.RouterGroup, cfg *config.Config, logger *zap.Logger) {
	admin := rg.Group("/admin")

	// Health check endpoint
	admin.GET("/health", func(c *gin.Context) {
		logger.Debug("Forwarding request to admin service health endpoint")
		forwardRequest(c, cfg.Services.AdminServiceURL+"/health", logger)
	})

	// Add more admin routes as needed
	// Examples:
	// admin.GET("/users", handleListUsers(cfg, logger))
	// admin.POST("/users/block", handleBlockUser(cfg, logger))
	// admin.GET("/dashboard", handleDashboard(cfg, logger))
}
