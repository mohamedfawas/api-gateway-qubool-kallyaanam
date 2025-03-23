package models

// Response provides a standardized structure for all API responses
type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// NewSuccessResponse creates a success response
func NewSuccessResponse(status int, message string, data interface{}) Response {
	return Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

// NewErrorResponse creates an error response
func NewErrorResponse(status int, message string, err interface{}) Response {
	return Response{
		Status:  status,
		Message: message,
		Error:   err,
	}
}
