package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/config"
)

// registerAdminRoutes registers all routes for the admin service
func registerAdminRoutes(rg *gin.RouterGroup, cfg *config.Config, logger *zap.Logger) {
	adminPath := rg.Group("/admin")

	// Health check endpoint
	adminPath.GET("/health", func(c *gin.Context) {
		serviceURL := cfg.Services.AdminServiceURL + "/health"
		logger.Debug("Forwarding request to admin service health endpoint",
			zap.String("url", serviceURL))
		forwardRequest(c, serviceURL, http.MethodGet, logger)
	})

	// Add more admin routes as needed
	// Examples:
	// admin.GET("/users", handleListUsers(cfg, logger))
	// admin.POST("/users/block", handleBlockUser(cfg, logger))
	// admin.GET("/dashboard", handleDashboard(cfg, logger))
}
