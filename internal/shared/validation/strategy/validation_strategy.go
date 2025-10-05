package strategy

import (
	"fmt"
	"strings"
)

// ValidationStrategy defines the interface for validation strategies
// This follows the Strategy pattern with type safety using generics
type ValidationStrategy[T any] interface {
	Validate(field string, value T) error
	GetName() string
}

// StringValidationStrategy handles string validation
type StringValidationStrategy struct {
	name string
}

// NewStringValidationStrategy creates a new string validation strategy
func NewStringValidationStrategy() ValidationStrategy[string] {
	return &StringValidationStrategy{
		name: "string",
	}
}

// Validate validates string values
func (s *StringValidationStrategy) Validate(field string, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("field %s is required", field)
	}

	return nil
}

// GetName returns the strategy name
func (s *StringValidationStrategy) GetName() string {
	return s.name
}

// NumericValidationStrategy handles numeric validation for specific numeric types
type NumericValidationStrategy[T interface{ ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 }] struct {
	name string
	min  *T
	max  *T
}

// NewNumericValidationStrategy creates a new numeric validation strategy
func NewNumericValidationStrategy[T interface{ ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 }](min, max *T) ValidationStrategy[T] {
	return &NumericValidationStrategy[T]{
		name: "numeric",
		min:  min,
		max:  max,
	}
}

// Validate validates numeric values with type safety
func (s *NumericValidationStrategy[T]) Validate(field string, value T) error {
	if s.min != nil && value < *s.min {
		return fmt.Errorf("field %s must be >= %v", field, *s.min)
	}
	if s.max != nil && value > *s.max {
		return fmt.Errorf("field %s must be <= %v", field, *s.max)
	}

	return nil
}

// GetName returns the strategy name
func (s *NumericValidationStrategy[T]) GetName() string {
	return s.name
}

// SliceValidationStrategy handles slice validation with type safety
type SliceValidationStrategy[T any] struct {
	name   string
	minLen int
	maxLen int
}

// NewSliceValidationStrategy creates a new slice validation strategy
func NewSliceValidationStrategy[T any](minLen, maxLen int) ValidationStrategy[[]T] {
	return &SliceValidationStrategy[T]{
		name:   "slice",
		minLen: minLen,
		maxLen: maxLen,
	}
}

// Validate validates slice values with type safety
func (s *SliceValidationStrategy[T]) Validate(field string, value []T) error {
	length := len(value)
	if length < s.minLen {
		return fmt.Errorf("field %s must have at least %d items", field, s.minLen)
	}
	if length > s.maxLen {
		return fmt.Errorf("field %s must have at most %d items", field, s.maxLen)
	}

	return nil
}

// GetName returns the strategy name
func (s *SliceValidationStrategy[T]) GetName() string {
	return s.name
}

// CompositeValidationStrategy combines multiple strategies of the same type
type CompositeValidationStrategy[T any] struct {
	name       string
	strategies []ValidationStrategy[T]
}

// NewCompositeValidationStrategy creates a new composite validation strategy
func NewCompositeValidationStrategy[T any](name string, strategies ...ValidationStrategy[T]) ValidationStrategy[T] {
	return &CompositeValidationStrategy[T]{
		name:       name,
		strategies: strategies,
	}
}

// Validate validates using all strategies
func (s *CompositeValidationStrategy[T]) Validate(field string, value T) error {
	for _, strategy := range s.strategies {
		if err := strategy.Validate(field, value); err != nil {
			return err
		}
	}
	return nil
}

// GetName returns the strategy name
func (s *CompositeValidationStrategy[T]) GetName() string {
	return s.name
}

