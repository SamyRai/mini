package validation

import "fmt"

// Validator represents a validation function
type Validator[T any] func(field string, value T) error

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
	Value   any
}

func (e ValidationError) Error() string {
	return e.Message
}

// This file has been refactored to use the new architecture.
// The functionality has been moved to:
// - internal/shared/validation/string_validators.go
// - internal/shared/validation/numeric_validators.go
// - internal/shared/validation/slice_validators.go
// - internal/shared/validation/validation_factory.go
//
// This file is kept for backward compatibility but should be removed
// once all references are updated.

// DEPRECATED: Use ValidationFactory instead
var factory = NewValidationFactory()

// Required validates that a field is not empty
func Required[T any](field string, value T) error {
	return factory.Required(field, value)
}

// MinLength creates a minimum length validator for strings
func MinLength(min int) Validator[string] {
	return factory.MinLength(min)
}

// MaxLength creates a maximum length validator for strings
func MaxLength(max int) Validator[string] {
	return factory.MaxLength(max)
}

// RangeInt creates a range validator for integers
func RangeInt(min, max int) Validator[int] {
	return factory.RangeInt(min, max)
}

// RangeFloat creates a range validator for floats
func RangeFloat(min, max float64) Validator[float64] {
	return factory.RangeFloat(min, max)
}

// Pattern creates a pattern validator for strings
func Pattern(pattern string) Validator[string] {
	return factory.Pattern(pattern)
}

// EnumString creates an enum validator for strings
func EnumString(values ...string) Validator[string] {
	return factory.EnumString(values...)
}

// URL validates that a field is a valid URL
func URL(field string, value string) error {
	return factory.URL(field, value)
}

// Email validates that a field is a valid email address
func Email(field string, value string) error {
	return factory.Email(field, value)
}

// Host validates that a field is a valid hostname
func Host(field string, value string) error {
	return factory.Host(field, value)
}

// Port validates that a field is a valid port number (string representation)
func Port(field string, value string) error {
	return factory.Port(field, value)
}

// PortInt validates that an integer field is a valid port number
func PortInt(field string, value int) error {
	return factory.PortInt(field, value)
}

// Timeout validates timeout values
func Timeout(field string, value int) error {
	return factory.Timeout(field, value)
}

// Path validates that a field is a safe file path
func Path(field string, value string) error {
	return factory.Path(field, value)
}

// Positive validates that an integer is positive
func Positive(field string, value int) error {
	return factory.Positive(field, value)
}

// NewInvalidFormatError creates a new invalid format error
func NewInvalidFormatError(field string, format string) error {
	return ValidationError{
		Field:   field,
		Message: fmt.Sprintf("invalid format: %s", format),
		Value:   nil,
	}
}

// NewMissingRequiredError creates a new missing required field error
func NewMissingRequiredError(field string) error {
	return ValidationError{
		Field:   field,
		Message: "field is required",
		Value:   nil,
	}
}

// StringRequired validates that a string field is not empty
func StringRequired(field string, value string) error {
	if value == "" {
		return ValidationError{
			Field:   field,
			Message: "field is required",
			Value:   value,
		}
	}
	return nil
}

// StringPath validates that a string is a safe file path
func StringPath(field string, value string) error {
	return Path(field, value)
}

// ValidateProxmoxArgs validates Proxmox connection arguments
func ValidateProxmoxArgs(host, user, password string) struct {
	IsValid bool
	Errors  []error
} {
	var errors []error
	
	if host == "" {
		errors = append(errors, ValidationError{Field: "host", Message: "host is required"})
	}
	if user == "" {
		errors = append(errors, ValidationError{Field: "user", Message: "user is required"})
	}
	if password == "" {
		errors = append(errors, ValidationError{Field: "password", Message: "password is required"})
	}
	
	return struct {
		IsValid bool
		Errors  []error
	}{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}