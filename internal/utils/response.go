// Package utils provides utility functions and helpers
package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// StandardResponse is the common response structure for all API endpoints
type StandardResponse struct {
	Status  bool        `json:"status"`  // Success status (true/false)
	Message string      `json:"message"` // Human-readable message
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// RespondWithSuccess sends a standardized success response
func RespondWithSuccess(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, StandardResponse{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

// RespondWithCreated sends a standardized created response
func RespondWithCreated(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, StandardResponse{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

// RespondWithError sends a standardized error response
func RespondWithError(c *gin.Context, statusCode int, message string, err interface{}) {
	c.JSON(statusCode, StandardResponse{
		Status:  false,
		Message: message,
		Error:   err,
	})
}

// RespondWithBadRequest sends a standardized bad request response
func RespondWithBadRequest(c *gin.Context, message string, err interface{}) {
	RespondWithError(c, http.StatusBadRequest, message, err)
}

// RespondWithUnauthorized sends a standardized unauthorized response
func RespondWithUnauthorized(c *gin.Context, message string, err interface{}) {
	RespondWithError(c, http.StatusUnauthorized, message, err)
}

// RespondWithForbidden sends a standardized forbidden response
func RespondWithForbidden(c *gin.Context, message string, err interface{}) {
	RespondWithError(c, http.StatusForbidden, message, err)
}

// RespondWithNotFound sends a standardized not found response
func RespondWithNotFound(c *gin.Context, message string, err interface{}) {
	RespondWithError(c, http.StatusNotFound, message, err)
}

// RespondWithInternalError sends a standardized internal server error response
func RespondWithInternalError(c *gin.Context, message string, err interface{}) {
	RespondWithError(c, http.StatusInternalServerError, message, err)
}

// RespondWithServiceUnavailable sends a standardized service unavailable response
func RespondWithServiceUnavailable(c *gin.Context, message string, err interface{}) {
	RespondWithError(c, http.StatusServiceUnavailable, message, err)
}

// RespondWithTooManyRequests sends a standardized too many requests response
func RespondWithTooManyRequests(c *gin.Context, message string, err interface{}) {
	RespondWithError(c, http.StatusTooManyRequests, message, err)
}
