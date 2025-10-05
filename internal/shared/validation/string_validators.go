package validation

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// StringValidator provides string-specific validation functions
type StringValidator struct{}

// NewStringValidator creates a new string validator
func NewStringValidator() *StringValidator {
	return &StringValidator{}
}

// Required validates that a string field is not empty
func (sv *StringValidator) Required(field string, value string) error {
	if strings.TrimSpace(value) == "" {
		return ValidationError{Field: field, Message: "field is required", Value: value}
	}
	return nil
}

// MinLength creates a minimum length validator for strings
func (sv *StringValidator) MinLength(min int) Validator[string] {
	return func(field string, value string) error {
		if len(strings.TrimSpace(value)) < min {
			return ValidationError{Field: field, Message: fmt.Sprintf("minimum length is %d", min), Value: value}
		}
		return nil
	}
}

// MaxLength creates a maximum length validator for strings
func (sv *StringValidator) MaxLength(max int) Validator[string] {
	return func(field string, value string) error {
		if len(strings.TrimSpace(value)) > max {
			return ValidationError{Field: field, Message: fmt.Sprintf("maximum length is %d", max), Value: value}
		}
		return nil
	}
}

// Pattern creates a pattern validator for strings
func (sv *StringValidator) Pattern(pattern string) Validator[string] {
	regex := regexp.MustCompile(pattern)
	return func(field string, value string) error {
		if value == "" {
			return nil // Skip validation for empty strings
		}

		if !regex.MatchString(value) {
			return ValidationError{Field: field, Message: fmt.Sprintf("value must match pattern: %s", pattern), Value: value}
		}
		return nil
	}
}

// URL validates that a field is a valid URL
func (sv *StringValidator) URL(field string, value string) error {
	if value == "" {
		return nil // Skip validation for empty strings
	}

	if _, err := url.ParseRequestURI(value); err != nil {
		return ValidationError{Field: field, Message: "value must be a valid URL", Value: value}
	}

	return nil
}

// Email validates that a field is a valid email address
func (sv *StringValidator) Email(field string, value string) error {
	if value == "" {
		return nil // Skip validation for empty strings
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		return ValidationError{Field: field, Message: "value must be a valid email address", Value: value}
	}

	return nil
}

// Host validates that a field is a valid hostname
func (sv *StringValidator) Host(field string, value string) error {
	if value == "" {
		return nil // Skip validation for empty strings
	}

	// Basic hostname validation
	hostRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`)
	if !hostRegex.MatchString(value) {
		return ValidationError{Field: field, Message: "value must be a valid hostname", Value: value}
	}

	return nil
}

// Port validates that a field is a valid port number (string representation)
func (sv *StringValidator) Port(field string, value string) error {
	if value == "" {
		return nil // Skip validation for empty strings
	}

	// Parse the port number
	var port int
	if _, err := fmt.Sscanf(value, "%d", &port); err != nil {
		return ValidationError{Field: field, Message: "port must be a valid number", Value: value}
	}

	// Validate port range
	if port < 1 || port > 65535 {
		return ValidationError{Field: field, Message: "port must be between 1 and 65535", Value: port}
	}

	return nil
}

// Path validates that a field is a safe file path
func (sv *StringValidator) Path(field string, value string) error {
	if value == "" {
		return nil // Skip validation for empty strings
	}

	// Basic path validation - more comprehensive validation is done in security layer
	if strings.Contains(value, "..") {
		return ValidationError{Field: field, Message: "path contains forbidden directory traversal", Value: value}
	}

	return nil
}
