package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/models"
)

// SimpleErrorHandler provides global error handling
func SimpleErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there were any errors during request processing
		if len(c.Errors) > 0 {
			// Use the last error as the response error
			lastErr := c.Errors.Last()

			// Determine HTTP status code
			statusCode := http.StatusInternalServerError
			if c.Writer.Status() != http.StatusOK {
				statusCode = c.Writer.Status()
			}

			c.JSON(statusCode, models.NewErrorResponse(
				statusCode,
				"Request failed",
				lastErr.Error(),
			))
			return
		}

		// If no handler wrote a response, return 404
		if c.Writer.Size() == -1 {
			c.JSON(http.StatusNotFound, models.NewErrorResponse(
				http.StatusNotFound,
				"Resource not found",
				"",
			))
		}
	}
}
