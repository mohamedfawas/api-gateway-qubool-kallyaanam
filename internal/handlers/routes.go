package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/middleware"
)

// SetupRoutes configures all routes for the API gateway
func SetupRoutes(router *gin.Engine, proxyHandler *ProxyHandler) {
	// Public routes
	public := router.Group("/")
	{
		// Auth service routes (no auth required)
		authRoutes := public.Group("/auth")
		authRoutes.Any("/*path", proxyHandler.ProxyToService("auth-service"))
	}

	// Protected routes (require authentication)
	protected := router.Group("/")
	protected.Use(middleware.JWTAuth())
	{
		// User service routes
		userRoutes := protected.Group("/user")
		userRoutes.Any("/*path", proxyHandler.ProxyToService("user-service"))
	}

	// Admin routes (require admin role)
	admin := router.Group("/admin")
	admin.Use(middleware.JWTAuth(), middleware.RBACMiddleware([]string{"admin"}))
	admin.Any("/*path", proxyHandler.ProxyToService("admin-service"))
}
