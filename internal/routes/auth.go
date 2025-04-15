package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/config"
)

// registerAuthRoutes registers all routes for the auth service
func registerAuthRoutes(rg *RouteGroup, cfg *config.Config, logger *zap.Logger) {
	// Health check endpoint
	rg.Router.GET("/health", func(c *gin.Context) {
		serviceURL := cfg.Services.AuthServiceURL + "/health"
		logger.Debug("Forwarding request to auth service health endpoint",
			zap.String("url", serviceURL))
		forwardRequest(c, serviceURL, http.MethodGet, logger)
	})

	// Add register endpoint
	rg.Router.POST("/register", func(c *gin.Context) {
		serviceURL := cfg.Services.AuthServiceURL + "/auth/register"
		logger.Debug("Forwarding request to auth service register endpoint",
			zap.String("url", serviceURL))
		forwardRequest(c, serviceURL, http.MethodPost, logger)
	})

	// Add verify email endpoint
	rg.Router.POST("/verify-email", func(c *gin.Context) {
		serviceURL := cfg.Services.AuthServiceURL + "/auth/verify-email"
		logger.Debug("Forwarding request to auth service verify email endpoint",
			zap.String("url", serviceURL))
		forwardRequest(c, serviceURL, http.MethodPost, logger)
	})

	// Add login endpoint
	rg.Router.POST("/login", func(c *gin.Context) {
		serviceURL := cfg.Services.AuthServiceURL + "/auth/login"
		logger.Debug("Forwarding request to auth service login endpoint",
			zap.String("url", serviceURL))
		forwardRequest(c, serviceURL, http.MethodPost, logger)
	})

	// Add more auth routes as needed
	// Examples:
	// auth.POST("/login", handleLogin(cfg, logger))
	// auth.POST("/register", handleRegister(cfg, logger))
	// auth.POST("/refresh-token", handleRefreshToken(cfg, logger))
}
