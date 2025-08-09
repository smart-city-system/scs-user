package errors

import (
	"time"
)

// ErrorResponse represents the structure of error responses sent to clients
type ErrorResponse struct {
	Error     ErrorDetail `json:"error"`
	RequestID string      `json:"request_id,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// ErrorDetail contains the error information
type ErrorDetail struct {
	Type    ErrorType   `json:"type"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// ValidationError represents validation error details
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

// NewErrorResponse creates a new error response
func NewErrorResponse(appErr *AppError, requestID string) *ErrorResponse {
	return &ErrorResponse{
		Error: ErrorDetail{
			Type:    appErr.Type,
			Message: appErr.Message,
			Details: appErr.Details,
		},
		RequestID: requestID,
		Timestamp: time.Now(),
	}
}

// NewValidationErrorResponse creates a validation error response
func NewValidationErrorResponse(errors ValidationErrors, requestID string) *ErrorResponse {
	return &ErrorResponse{
		Error: ErrorDetail{
			Type:    ErrorTypeValidation,
			Message: "Validation failed",
			Details: errors,
		},
		RequestID: requestID,
		Timestamp: time.Now(),
	}
}
