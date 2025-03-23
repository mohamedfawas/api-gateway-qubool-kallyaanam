package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/models"
	"go.uber.org/zap"
)

// ErrorHandler middleware catches panics and returns a standardized error response
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error with stack trace
				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("request-id", c.GetString("requestID")),
				)

				// Return error response
				c.AbortWithStatusJSON(http.StatusInternalServerError, models.NewErrorResponse(
					http.StatusInternalServerError,
					"Internal server error",
					nil, // Don't expose the actual error details to clients
				))
			}
		}()

		c.Next()
	}
}
