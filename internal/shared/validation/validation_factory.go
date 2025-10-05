package validation

import (
	"fmt"
	"time"
)

// ValidationFactory provides a factory for creating validators
type ValidationFactory struct {
	stringValidator *StringValidator
	numericValidator *NumericValidator
	sliceValidator *SliceValidator
}

// NewValidationFactory creates a new validation factory
func NewValidationFactory() *ValidationFactory {
	return &ValidationFactory{
		stringValidator:  NewStringValidator(),
		numericValidator: NewNumericValidator(),
		sliceValidator:    NewSliceValidator(),
	}
}

// String returns the string validator
func (vf *ValidationFactory) String() *StringValidator {
	return vf.stringValidator
}

// Numeric returns the numeric validator
func (vf *ValidationFactory) Numeric() *NumericValidator {
	return vf.numericValidator
}

// Slice returns the slice validator
func (vf *ValidationFactory) Slice() *SliceValidator {
	return vf.sliceValidator
}

// Common validation functions that delegate to specific validators
func (vf *ValidationFactory) Required(field string, value any) error {
	return Required(field, value)
}

func (vf *ValidationFactory) MinLength(min int) Validator[string] {
	return vf.stringValidator.MinLength(min)
}

func (vf *ValidationFactory) MaxLength(max int) Validator[string] {
	return vf.stringValidator.MaxLength(max)
}

func (vf *ValidationFactory) RangeInt(min, max int) Validator[int] {
	return vf.numericValidator.RangeInt(min, max)
}

func (vf *ValidationFactory) RangeFloat(min, max float64) Validator[float64] {
	return vf.numericValidator.RangeFloat(min, max)
}

func (vf *ValidationFactory) Pattern(pattern string) Validator[string] {
	return vf.stringValidator.Pattern(pattern)
}

func (vf *ValidationFactory) EnumString(values ...string) Validator[string] {
	return Enum(values...)
}

func (vf *ValidationFactory) URL(field string, value string) error {
	return vf.stringValidator.URL(field, value)
}

func (vf *ValidationFactory) Email(field string, value string) error {
	return vf.stringValidator.Email(field, value)
}

func (vf *ValidationFactory) Host(field string, value string) error {
	return vf.stringValidator.Host(field, value)
}

func (vf *ValidationFactory) Port(field string, value string) error {
	return vf.stringValidator.Port(field, value)
}

func (vf *ValidationFactory) PortInt(field string, value int) error {
	return vf.numericValidator.PortInt(field, value)
}

func (vf *ValidationFactory) Timeout(field string, value int) error {
	return vf.numericValidator.Timeout(field, value)
}

func (vf *ValidationFactory) DurationPositive(field string, value time.Duration) error {
	if value <= 0 {
		return ValidationError{
			Field:   field,
			Message: "duration must be positive",
			Value:   value,
		}
	}
	return nil
}

func (vf *ValidationFactory) Path(field string, value string) error {
	return vf.stringValidator.Path(field, value)
}

func (vf *ValidationFactory) Positive(field string, value int) error {
	return vf.numericValidator.Positive(field, value)
}

func (vf *ValidationFactory) RangeInt64(min, max int64) Validator[int64] {
	return func(field string, value int64) error {
		if value < min || value > max {
			return ValidationError{
				Field:   field,
				Message: fmt.Sprintf("value must be between %d and %d", min, max),
				Value:   value,
			}
		}
		return nil
	}
}

func (vf *ValidationFactory) NewStringValidator() *StringValidator {
	return vf.stringValidator
}

func (vf *ValidationFactory) NewInt64Validator() *NumericValidator {
	return vf.numericValidator
}
