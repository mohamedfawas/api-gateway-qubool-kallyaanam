package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/config"
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
	// Protected endpoints...
	// For now, we won't add any since we're focusing on health checks
}
