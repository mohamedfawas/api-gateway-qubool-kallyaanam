package models

// Response represents the standard API response format
type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// NewSuccessResponse creates a standard success response
func NewSuccessResponse(status int, message string, data interface{}) Response {
	return Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

// NewErrorResponse creates a standard error response
func NewErrorResponse(status int, message string, errorDetails interface{}) Response {
	return Response{
		Status:  status,
		Message: message,
		Error:   errorDetails,
	}
}
