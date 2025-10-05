package security

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandValidator_ValidateCommand(t *testing.T) {
	config := &SecurityConfig{
		AllowedCommands: []string{"ls", "cat", "echo"},
		BlockedPaths:    []string{"/etc/passwd", "/root"},
	}

	validator := NewCommandValidator(config)

	tests := []struct {
		name      string
		command   string
		wantError bool
		errorCode string
	}{
		{
			name:      "valid command",
			command:   "ls -la",
			wantError: false,
		},
		{
			name:      "empty command",
			command:   "",
			wantError: true,
			errorCode: ErrCodeInvalidInput,
		},
		{
			name:      "dangerous pattern - rm -rf",
			command:   "rm -rf /",
			wantError: true,
			errorCode: ErrCodeDangerousPattern,
		},
		{
			name:      "dangerous pattern - dd",
			command:   "dd if=/dev/zero of=/dev/sda",
			wantError: true,
			errorCode: ErrCodeDangerousPattern,
		},
		{
			name:      "dangerous pattern - chmod 777",
			command:   "chmod 777 /etc/shadow",
			wantError: true,
			errorCode: ErrCodeDangerousPattern,
		},
		{
			name:      "path traversal",
			command:   "ls ../../../etc",
			wantError: true,
			errorCode: ErrCodePathTraversal,
		},
		{
			name:      "blocked path",
			command:   "cat /etc/passwd",
			wantError: true,
			errorCode: ErrCodePathBlocked,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCommand(tt.command)

			if tt.wantError {
				require.Error(t, err)
				var secErr *SecurityError
				if errors.As(err, &secErr) {
					assert.Equal(t, tt.errorCode, secErr.Code)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCommandValidator_IsCommandAllowed(t *testing.T) {
	config := &SecurityConfig{
		AllowedCommands: []string{"ls", "cat", "echo"},
	}

	validator := NewCommandValidator(config)

	tests := []struct {
		name     string
		command  string
		expected bool
	}{
		{"allowed command", "ls", true},
		{"allowed command with args", "ls -la", true},
		{"not allowed command", "rm", false},
		{"empty command", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.IsCommandAllowed(tt.command)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPathValidator_ValidatePath(t *testing.T) {
	config := &SecurityConfig{
		AllowedPaths: []string{"/tmp", "/var/log"},
		BlockedPaths: []string{"/etc", "/root"},
	}

	validator := NewPathValidator(config)

	tests := []struct {
		name      string
		path      string
		wantError bool
		errorCode string
	}{
		{
			name:      "allowed path",
			path:      "/tmp/test",
			wantError: false,
		},
		{
			name:      "empty path",
			path:      "",
			wantError: true,
			errorCode: ErrCodeInvalidInput,
		},
		{
			name:      "blocked path",
			path:      "/etc/passwd",
			wantError: true,
			errorCode: ErrCodePathBlocked,
		},
		{
			name:      "path traversal",
			path:      "/tmp/../../../etc",
			wantError: true,
			errorCode: ErrCodePathTraversal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidatePath(tt.path)

			if tt.wantError {
				require.Error(t, err)
				var secErr *SecurityError
				if errors.As(err, &secErr) {
					assert.Equal(t, tt.errorCode, secErr.Code)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestInputSanitizer_Sanitize(t *testing.T) {
	sanitizer := NewInputSanitizer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal input",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "input with null bytes",
			input:    "hello\x00world",
			expected: "helloworld",
		},
		{
			name:     "input with control characters",
			input:    "hello\x01\x02world\x7f",
			expected: "helloworld",
		},
		{
			name:     "long input gets truncated",
			input:    strings.Repeat("a", 15000),
			expected: strings.Repeat("a", 10000),
		},
		{
			name:     "input with newlines and tabs preserved",
			input:    "hello\nworld\ttab",
			expected: "hello\nworld\ttab",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizer.Sanitize(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSecurityError_Error(t *testing.T) {
	err := SecurityError{
		Code:    "TEST_ERROR",
		Message: "Test error message",
		Cause:   assert.AnError,
	}

	expected := "security error [TEST_ERROR]: Test error message: assert.AnError general error for testing"
	assert.Equal(t, expected, err.Error())
}

func TestSecurityError_Unwrap(t *testing.T) {
	cause := assert.AnError
	err := SecurityError{
		Code:    "TEST_ERROR",
		Message: "Test error message",
		Cause:   cause,
	}

	assert.Equal(t, cause, err.Unwrap())
}

func TestDefaultSecurityConfig(t *testing.T) {
	config := DefaultSecurityConfig()

	assert.NotNil(t, config)
	assert.NotEmpty(t, config.AllowedCommands)
	assert.NotEmpty(t, config.WorkingDirectory)
	assert.Greater(t, config.CommandTimeout, time.Duration(0))
	assert.Greater(t, config.MaxOutputSize, int64(0))
	assert.NotEmpty(t, config.AllowedEnvVars)
	assert.NotEmpty(t, config.AllowedPaths)
	assert.NotEmpty(t, config.BlockedPaths)
}
