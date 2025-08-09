package errors

import (
	"fmt"
	"net/http"
)

// ErrorType represents the type of error
type ErrorType string

const (
	// Client errors (4xx)
	ErrorTypeValidation    ErrorType = "VALIDATION_ERROR"
	ErrorTypeNotFound      ErrorType = "NOT_FOUND"
	ErrorTypeUnauthorized  ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden     ErrorType = "FORBIDDEN"
	ErrorTypeBadRequest    ErrorType = "BAD_REQUEST"
	ErrorTypeConflict      ErrorType = "CONFLICT"

	// Server errors (5xx)
	ErrorTypeInternal      ErrorType = "INTERNAL_ERROR"
	ErrorTypeDatabase      ErrorType = "DATABASE_ERROR"
	ErrorTypeExternal      ErrorType = "EXTERNAL_SERVICE_ERROR"
	ErrorTypeTimeout       ErrorType = "TIMEOUT_ERROR"
)

// AppError represents a custom application error
type AppError struct {
	Type       ErrorType   `json:"type"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
	StatusCode int         `json:"-"`
	Err        error       `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new application error
func NewAppError(errorType ErrorType, message string, err error) *AppError {
	appErr := &AppError{
		Type:    errorType,
		Message: message,
		Err:     err,
	}
	
	// Set status code based on error type
	switch errorType {
	case ErrorTypeValidation, ErrorTypeBadRequest:
		appErr.StatusCode = http.StatusBadRequest
	case ErrorTypeNotFound:
		appErr.StatusCode = http.StatusNotFound
	case ErrorTypeUnauthorized:
		appErr.StatusCode = http.StatusUnauthorized
	case ErrorTypeForbidden:
		appErr.StatusCode = http.StatusForbidden
	case ErrorTypeConflict:
		appErr.StatusCode = http.StatusConflict
	case ErrorTypeDatabase, ErrorTypeInternal, ErrorTypeExternal, ErrorTypeTimeout:
		appErr.StatusCode = http.StatusInternalServerError
	default:
		appErr.StatusCode = http.StatusInternalServerError
	}
	
	return appErr
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details interface{}) *AppError {
	e.Details = details
	return e
}

// Convenience functions for common errors

// NewValidationError creates a validation error
func NewValidationError(message string, details interface{}) *AppError {
	return NewAppError(ErrorTypeValidation, message, nil).WithDetails(details)
}

// NewNotFoundError creates a not found error
func NewNotFoundError(resource string) *AppError {
	return NewAppError(ErrorTypeNotFound, fmt.Sprintf("%s not found", resource), nil)
}

// NewDatabaseError creates a database error
func NewDatabaseError(operation string, err error) *AppError {
	return NewAppError(ErrorTypeDatabase, fmt.Sprintf("Database operation failed: %s", operation), err)
}

// NewInternalError creates an internal server error
func NewInternalError(message string, err error) *AppError {
	return NewAppError(ErrorTypeInternal, message, err)
}

// NewBadRequestError creates a bad request error
func NewBadRequestError(message string) *AppError {
	return NewAppError(ErrorTypeBadRequest, message, nil)
}

// NewConflictError creates a conflict error
func NewConflictError(message string) *AppError {
	return NewAppError(ErrorTypeConflict, message, nil)
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(message string) *AppError {
	return NewAppError(ErrorTypeUnauthorized, message, nil)
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) (*AppError, bool) {
	if appErr, ok := err.(*AppError); ok {
		return appErr, true
	}
	return nil, false
}
