package constants

// HTTP Headers
const (
	HeaderContentType     = "Content-Type"
	HeaderApplicationJSON = "application/json"
	HeaderRequestID       = "X-Request-ID"
	HeaderAuthorization   = "Authorization"
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

// Error messages
const (
	ErrInvalidRequest     = "Invalid request"
	ErrServiceUnavailable = "Service unavailable"
	ErrInternalServer     = "Internal server error"
	ErrResourceNotFound   = "Resource not found"
)

// Authentication constants
const (
	// Roles
	RoleAdmin = "admin"
	RoleUser  = "user"

	// Context keys
	ContextKeyUser = "user"
)
