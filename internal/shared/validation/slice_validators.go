package validation

import (
	"fmt"
)

// SliceValidator provides slice-specific validation functions
type SliceValidator struct{}

// NewSliceValidator creates a new slice validator
func NewSliceValidator() *SliceValidator {
	return &SliceValidator{}
}

// Required validates that a slice is not empty
func RequiredSlice[T any](field string, value []T) error {
	if len(value) == 0 {
		return ValidationError{Field: field, Message: "field is required", Value: value}
	}
	return nil
}

// MinLengthSlice creates a minimum length validator for slices
func MinLengthSlice[T any](min int) Validator[[]T] {
	return func(field string, value []T) error {
		if len(value) < min {
			return ValidationError{Field: field, Message: fmt.Sprintf("minimum length is %d", min), Value: value}
		}
		return nil
	}
}

// MaxLengthSlice creates a maximum length validator for slices
func MaxLengthSlice[T any](max int) Validator[[]T] {
	return func(field string, value []T) error {
		if len(value) > max {
			return ValidationError{Field: field, Message: fmt.Sprintf("maximum length is %d", max), Value: value}
		}
		return nil
	}
}

// Enum creates an enum validator for comparable types
func Enum[T comparable](values ...T) Validator[T] {
	valueSet := make(map[T]bool)
	for _, v := range values {
		valueSet[v] = true
	}

	return func(field string, value T) error {
		if !valueSet[value] {
			return ValidationError{Field: field, Message: fmt.Sprintf("value must be one of: %v", values), Value: value}
		}
		return nil
	}
}
