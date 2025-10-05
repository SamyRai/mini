package errors

import (
	"fmt"
)

// ErrorAggregators provides error aggregation functionality
type ErrorAggregators struct{}

// NewErrorAggregators creates a new error aggregators instance
func NewErrorAggregators() *ErrorAggregators {
	return &ErrorAggregators{}
}

// ErrorAggregator aggregates multiple errors
type ErrorAggregator struct {
	errors []*ErrorResponse
}

// NewErrorAggregator creates a new error aggregator
func (ea *ErrorAggregators) NewErrorAggregator() *ErrorAggregator {
	return &ErrorAggregator{
		errors: make([]*ErrorResponse, 0),
	}
}

// Add adds an error to the aggregator
func (ea *ErrorAggregator) Add(err *ErrorResponse) {
	ea.errors = append(ea.errors, err)
}

// HasErrors returns true if there are any errors
func (ea *ErrorAggregator) HasErrors() bool {
	return len(ea.errors) > 0
}

// GetErrors returns all collected errors
func (ea *ErrorAggregator) GetErrors() []*ErrorResponse {
	return ea.errors
}

// ToCombinedError creates a combined error response
func (ea *ErrorAggregator) ToCombinedError() *ErrorResponse {
	if len(ea.errors) == 0 {
		return nil
	}

	if len(ea.errors) == 1 {
		return ea.errors[0]
	}

	// Create a combined error
	combined := NewErrorResponse(ErrorCodeInternalError, fmt.Sprintf("Multiple errors occurred (%d total)", len(ea.errors)))
	combined.Details["error_count"] = len(ea.errors)
	combined.Details["errors"] = ea.errors

	return combined
}

// ErrorCollector collects and manages multiple errors
type ErrorCollector struct {
	errors []*ErrorResponse
}

// NewErrorCollector creates a new error collector
func (ea *ErrorAggregators) NewErrorCollector() *ErrorCollector {
	return &ErrorCollector{
		errors: make([]*ErrorResponse, 0),
	}
}

// Collect adds an error to the collection
func (ec *ErrorCollector) Collect(err error) {
	if err == nil {
		return
	}

	if errResp, ok := err.(*ErrorResponse); ok {
		ec.errors = append(ec.errors, errResp)
	} else {
		// Wrap generic error
		ec.errors = append(ec.errors, WrapError(err, ErrorCodeInternalError, "An unexpected error occurred"))
	}
}

// HasErrors returns true if any errors were collected
func (ec *ErrorCollector) HasErrors() bool {
	return len(ec.errors) > 0
}

// GetErrors returns all collected errors
func (ec *ErrorCollector) GetErrors() []*ErrorResponse {
	return ec.errors
}

// GetCombinedError returns a combined error if there are multiple errors
func (ec *ErrorCollector) GetCombinedError() error {
	if len(ec.errors) == 0 {
		return nil
	}

	if len(ec.errors) == 1 {
		return ec.errors[0]
	}

	// Create a combined error
	combined := NewErrorResponse(ErrorCodeInternalError, fmt.Sprintf("Multiple errors occurred (%d total)", len(ec.errors)))
	combined.Details["error_count"] = len(ec.errors)
	combined.Details["errors"] = ec.errors

	return combined
}
