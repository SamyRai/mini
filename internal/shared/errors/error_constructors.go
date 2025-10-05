package errors

import (
	"fmt"
	"time"
)

// ErrorConstructors provides common error constructors
type ErrorConstructors struct{}

// NewErrorConstructors creates a new error constructors instance
func NewErrorConstructors() *ErrorConstructors {
	return &ErrorConstructors{}
}

// Authentication errors
func (ec *ErrorConstructors) NewUnauthorizedError(message string) *ErrorResponse {
	return NewErrorResponse(ErrorCodeUnauthorized, message)
}

func (ec *ErrorConstructors) NewInvalidAPIKeyError() *ErrorResponse {
	return NewErrorResponse(ErrorCodeInvalidAPIKey, "Invalid or missing API key")
}

func (ec *ErrorConstructors) NewRateLimitExceededError() *ErrorResponse {
	return NewErrorResponse(ErrorCodeRateLimitExceeded, "Rate limit exceeded")
}

// Command execution errors
func (ec *ErrorConstructors) NewCommandNotFoundError(command string) *ErrorResponse {
	return NewErrorResponse(ErrorCodeCommandNotFound, fmt.Sprintf("Command not found: %s", command))
}

func (ec *ErrorConstructors) NewCommandTimeoutError(timeout time.Duration) *ErrorResponse {
	return NewErrorResponse(ErrorCodeCommandTimeout, fmt.Sprintf("Command timed out after %v", timeout))
}

func (ec *ErrorConstructors) NewCommandFailedError(command string, exitCode int, output string) *ErrorResponse {
	return NewErrorResponse(ErrorCodeCommandFailed, fmt.Sprintf("Command failed: %s", command)).
		WithDetails(map[string]any{
			"command":   command,
			"exit_code": exitCode,
			"output":    output,
		}).
		WithSuggestions("Check command syntax", "Verify command permissions", "Check system resources")
}

func (ec *ErrorConstructors) NewCommandBlockedError(command string, reason string) *ErrorResponse {
	return NewErrorResponse(ErrorCodeCommandBlocked, fmt.Sprintf("Command blocked: %s", command)).
		WithDetails(map[string]any{
			"command": command,
			"reason":  reason,
		})
}

// File system errors
func (ec *ErrorConstructors) NewFileNotFoundError(path string) *ErrorResponse {
	return NewErrorResponse(ErrorCodeFileNotFound, fmt.Sprintf("File not found: %s", path))
}

func (ec *ErrorConstructors) NewPermissionDeniedError(path string) *ErrorResponse {
	return NewErrorResponse(ErrorCodePermissionDenied, fmt.Sprintf("Permission denied: %s", path))
}

func (ec *ErrorConstructors) NewPathBlockedError(path string) *ErrorResponse {
	return NewErrorResponse(ErrorCodePathBlocked, fmt.Sprintf("Path blocked: %s", path))
}

func (ec *ErrorConstructors) NewFileAccessError(path string, operation string, details string) *ErrorResponse {
	return NewErrorResponse(ErrorCodePermissionDenied, fmt.Sprintf("File access denied: %s operation on %s", operation, path)).
		WithDetails(map[string]any{
			"path":      path,
			"operation": operation,
			"details":   details,
		}).
		WithSuggestions("Check user permissions")
}

// Validation errors
func (ec *ErrorConstructors) NewInvalidInputError(message string) *ErrorResponse {
	return NewErrorResponse(ErrorCodeInvalidInput, message)
}

func (ec *ErrorConstructors) NewMissingRequiredError(field string) *ErrorResponse {
	return NewErrorResponse(ErrorCodeMissingRequired, fmt.Sprintf("Required field missing: %s", field))
}

func (ec *ErrorConstructors) NewInvalidFormatError(field string, format string) *ErrorResponse {
	return NewErrorResponse(ErrorCodeInvalidFormat, fmt.Sprintf("Invalid format for field %s: expected %s", field, format))
}

// System errors
func (ec *ErrorConstructors) NewInternalError(message string) *ErrorResponse {
	return NewErrorResponse(ErrorCodeInternalError, message)
}

func (ec *ErrorConstructors) NewServiceUnavailableError(service string) *ErrorResponse {
	return NewErrorResponse(ErrorCodeServiceUnavailable, fmt.Sprintf("Service unavailable: %s", service))
}

func (ec *ErrorConstructors) NewResourceExhaustedError(resource string) *ErrorResponse {
	return NewErrorResponse(ErrorCodeResourceExhausted, fmt.Sprintf("Resource exhausted: %s", resource))
}
