// Package middleware provides HTTP middleware components for the API Gateway
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/config"
	"go.uber.org/zap"
)

// RegisterMiddlewares registers all middleware components with the router
func RegisterMiddlewares(router *gin.Engine, cfg *config.Config, logger *zap.Logger) {
	// Add recovery middleware first
	router.Use(gin.Recovery())

	// Add logger middleware
	router.Use(LoggerMiddleware(logger))

	// Add CORS middleware early in the chain
	if cfg.CORS.Enabled {
		router.Use(CORSMiddleware(cfg, logger))
	}

	// Add rate limiter middleware
	if cfg.RateLimiting.Enabled {
		router.Use(RateLimiterMiddleware(cfg, logger))
	}

	// Add error handler middleware
	router.Use(ErrorHandlerMiddleware(logger))

	// Add standard response formatter middleware
	router.Use(ResponseWrapper())
}
