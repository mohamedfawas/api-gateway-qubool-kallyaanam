package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/constants"
)

// StandardResponse represents the standard API response format
type StandardResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// ResponseBodyWriter wraps gin.ResponseWriter to capture response body
type ResponseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write captures response body while writing it to the client
func (w *ResponseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// StandardResponseMiddleware formats all successful responses (2xx) to match the standard format
func StandardResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create response body capturer
		bodyBuffer := new(bytes.Buffer)
		bodyWriter := &ResponseBodyWriter{
			ResponseWriter: c.Writer,
			body:           bodyBuffer,
		}
		c.Writer = bodyWriter

		// Process the request
		c.Next()

		// Only format successful responses
		if !isSuccessStatus(c.Writer.Status()) {
			return
		}

		// Skip if response body is empty
		responseBody := bodyBuffer.Bytes()
		if len(responseBody) == 0 {
			return
		}

		// Skip if already in standard format
		if isAlreadyStandardFormat(responseBody) {
			return
		}

		// Create standardized response
		data := parseResponseData(responseBody)
		standardResponse := StandardResponse{
			Status:  true,
			Message: constants.MessageSuccess,
			Data:    data,
		}

		// Send standardized response
		c.Writer.Header().Set(constants.HeaderContentType, constants.HeaderApplicationJSON)
		c.JSON(c.Writer.Status(), standardResponse)
	}
}

// isSuccessStatus checks if HTTP status code is in the 2xx range
func isSuccessStatus(statusCode int) bool {
	return statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices
}

// isAlreadyStandardFormat checks if response is already in standard format
func isAlreadyStandardFormat(body []byte) bool {
	if len(body) == 0 || body[0] != '{' {
		return false
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return false
	}

	_, hasStatus := response["status"]
	return hasStatus
}

// parseResponseData converts response body to appropriate data format
func parseResponseData(body []byte) interface{} {
	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		// If not valid JSON, use as string
		return string(body)
	}
	return data
}
