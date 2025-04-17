package errors

import (
	"fmt"
	"net/http"
)

// ErrorType is a custom type to categorize different kinds of errors
type ErrorType string

// Define error types
const (
	ErrorTypeValidation         ErrorType = "VALIDATION_ERROR"
	ErrorTypeBadRequest         ErrorType = "BAD_REQUEST"
	ErrorTypeUnauthorized       ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden          ErrorType = "FORBIDDEN"
	ErrorTypeNotFound           ErrorType = "NOT_FOUND"
	ErrorTypeInternal           ErrorType = "INTERNAL_ERROR"
	ErrorTypeServiceUnavailable ErrorType = "SERVICE_UNAVAILABLE"
	ErrorTypeRateLimited        ErrorType = "RATE_LIMITED"
)

// APIError is a struct that represents an error in a standard format for APIs
type APIError struct {
	Type    ErrorType   `json:"type"`              // The type of error (e.g., BAD_REQUEST)
	Message string      `json:"message"`           // A human-readable message describing the error
	Details interface{} `json:"details,omitempty"` // Optional extra information about the error
	Err     error       `json:"-"`                 // The original error (hidden in JSON responses), not exposed
}

// This method allows APIError to satisfy Go's built-in `error` interface
func (e *APIError) Error() string {
	// If an internal error exists, include it in the error string
	if e.Err != nil {
		return fmt.Sprintf("%s: %s - %v", e.Type, e.Message, e.Err)
	}

	// If there's no internal error, return only the type and message
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// StatusCode returns the correct HTTP status code based on the error type
func (e *APIError) StatusCode() int {
	switch e.Type {
	case ErrorTypeValidation, ErrorTypeBadRequest:
		return http.StatusBadRequest
	case ErrorTypeUnauthorized:
		return http.StatusUnauthorized
	case ErrorTypeForbidden:
		return http.StatusForbidden
	case ErrorTypeNotFound:
		return http.StatusNotFound
	case ErrorTypeServiceUnavailable:
		return http.StatusServiceUnavailable
	case ErrorTypeRateLimited:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}

// ToResponse converts the APIError into a map that can be easily converted to JSON for API responses
func (e *APIError) ToResponse() map[string]interface{} {
	// Basic structure of the response
	response := map[string]interface{}{
		"type":    string(e.Type),
		"message": e.Message,
	}
	// Include details if they exist
	if e.Details != nil {
		response["details"] = e.Details
	}
	return response
}

// New creates a new APIError with the given type, message, and internal error
func New(errorType ErrorType, message string, err error) *APIError {
	return &APIError{
		Type:    errorType,
		Message: message,
		Err:     err,
	}
}

// NewWithDetails creates a new APIError with details
func NewWithDetails(errorType ErrorType, message string, details interface{}, err error) *APIError {
	return &APIError{
		Type:    errorType,
		Message: message,
		Details: details,
		Err:     err,
	}
}

// ValidationError creates a new validation error
func ValidationError(message string, details interface{}) *APIError {
	return NewWithDetails(ErrorTypeValidation, message, details, nil)
}

// BadRequestError creates a new bad request error
func BadRequestError(message string, err error) *APIError {
	return New(ErrorTypeBadRequest, message, err)
}

// NotFoundError creates a new not found error
func NotFoundError(message string) *APIError {
	return New(ErrorTypeNotFound, message, nil)
}

// ServiceUnavailableError creates a new service unavailable error
func ServiceUnavailableError(message string, err error) *APIError {
	return New(ErrorTypeServiceUnavailable, message, err)
}

// InternalError creates a new internal server error
func InternalError(message string, err error) *APIError {
	return New(ErrorTypeInternal, message, err)
}

// RateLimitedError creates a new rate limited error
func RateLimitedError(message string) *APIError {
	return New(ErrorTypeRateLimited, message, nil)
}
