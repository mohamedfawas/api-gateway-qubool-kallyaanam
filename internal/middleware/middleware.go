// Package middleware provides HTTP middleware components for the API Gateway
package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RegisterMiddlewares registers all middleware components with the router
func RegisterMiddlewares(router *gin.Engine, logger *zap.Logger) {
	// Add recovery middleware first
	router.Use(gin.Recovery())

	// Add error handler middleware
	router.Use(ErrorHandlerMiddleware(logger))

	// Add logger middleware
	router.Use(LoggerMiddleware(logger))

	// Add standard response formatter middleware
	router.Use(StandardResponseMiddleware())
}
