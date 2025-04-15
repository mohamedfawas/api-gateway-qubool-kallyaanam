// api-gateway-qubool-kallyaanam/internal/middleware/response.go
package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// StandardResponse is the common response structure for all API endpoints
type StandardResponse struct {
	Status  int         `json:"status"` // HTTP status code
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ResponseWriter extends gin.ResponseWriter to capture the response body
type ResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write captures the response and writes to the underlying ResponseWriter
func (w *ResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// WriteString captures the response string and writes to the underlying ResponseWriter
func (w *ResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// ResponseWrapper standardizes all API responses
func ResponseWrapper() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a buffer to capture the response
		responseBuffer := &bytes.Buffer{}

		// Create a custom ResponseWriter that captures the response
		writer := &ResponseWriter{
			ResponseWriter: c.Writer,
			body:           responseBuffer,
		}

		// Replace the ResponseWriter with our custom one
		c.Writer = writer

		// Process the request
		c.Next()

		// If we have errors, don't modify the response (let error handler do its job)
		if len(c.Errors) > 0 {
			return
		}

		// Only process JSON responses with success status codes
		contentType := writer.Header().Get("Content-Type")
		statusCode := writer.Status()

		if statusCode >= 200 && statusCode < 400 && contentType == "application/json" {
			// Try to parse the response
			var originalResponse map[string]interface{}

			// Check if we can unmarshal the response body
			if err := json.Unmarshal(responseBuffer.Bytes(), &originalResponse); err == nil {
				// Skip if already standardized
				if _, hasStatus := originalResponse["status"]; hasStatus {
					return
				}

				// Create standardized response
				standardResp := StandardResponse{
					Status:  statusCode,
					Message: http.StatusText(statusCode),
					Data:    originalResponse,
				}

				// Convert to JSON
				newBody, err := json.Marshal(standardResp)
				if err != nil {
					return
				}

				// Reset buffer and headers
				c.Writer.Header().Set("Content-Length", string(len(newBody)))
				c.Writer.Header().Set("Content-Type", "application/json")

				// Write the standardized response
				c.Writer.WriteHeader(statusCode)
				c.Writer.Write(newBody)
			}
		}
	}
}
