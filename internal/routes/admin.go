package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/config"
)

// registerAdminPublicRoutes registers public routes for the admin service
func registerAdminPublicRoutes(rg *RouteGroup, cfg *config.Config, logger *zap.Logger) {
	// Health check endpoint
	rg.Router.GET("/health", func(c *gin.Context) {
		serviceURL := cfg.Services.AdminServiceURL + "/health"
		logger.Debug("Forwarding request to admin service health endpoint",
			zap.String("url", serviceURL))
		forwardRequest(c, serviceURL, http.MethodGet, logger)
	})

	// Other public endpoints...
}

// registerAdminProtectedRoutes registers protected routes for the admin service
func registerAdminProtectedRoutes(rg *RouteGroup, cfg *config.Config, logger *zap.Logger) {
	// Protected endpoints...
	// For now, we won't add any since we're focusing on health checks
}
