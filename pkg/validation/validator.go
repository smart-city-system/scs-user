package validation

import (
	"fmt"
	"reflect"
	"scs-user/pkg/errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wraps the go-playground validator
type Validator struct {
	validator *validator.Validate
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	v := validator.New()

	// Register custom tag name function to use json tags
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register custom validators
	registerCustomValidators(v)

	return &Validator{validator: v}
}

// Validate validates a struct and returns validation errors
func (v *Validator) Validate(s interface{}) error {
	err := v.validator.Struct(s)
	if err == nil {
		return nil
	}

	var validationErrors errors.ValidationErrors

	if validatorErrors, ok := err.(validator.ValidationErrors); ok {
		for _, validatorError := range validatorErrors {
			validationError := errors.ValidationError{
				Field:   validatorError.Field(),
				Message: getErrorMessage(validatorError),
				Value:   validatorError.Value(),
			}
			validationErrors = append(validationErrors, validationError)
		}
	}

	return errors.NewValidationError("Validation failed", validationErrors)
}

// getErrorMessage returns a human-readable error message for validation errors
func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", fe.Field(), fe.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", fe.Field(), fe.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", fe.Field(), fe.Param())
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", fe.Field())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", fe.Field())
	case "numeric":
		return fmt.Sprintf("%s must be numeric", fe.Field())
	case "alpha":
		return fmt.Sprintf("%s must contain only letters", fe.Field())
	case "alphanum":
		return fmt.Sprintf("%s must contain only letters and numbers", fe.Field())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", fe.Field(), fe.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", fe.Field(), fe.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", fe.Field(), fe.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("%s is invalid", fe.Field())
	}
}

// registerCustomValidators registers custom validation rules
func registerCustomValidators(v *validator.Validate) {
	// Register UUID validator
	v.RegisterValidation("uuid", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		if value == "" {
			return true // Let required handle empty values
		}
		// Simple UUID format validation
		return len(value) == 36 &&
			value[8] == '-' &&
			value[13] == '-' &&
			value[18] == '-' &&
			value[23] == '-'
	})

	// Register role validator for user roles
	v.RegisterValidation("role", func(fl validator.FieldLevel) bool {
		role := fl.Field().String()
		validRoles := []string{"admin", "guard", "operator"}
		for _, validRole := range validRoles {
			if role == validRole {
				return true
			}
		}
		return false
	})
}

// ValidateStruct is a convenience function to validate a struct
func ValidateStruct(s interface{}) error {
	validator := NewValidator()
	return validator.Validate(s)
}
