package validation

import (
	"fmt"
	"time"
)

// NumericValidator provides numeric-specific validation functions
type NumericValidator struct{}

// NewNumericValidator creates a new numeric validator
func NewNumericValidator() *NumericValidator {
	return &NumericValidator{}
}

// Range creates a range validator for numeric types
func Range[T interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}](min, max T) Validator[T] {
	return func(field string, value T) error {
		if value < min || value > max {
			return ValidationError{Field: field, Message: fmt.Sprintf("value must be between %v and %v", min, max), Value: value}
		}
		return nil
	}
}

// RangeInt creates a range validator for integers
func (nv *NumericValidator) RangeInt(min, max int) Validator[int] {
	return func(field string, value int) error {
		if value < min || value > max {
			return ValidationError{Field: field, Message: fmt.Sprintf("value must be between %d and %d", min, max), Value: value}
		}
		return nil
	}
}

// RangeFloat creates a range validator for floats
func (nv *NumericValidator) RangeFloat(min, max float64) Validator[float64] {
	return func(field string, value float64) error {
		if value < min || value > max {
			return ValidationError{Field: field, Message: fmt.Sprintf("value must be between %v and %v", min, max), Value: value}
		}
		return nil
	}
}

// PortInt validates that an integer field is a valid port number
func (nv *NumericValidator) PortInt(field string, value int) error {
	if value < 1 || value > 65535 {
		return ValidationError{Field: field, Message: "port must be between 1 and 65535", Value: value}
	}
	return nil
}

// Timeout validates timeout values
func (nv *NumericValidator) Timeout(field string, value int) error {
	if value < 1 || value > 300 {
		return ValidationError{Field: field, Message: "timeout must be between 1 and 300 seconds", Value: value}
	}
	return nil
}

// Positive validates that an integer is positive
func (nv *NumericValidator) Positive(field string, value int) error {
	if value <= 0 {
		return ValidationError{Field: field, Message: "value must be positive", Value: value}
	}
	return nil
}

// DurationPositive validates that a duration is positive
func (nv *NumericValidator) DurationPositive(field string, value time.Duration) error {
	if value <= 0 {
		return ValidationError{Field: field, Message: "duration must be positive", Value: value}
	}
	return nil
}
