package security

import (
	"fmt"
	"regexp"
	"strings"
)

// SecurityError represents security-related errors
type SecurityError struct {
	Code    string
	Message string
	Cause   error
}

func (e SecurityError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("security error [%s]: %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("security error [%s]: %s", e.Code, e.Message)
}

func (e SecurityError) Unwrap() error {
	return e.Cause
}

// Error codes
const (
	ErrCodeCommandNotAllowed    = "COMMAND_NOT_ALLOWED"
	ErrCodeDangerousPattern     = "DANGEROUS_PATTERN"
	ErrCodePathTraversal        = "PATH_TRAVERSAL"
	ErrCodePathBlocked         = "PATH_BLOCKED"
	ErrCodePathNotAllowed      = "PATH_NOT_ALLOWED"
	ErrCodeInvalidInput        = "INVALID_INPUT"
	ErrCodeCommandTooLong      = "COMMAND_TOO_LONG"
)

// CommandValidator interface for command validation
type CommandValidator interface {
	ValidateCommand(command string) error
	IsCommandAllowed(command string) bool
}

// PathValidator interface for path validation
type PathValidator interface {
	ValidatePath(path string) error
	IsPathAllowed(path string) bool
}

// InputSanitizer interface for input sanitization
type InputSanitizer interface {
	Sanitize(input string) string
}

// CommandValidatorImpl implements command validation
type CommandValidatorImpl struct {
	config *SecurityConfig
}

// NewCommandValidator creates a new command validator
func NewCommandValidator(config *SecurityConfig) CommandValidator {
	return &CommandValidatorImpl{config: config}
}

// ValidateCommand validates a command for security
func (v *CommandValidatorImpl) ValidateCommand(command string) error {
	if command == "" {
		return SecurityError{
			Code:    ErrCodeInvalidInput,
			Message: "empty command",
		}
	}

	// Check for dangerous patterns
	dangerousPatterns := []string{
		"rm -rf", "dd if=", "mkfs", "fdisk", "parted",
		"chmod 777", "chown root", "sudo", "su -",
		"wget", "curl", "nc ", "netcat", "telnet",
		"> /dev/null", "2>&1", "|", ";", "&",
		"$(", "`", "${", "$(", "${",
	}

	commandLower := strings.ToLower(command)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(commandLower, pattern) {
			return SecurityError{
				Code:    ErrCodeDangerousPattern,
				Message: fmt.Sprintf("command contains dangerous pattern: %s", pattern),
			}
		}
	}

	// Check for path traversal attempts
	if strings.Contains(command, "..") {
		return SecurityError{
			Code:    ErrCodePathTraversal,
			Message: "command contains path traversal attempts",
		}
	}

	// Validate against blocked paths
	for _, blockedPath := range v.config.BlockedPaths {
		if strings.Contains(command, blockedPath) {
			return SecurityError{
				Code:    ErrCodePathBlocked,
				Message: fmt.Sprintf("command references blocked path: %s", blockedPath),
			}
		}
	}

	return nil
}

// IsCommandAllowed checks if a command is allowed
func (v *CommandValidatorImpl) IsCommandAllowed(command string) bool {
	if command == "" {
		return false
	}

	// Split command into parts
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return false
	}

	// Check if base command is in allowlist
	allowedCommands := make(map[string]bool)
	for _, cmd := range v.config.AllowedCommands {
		allowedCommands[cmd] = true
	}

	return allowedCommands[parts[0]]
}

// PathValidatorImpl implements path validation
type PathValidatorImpl struct {
	config *SecurityConfig
}

// NewPathValidator creates a new path validator
func NewPathValidator(config *SecurityConfig) PathValidator {
	return &PathValidatorImpl{config: config}
}

// ValidatePath validates a file path for security
func (v *PathValidatorImpl) ValidatePath(path string) error {
	if path == "" {
		return SecurityError{
			Code:    ErrCodeInvalidInput,
			Message: "empty path",
		}
	}

	// Basic path traversal check
	if strings.Contains(path, "..") {
		return SecurityError{
			Code:    ErrCodePathTraversal,
			Message: "path contains traversal attempts",
		}
	}

	// Check against blocked paths
	for _, blockedPath := range v.config.BlockedPaths {
		if strings.HasPrefix(path, blockedPath) || strings.Contains(path, blockedPath) {
			return SecurityError{
				Code:    ErrCodePathBlocked,
				Message: fmt.Sprintf("path is blocked: %s", path),
			}
		}
	}

	// Check if path is in allowed paths (if specified)
	if len(v.config.AllowedPaths) > 0 {
		allowed := false
		for _, allowedPath := range v.config.AllowedPaths {
			if strings.HasPrefix(path, allowedPath) {
				allowed = true
				break
			}
		}
		if !allowed {
			return SecurityError{
				Code:    ErrCodePathNotAllowed,
				Message: fmt.Sprintf("path not in allowed paths: %s", path),
			}
		}
	}

	return nil
}

// IsPathAllowed checks if a path is allowed
func (v *PathValidatorImpl) IsPathAllowed(path string) bool {
	return v.ValidatePath(path) == nil
}

// InputSanitizerImpl implements input sanitization
type InputSanitizerImpl struct{}

// NewInputSanitizer creates a new input sanitizer
func NewInputSanitizer() InputSanitizer {
	return &InputSanitizerImpl{}
}

// Sanitize removes potentially dangerous characters from input
func (s *InputSanitizerImpl) Sanitize(input string) string {
	if input == "" {
		return input
	}

	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")

	// Remove control characters except newlines and tabs
	// Keep newlines (\n) and tabs (\t) for legitimate use cases
	re := regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F\x7F]`)
	input = re.ReplaceAllString(input, "")

	// Limit length to prevent DoS
	maxLength := 10000
	if len(input) > maxLength {
		input = input[:maxLength]
	}

	return strings.TrimSpace(input)
}
