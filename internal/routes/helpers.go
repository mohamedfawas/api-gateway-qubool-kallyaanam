// Package routes provides routing definitions for the API Gateway
package routes

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/config"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/middleware"
)

// RouteGroup represents a group of routes with authentication settings
type RouteGroup struct {
	Router      *gin.RouterGroup
	Config      *config.Config
	Logger      *zap.Logger
	RequireAuth bool
	Roles       []string
}

// NewPublicRouteGroup creates a new public route group without authentication
func NewPublicRouteGroup(router *gin.RouterGroup, cfg *config.Config, logger *zap.Logger) *RouteGroup {
	return &RouteGroup{
		Router:      router,
		Config:      cfg,
		Logger:      logger,
		RequireAuth: false,
	}
}

// NewProtectedRouteGroup creates a new protected route group with authentication
func NewProtectedRouteGroup(router *gin.RouterGroup, cfg *config.Config, logger *zap.Logger, roles []string) *RouteGroup {
	// Apply JWT middleware to this group
	router.Use(middleware.JWTAuthMiddleware(cfg, logger))

	// If roles are specified, apply role middleware
	if len(roles) > 0 {
		router.Use(middleware.RoleAuthMiddleware(roles, logger))
	}

	return &RouteGroup{
		Router:      router,
		Config:      cfg,
		Logger:      logger,
		RequireAuth: true,
		Roles:       roles,
	}
}
