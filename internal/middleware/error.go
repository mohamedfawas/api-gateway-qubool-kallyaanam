package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/constants"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/errors"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/utils"
	"go.uber.org/zap"
)

// ErrorHandlerMiddleware is a function that returns a Gin middleware.
// This middleware catches and handles all errors in one place.
func ErrorHandlerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Call the next handler in the chain. This could be another middleware or the actual endpoint handler.
		c.Next()

		// If there are no errors in the context, just exit and do nothing.
		if len(c.Errors) == 0 {
			return
		}

		// If there were errors, get the last error that happened
		err := c.Errors.Last()

		// Log the error using zap logger
		logger.Error("Request error", zap.Error(err.Err))

		// Check if the error is a custom APIError type
		if apiErr, ok := err.Err.(*errors.APIError); ok {
			// Get the HTTP status code from the custom error
			statusCode := apiErr.StatusCode()

			// Respond with a proper error message, using a helper function
			utils.RespondWithError(c,
				statusCode,                          // HTTP status code like 400, 401, etc.
				getMessageForStatusCode(statusCode), // A readable error message string
				apiErr.ToResponse())                 // The full error response
			return
		}

		// If it's not a custom API error, treat it as an internal server error (generic case)
		utils.RespondWithInternalError(c,
			constants.StatusInternalServerError,
			map[string]string{ // Build a basic error response
				"type":    string(errors.ErrorTypeInternal), // error type like "internal"
				"message": err.Error(),                      // the error message
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
