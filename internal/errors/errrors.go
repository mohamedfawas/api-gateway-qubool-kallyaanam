// Package errors provides custom error types and handling for the API Gateway
package errors

import (
	"fmt"
	"net/http"
)

// ErrorType represents the category of an error
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

// APIError represents a standard API error
type APIError struct {
	Type    ErrorType   `json:"type"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
	Err     error       `json:"-"` // Internal error, not exposed
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s - %v", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// StatusCode returns the HTTP status code for this error
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

// New creates a new APIError
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

// Add a convenience function for rate limited errors
func RateLimitedError(message string) *APIError {
	return New(ErrorTypeRateLimited, message, nil)
}
