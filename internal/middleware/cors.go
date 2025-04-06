// Package middleware provides HTTP middleware components for the API Gateway
package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/config"
)

// CORSMiddleware creates a middleware for handling CORS
func CORSMiddleware(cfg *config.Config, logger *zap.Logger) gin.HandlerFunc {
	corsConfig := cors.Config{
		AllowOrigins:     cfg.CORS.AllowOrigins,
		AllowMethods:     cfg.CORS.AllowMethods,
		AllowHeaders:     cfg.CORS.AllowHeaders,
		ExposeHeaders:    cfg.CORS.ExposeHeaders,
		AllowCredentials: cfg.CORS.AllowCredentials,
		MaxAge:           cfg.CORS.MaxAge,
	}

	// Log CORS configuration
	logger.Info("Configured CORS middleware",
		zap.Strings("allowOrigins", cfg.CORS.AllowOrigins),
		zap.Strings("allowMethods", cfg.CORS.AllowMethods),
		zap.Strings("allowHeaders", cfg.CORS.AllowHeaders),
		zap.Strings("exposeHeaders", cfg.CORS.ExposeHeaders),
		zap.Bool("allowCredentials", cfg.CORS.AllowCredentials),
		zap.Duration("maxAge", cfg.CORS.MaxAge),
	)

	return cors.New(corsConfig)
}
