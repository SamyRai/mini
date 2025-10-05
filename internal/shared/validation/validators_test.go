package validation

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidationError_Error(t *testing.T) {
	err := ValidationError{
		Field:   "test_field",
		Message: "test message",
		Value:   "test_value",
	}

	assert.Equal(t, "test message", err.Error())
}

func TestValidationError_Fields(t *testing.T) {
	err := ValidationError{
		Field:   "test_field",
		Message: "test message",
		Value:   "test_value",
	}
	assert.Equal(t, "test_field", err.Field)
	assert.Equal(t, "test message", err.Message)
	assert.Equal(t, "test_value", err.Value)
}

func TestValidationFactory_DurationPositive(t *testing.T) {
	vf := NewValidationFactory()

	// Test valid duration
	err := vf.DurationPositive("timeout", time.Second*30)
	assert.NoError(t, err)

	// Test invalid duration
	err = vf.DurationPositive("timeout", time.Second*-1)
	assert.Error(t, err)
}

func TestValidationFactory_RangeInt64(t *testing.T) {
	vf := NewValidationFactory()

	validator := vf.RangeInt64(10, 100)

	// Test valid value
	err := validator("test_field", int64(50))
	assert.NoError(t, err)

	// Test value too low
	err = validator("test_field", int64(5))
	assert.Error(t, err)

	// Test value too high
	err = validator("test_field", int64(150))
	assert.Error(t, err)
}

func TestValidationFactory_Positive(t *testing.T) {
	vf := NewValidationFactory()

	// Test valid positive value
	err := vf.Positive("count", 10)
	assert.NoError(t, err)

	// Test invalid negative value
	err = vf.Positive("count", -1)
	assert.Error(t, err)

	// Test zero
	err = vf.Positive("count", 0)
	assert.Error(t, err)
}
