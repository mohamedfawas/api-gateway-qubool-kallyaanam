package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/config"
)

// registerAuthRoutes registers all routes for the auth service
func registerAuthRoutes(rg *gin.RouterGroup, cfg *config.Config, logger *zap.Logger) {
	authPath := rg.Group("/auth")

	// Health check endpoint
	authPath.GET("/health", func(c *gin.Context) {
		serviceURL := cfg.Services.AuthServiceURL + "/health"
		logger.Debug("Forwarding request to auth service health endpoint",
			zap.String("url", serviceURL))
		forwardRequest(c, serviceURL, http.MethodGet, logger)
	})

	// Add more auth routes as needed
	// Examples:
	// auth.POST("/login", handleLogin(cfg, logger))
	// auth.POST("/register", handleRegister(cfg, logger))
	// auth.POST("/refresh-token", handleRefreshToken(cfg, logger))
}
