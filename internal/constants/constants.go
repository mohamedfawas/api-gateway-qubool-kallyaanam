// Package constants provides constant values used throughout the application
package constants

// HTTP Headers
const (
	HeaderContentType     = "Content-Type"
	HeaderApplicationJSON = "application/json"
	HeaderRequestID       = "X-Request-ID"
)

// Response messages
const (
	MessageSuccess = "Success"
	MessageError   = "Error"
)

// Status codes descriptions
const (
	StatusOK                  = "OK"
	StatusBadRequest          = "Bad Request"
	StatusUnauthorized        = "Unauthorized"
	StatusForbidden           = "Forbidden"
	StatusNotFound            = "Not Found"
	StatusInternalServerError = "Internal Server Error"
	StatusServiceUnavailable  = "Service Unavailable"
)
