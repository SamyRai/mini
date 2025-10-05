package errors

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"
)

// This file has been refactored to use the new architecture.
// The functionality has been moved to:
// - internal/shared/errors/error_constructors.go
// - internal/shared/errors/error_aggregators.go
//
// This file is kept for backward compatibility but should be removed
// once all references are updated.

// ErrorCode represents different types of errors
type ErrorCode string

const (
	// Authentication errors
	ErrorCodeUnauthorized     ErrorCode = "UNAUTHORIZED"
	ErrorCodeInvalidAPIKey    ErrorCode = "INVALID_API_KEY"
	ErrorCodeRateLimitExceeded ErrorCode = "RATE_LIMIT_EXCEEDED"
	
	// Command execution errors
	ErrorCodeCommandNotFound  ErrorCode = "COMMAND_NOT_FOUND"
	ErrorCodeCommandTimeout   ErrorCode = "COMMAND_TIMEOUT"
	ErrorCodeCommandFailed    ErrorCode = "COMMAND_FAILED"
	ErrorCodeCommandBlocked   ErrorCode = "COMMAND_BLOCKED"
	
	// File system errors
	ErrorCodeFileNotFound     ErrorCode = "FILE_NOT_FOUND"
	ErrorCodePermissionDenied ErrorCode = "PERMISSION_DENIED"
	ErrorCodePathBlocked      ErrorCode = "PATH_BLOCKED"
	
	// Validation errors
	ErrorCodeInvalidInput     ErrorCode = "INVALID_INPUT"
	ErrorCodeMissingRequired  ErrorCode = "MISSING_REQUIRED"
	ErrorCodeInvalidFormat    ErrorCode = "INVALID_FORMAT"
	
	// System errors
	ErrorCodeInternalError    ErrorCode = "INTERNAL_ERROR"
	ErrorCodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
	ErrorCodeResourceExhausted ErrorCode = "RESOURCE_EXHAUSTED"
)

// ErrorResponse represents a structured error response
type ErrorResponse struct {
	Code        ErrorCode            `json:"code"`
	Message     string               `json:"message"`
	Details     map[string]any       `json:"details,omitempty"`
	Retryable   bool                 `json:"retryable"`
	RequestID   string               `json:"request_id,omitempty"`
	Timestamp   time.Time            `json:"timestamp"`
	Suggestions []string             `json:"suggestions,omitempty"`
	Stack       string               `json:"stack,omitempty"`
	Cause       error                `json:"-"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(code ErrorCode, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:      code,
		Message:   message,
		Details:   make(map[string]any),
		Retryable: isRetryable(code),
		Timestamp: time.Now(),
		Stack:     captureStackTrace(),
	}
}

// NewErrorResponseWithCause creates a new error response with a cause
func NewErrorResponseWithCause(code ErrorCode, message string, cause error) *ErrorResponse {
	return &ErrorResponse{
		Code:      code,
		Message:   message,
		Details:   make(map[string]any),
		Retryable: isRetryable(code),
		Timestamp: time.Now(),
		Stack:     captureStackTrace(),
		Cause:     cause,
	}
}

// captureStackTrace captures the current stack trace
func captureStackTrace() string {
	const depth = 10
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:]) // Skip runtime.Callers, captureStackTrace, and caller

	frames := runtime.CallersFrames(pcs[:n])
	var stack strings.Builder

	for {
		frame, more := frames.Next()
		stack.WriteString(fmt.Sprintf("%s:%d %s\n", frame.File, frame.Line, frame.Function))

		if !more {
			break
		}
	}

	return stack.String()
}

// WithDetails adds details to the error response
func (e *ErrorResponse) WithDetails(details map[string]any) *ErrorResponse {
	e.Details = details
	return e
}

// WithRequestID adds a request ID to the error response
func (e *ErrorResponse) WithRequestID(requestID string) *ErrorResponse {
	e.RequestID = requestID
	return e
}

// WithSuggestions adds suggestions to the error response
func (e *ErrorResponse) WithSuggestions(suggestions ...string) *ErrorResponse {
	e.Suggestions = suggestions
	return e
}

// Error implements the error interface
func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// ToJSON converts the error response to JSON
func (e *ErrorResponse) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// isRetryable determines if an error code is retryable
func isRetryable(code ErrorCode) bool {
	retryableCodes := map[ErrorCode]bool{
		ErrorCodeRateLimitExceeded: true,
		ErrorCodeCommandTimeout:    true,
		ErrorCodeServiceUnavailable: true,
		ErrorCodeResourceExhausted: true,
	}
	
	return retryableCodes[code]
}

// DEPRECATED: Use ErrorConstructors instead
var constructors = NewErrorConstructors()

// Common error constructors (delegated to new architecture)
func NewUnauthorizedError(message string) *ErrorResponse {
	return constructors.NewUnauthorizedError(message)
}

func NewInvalidAPIKeyError() *ErrorResponse {
	return constructors.NewInvalidAPIKeyError()
}

func NewRateLimitExceededError() *ErrorResponse {
	return constructors.NewRateLimitExceededError()
}

func NewCommandNotFoundError(command string) *ErrorResponse {
	return constructors.NewCommandNotFoundError(command)
}

func NewCommandTimeoutError(timeout time.Duration) *ErrorResponse {
	return constructors.NewCommandTimeoutError(timeout)
}

func NewCommandFailedError(command string, exitCode int, output string) *ErrorResponse {
	return constructors.NewCommandFailedError(command, exitCode, output)
}

func NewCommandBlockedError(command string, reason string) *ErrorResponse {
	return constructors.NewCommandBlockedError(command, reason)
}

func NewFileNotFoundError(path string) *ErrorResponse {
	return constructors.NewFileNotFoundError(path)
}

func NewPermissionDeniedError(path string) *ErrorResponse {
	return constructors.NewPermissionDeniedError(path)
}

func NewPathBlockedError(path string) *ErrorResponse {
	return constructors.NewPathBlockedError(path)
}

func NewFileAccessError(path string, operation string, details string) *ErrorResponse {
	return constructors.NewFileAccessError(path, operation, details)
}

func NewInvalidInputError(message string) *ErrorResponse {
	return constructors.NewInvalidInputError(message)
}

func NewMissingRequiredError(field string) *ErrorResponse {
	return constructors.NewMissingRequiredError(field)
}

func NewInvalidFormatError(field string, format string) *ErrorResponse {
	return constructors.NewInvalidFormatError(field, format)
}

func NewInternalError(message string) *ErrorResponse {
	return constructors.NewInternalError(message)
}

func NewServiceUnavailableError(service string) *ErrorResponse {
	return constructors.NewServiceUnavailableError(service)
}

func NewResourceExhaustedError(resource string) *ErrorResponse {
	return constructors.NewResourceExhaustedError(resource)
}

// WrapError wraps an existing error with additional context
func WrapError(err error, code ErrorCode, message string) *ErrorResponse {
	if err == nil {
		return NewErrorResponse(code, message)
	}

	return NewErrorResponseWithCause(code, message, err).
		WithSuggestions("Check logs for more details", "Contact support if issue persists")
}

// IsRetryable checks if an error is retryable
func IsRetryable(err error) bool {
	if errResp, ok := err.(*ErrorResponse); ok {
		return errResp.Retryable
	}
	return false
}

// GetErrorCode extracts error code from an error
func GetErrorCode(err error) ErrorCode {
	if errResp, ok := err.(*ErrorResponse); ok {
		return errResp.Code
	}
	return ErrorCodeInternalError
}

// GetErrorMessage extracts error message from an error
func GetErrorMessage(err error) string {
	if errResp, ok := err.(*ErrorResponse); ok {
		return errResp.Message
	}
	return err.Error()
}