// Package middleware provides HTTP middleware components for the API Gateway
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/constants"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/errors"
	"go.uber.org/zap"
)

// ErrorResponse is the standard error response format
type ErrorResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Error   interface{} `json:"error,omitempty"`
}

func ErrorHandlerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last()
		logger.Error("Request error", zap.Error(err.Err))

		var statusCode int
		var errDetails interface{}

		if apiErr, ok := err.Err.(*errors.APIError); ok {
			statusCode = apiErr.StatusCode()
			errDetails = apiErr
		} else {
			statusCode = http.StatusInternalServerError
			errDetails = map[string]string{
				"type":    string(errors.ErrorTypeInternal),
				"message": err.Error(),
			}
		}

		c.JSON(statusCode, ErrorResponse{
			Status:  false,
			Message: getMessageForStatusCode(statusCode),
			Error:   errDetails,
		})

		c.Abort()
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
	default:
		return constants.StatusInternalServerError
	}
}
