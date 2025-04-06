package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RegisterMiddlewares registers all middleware components with the router
func RegisterMiddlewares(router *gin.Engine, logger *zap.Logger) {
	// Add logger middleware
	router.Use(LoggerMiddleware(logger))

	// Add standard response formatter
	router.Use(StandardResponseMiddleware())
}
