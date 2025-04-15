// Package middleware provides HTTP middleware components for the API Gateway
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/constants"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/errors"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/utils"
	"go.uber.org/zap"
)

// ErrorHandlerMiddleware handles all errors in a standardized way
func ErrorHandlerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last()
		logger.Error("Request error", zap.Error(err.Err))

		// Handle custom API errors
		if apiErr, ok := err.Err.(*errors.APIError); ok {
			statusCode := apiErr.StatusCode()
			utils.RespondWithError(c, statusCode, getMessageForStatusCode(statusCode), apiErr.ToResponse())
			return
		}

		// Handle generic errors
		utils.RespondWithInternalError(c, constants.StatusInternalServerError, map[string]string{
			"type":    string(errors.ErrorTypeInternal),
			"message": err.Error(),
		})
	}
}

// getMessageForStatusCode returns an appropriate error message for the status code
func getMessageForStatusCode(statusCode int) string {
	switch statusCode {
	case http.StatusBadRequest:
		return constants.StatusBadRequest
	case http.StatusUnauthorized:
		return constants.StatusUnauthorized
	case http.StatusForbidden:
		return constants.StatusForbidden
	case http.StatusNotFound:
		return constants.StatusNotFound
	case http.StatusServiceUnavailable:
		return constants.StatusServiceUnavailable
	case http.StatusTooManyRequests:
		return constants.StatusTooManyRequests
	default:
		return constants.StatusInternalServerError
	}
}
