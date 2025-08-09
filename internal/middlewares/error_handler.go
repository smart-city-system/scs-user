package middleware

import (
	"net/http"
	"scs-user/pkg/errors"
	"scs-user/pkg/utils"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// ErrorHandlerMiddleware handles all errors in a centralized way
func (mw *MiddlewareManager) ErrorHandlerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err == nil {
			return nil
		}

		requestID := utils.GetRequestID(c)

		// Log the error
		mw.logger.Errorf("Request %s failed: %v", requestID, err)

		// Handle different types of errors
		if appErr, ok := errors.IsAppError(err); ok {
			return mw.handleAppError(c, appErr, requestID)
		}

		// Handle Echo HTTP errors
		if httpErr, ok := err.(*echo.HTTPError); ok {
			return mw.handleEchoHTTPError(c, httpErr, requestID)
		}

		// Handle GORM errors
		if mw.isGormError(err) {
			return mw.handleGormError(c, err, requestID)
		}

		// Handle unknown errors
		return mw.handleUnknownError(c, err, requestID)
	}
}

// handleAppError handles custom application errors
func (mw *MiddlewareManager) handleAppError(c echo.Context, appErr *errors.AppError, requestID string) error {
	response := errors.NewErrorResponse(appErr, requestID)
	return c.JSON(appErr.StatusCode, response)
}

// handleEchoHTTPError handles Echo framework HTTP errors
func (mw *MiddlewareManager) handleEchoHTTPError(c echo.Context, httpErr *echo.HTTPError, requestID string) error {
	var errorType errors.ErrorType
	var message string

	switch httpErr.Code {
	case http.StatusBadRequest:
		errorType = errors.ErrorTypeBadRequest
		message = "Bad request"
	case http.StatusUnauthorized:
		errorType = errors.ErrorTypeUnauthorized
		message = "Unauthorized"
	case http.StatusForbidden:
		errorType = errors.ErrorTypeForbidden
		message = "Forbidden"
	case http.StatusNotFound:
		errorType = errors.ErrorTypeNotFound
		message = "Resource not found"
	case http.StatusConflict:
		errorType = errors.ErrorTypeConflict
		message = "Conflict"
	default:
		errorType = errors.ErrorTypeInternal
		message = "Internal server error"
	}

	if httpErr.Message != nil {
		if msg, ok := httpErr.Message.(string); ok {
			message = msg
		}
	}

	appErr := errors.NewAppError(errorType, message, httpErr)
	response := errors.NewErrorResponse(appErr, requestID)
	return c.JSON(httpErr.Code, response)
}

// handleGormError handles GORM database errors
func (mw *MiddlewareManager) handleGormError(c echo.Context, err error, requestID string) error {
	var appErr *errors.AppError

	switch err {
	case gorm.ErrRecordNotFound:
		appErr = errors.NewNotFoundError("record")
	case gorm.ErrInvalidTransaction:
		appErr = errors.NewDatabaseError("invalid transaction", err)
	case gorm.ErrNotImplemented:
		appErr = errors.NewInternalError("database operation not implemented", err)
	case gorm.ErrMissingWhereClause:
		appErr = errors.NewBadRequestError("missing where clause in query")
	case gorm.ErrUnsupportedRelation:
		appErr = errors.NewInternalError("unsupported database relation", err)
	case gorm.ErrPrimaryKeyRequired:
		appErr = errors.NewBadRequestError("primary key required")
	case gorm.ErrModelValueRequired:
		appErr = errors.NewBadRequestError("model value required")
	case gorm.ErrInvalidData:
		appErr = errors.NewBadRequestError("invalid data provided")
	case gorm.ErrUnsupportedDriver:
		appErr = errors.NewInternalError("unsupported database driver", err)
	case gorm.ErrRegistered:
		appErr = errors.NewInternalError("database callback already registered", err)
	case gorm.ErrInvalidField:
		appErr = errors.NewBadRequestError("invalid field in query")
	case gorm.ErrEmptySlice:
		appErr = errors.NewBadRequestError("empty slice provided")
	case gorm.ErrDryRunModeUnsupported:
		appErr = errors.NewInternalError("dry run mode unsupported", err)
	case gorm.ErrInvalidDB:
		appErr = errors.NewInternalError("invalid database connection", err)
	case gorm.ErrInvalidValue:
		appErr = errors.NewBadRequestError("invalid value provided")
	case gorm.ErrInvalidValueOfLength:
		appErr = errors.NewBadRequestError("invalid value length")
	default:
		// Check for constraint violations and other database-specific errors
		if mw.isDuplicateKeyError(err) {
			appErr = errors.NewConflictError("resource already exists")
		} else if mw.isForeignKeyError(err) {
			appErr = errors.NewBadRequestError("invalid reference to related resource")
		} else {
			appErr = errors.NewDatabaseError("unknown database error", err)
		}
	}

	response := errors.NewErrorResponse(appErr, requestID)
	return c.JSON(appErr.StatusCode, response)
}

// handleUnknownError handles any other unknown errors
func (mw *MiddlewareManager) handleUnknownError(c echo.Context, err error, requestID string) error {
	appErr := errors.NewInternalError("An unexpected error occurred", err)
	response := errors.NewErrorResponse(appErr, requestID)
	return c.JSON(http.StatusInternalServerError, response)
}

// isGormError checks if the error is a GORM error
func (mw *MiddlewareManager) isGormError(err error) bool {
	gormErrors := []error{
		gorm.ErrRecordNotFound,
		gorm.ErrInvalidTransaction,
		gorm.ErrNotImplemented,
		gorm.ErrMissingWhereClause,
		gorm.ErrUnsupportedRelation,
		gorm.ErrPrimaryKeyRequired,
		gorm.ErrModelValueRequired,
		gorm.ErrInvalidData,
		gorm.ErrUnsupportedDriver,
		gorm.ErrRegistered,
		gorm.ErrInvalidField,
		gorm.ErrEmptySlice,
		gorm.ErrDryRunModeUnsupported,
		gorm.ErrInvalidDB,
		gorm.ErrInvalidValue,
		gorm.ErrInvalidValueOfLength,
	}

	for _, gormErr := range gormErrors {
		if err == gormErr {
			return true
		}
	}

	// Check for database-specific errors
	return mw.isDuplicateKeyError(err) || mw.isForeignKeyError(err)
}

// isDuplicateKeyError checks if the error is a duplicate key constraint violation
func (mw *MiddlewareManager) isDuplicateKeyError(err error) bool {
	errStr := err.Error()
	// PostgreSQL duplicate key error patterns
	return contains(errStr, "duplicate key value violates unique constraint") ||
		contains(errStr, "UNIQUE constraint failed") ||
		contains(errStr, "duplicate entry")
}

// isForeignKeyError checks if the error is a foreign key constraint violation
func (mw *MiddlewareManager) isForeignKeyError(err error) bool {
	errStr := err.Error()
	// PostgreSQL foreign key error patterns
	return contains(errStr, "violates foreign key constraint") ||
		contains(errStr, "FOREIGN KEY constraint failed") ||
		contains(errStr, "foreign key constraint fails")
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
