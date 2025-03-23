package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/models"
)

// ErrorHandler provides global error handling
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there were any errors during request processing
		if len(c.Errors) > 0 {
			// Take the last error
			err := c.Errors.Last()

			// Return appropriate status code based on error
			statusCode := http.StatusInternalServerError
			if c.Writer.Status() != http.StatusOK {
				statusCode = c.Writer.Status()
			}

			c.JSON(statusCode, models.NewErrorResponse(
				statusCode,
				"Request failed",
				err.Error(),
			))
			return
		}

		// If no response has been sent and we've reached this point
		if c.Writer.Size() == -1 {
			c.JSON(http.StatusNotFound, models.NewErrorResponse(
				http.StatusNotFound,
				"Resource not found",
				nil,
			))
		}
	}
}
