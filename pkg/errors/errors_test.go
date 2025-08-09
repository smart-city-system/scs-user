package errors

import (
	"errors"
	"testing"
)

func TestNewAppError(t *testing.T) {
	tests := []struct {
		name           string
		errorType      ErrorType
		message        string
		err            error
		expectedStatus int
	}{
		{
			name:           "Validation Error",
			errorType:      ErrorTypeValidation,
			message:        "Invalid input",
			err:            nil,
			expectedStatus: 400,
		},
		{
			name:           "Not Found Error",
			errorType:      ErrorTypeNotFound,
			message:        "Resource not found",
			err:            nil,
			expectedStatus: 404,
		},
		{
			name:           "Database Error",
			errorType:      ErrorTypeDatabase,
			message:        "Database connection failed",
			err:            errors.New("connection timeout"),
			expectedStatus: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appErr := NewAppError(tt.errorType, tt.message, tt.err)
			
			if appErr.Type != tt.errorType {
				t.Errorf("Expected error type %s, got %s", tt.errorType, appErr.Type)
			}
			
			if appErr.Message != tt.message {
				t.Errorf("Expected message %s, got %s", tt.message, appErr.Message)
			}
			
			if appErr.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, appErr.StatusCode)
			}
			
			if appErr.Err != tt.err {
				t.Errorf("Expected underlying error %v, got %v", tt.err, appErr.Err)
			}
		})
	}
}

func TestConvenienceFunctions(t *testing.T) {
	// Test NewValidationError
	validationErr := NewValidationError("Invalid field", map[string]string{"field": "name"})
	if validationErr.Type != ErrorTypeValidation {
		t.Errorf("Expected validation error type")
	}

	// Test NewNotFoundError
	notFoundErr := NewNotFoundError("user")
	if notFoundErr.Type != ErrorTypeNotFound {
		t.Errorf("Expected not found error type")
	}
	if notFoundErr.Message != "user not found" {
		t.Errorf("Expected 'user not found', got %s", notFoundErr.Message)
	}

	// Test NewDatabaseError
	dbErr := NewDatabaseError("insert", errors.New("constraint violation"))
	if dbErr.Type != ErrorTypeDatabase {
		t.Errorf("Expected database error type")
	}
}

func TestIsAppError(t *testing.T) {
	// Test with AppError
	appErr := NewValidationError("test", nil)
	if result, ok := IsAppError(appErr); !ok || result != appErr {
		t.Errorf("Expected IsAppError to return true for AppError")
	}

	// Test with regular error
	regularErr := errors.New("regular error")
	if _, ok := IsAppError(regularErr); ok {
		t.Errorf("Expected IsAppError to return false for regular error")
	}
}
