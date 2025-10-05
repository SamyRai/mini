package command

import (
	"context"
	"errors"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var (
	ErrInvalidCommand = errors.New("invalid command")
	ErrInvalidTimeout = errors.New("invalid timeout")
)

// Command represents a command to be executed
type Command struct {
	Command string
	Timeout int
	Args    []string
}

// Result represents the result of command execution
type Result struct {
	Output    string
	Error     string
	ExitCode  int
	Duration  time.Duration
	Timestamp time.Time
}

// Repository defines the interface for command operations
type Repository interface {
	Execute(ctx context.Context, cmd *Command) (*Result, error)
	Validate(ctx context.Context, cmd *Command) error
	Sanitize(ctx context.Context, cmd *Command) error
}

// RepositoryImpl implements the Repository interface
type RepositoryImpl struct {
	// Command repository implementation with security validation
}

// NewRepository creates a new command repository
func NewRepository() Repository {
	return &RepositoryImpl{}
}

// Execute executes a command
func (r *RepositoryImpl) Execute(ctx context.Context, cmd *Command) (*Result, error) {
	// Create context with timeout
	timeout := time.Duration(cmd.Timeout) * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Execute command
	start := time.Now()
	execCmd := exec.CommandContext(ctx, "sh", "-c", cmd.Command)
	output, err := execCmd.Output()

	duration := time.Since(start)

	result := &Result{
		Output:    string(output),
		Duration:  duration,
		Timestamp: start,
	}

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitError.ExitCode()
			result.Error = string(exitError.Stderr)
		} else {
			result.Error = err.Error()
		}
		return result, err
	}

	result.ExitCode = 0
	return result, nil
}

// Validate validates a command
func (r *RepositoryImpl) Validate(ctx context.Context, cmd *Command) error {
	if cmd.Command == "" {
		return ErrInvalidCommand
	}
	if cmd.Timeout < 0 {
		return ErrInvalidTimeout
	}
	return nil
}

// Sanitize sanitizes command input
func (r *RepositoryImpl) Sanitize(ctx context.Context, cmd *Command) error {
	// Comprehensive command sanitization
	cmd.Command = sanitizeCommand(cmd.Command)
	cmd.Args = sanitizeArgs(cmd.Args)
	return nil
}

// sanitizeCommand performs comprehensive command sanitization
func sanitizeCommand(command string) string {
	if command == "" {
		return command
	}

	// Remove null bytes and control characters (except newlines and tabs for legitimate use)
	command = removeControlChars(command)

	// Normalize whitespace
	command = normalizeWhitespace(command)

	// Remove potentially dangerous patterns
	command = removeDangerousPatterns(command)

	// Limit command length
	maxCommandLength := 10000
	if len(command) > maxCommandLength {
		command = command[:maxCommandLength]
	}

	return strings.TrimSpace(command)
}

// sanitizeArgs sanitizes command arguments
func sanitizeArgs(args []string) []string {
	if len(args) == 0 {
		return args
	}

	sanitized := make([]string, len(args))
	for i, arg := range args {
		sanitized[i] = sanitizeArg(arg)
	}

	return sanitized
}

// sanitizeArg sanitizes a single argument
func sanitizeArg(arg string) string {
	if arg == "" {
		return arg
	}

	// Remove control characters
	arg = removeControlChars(arg)

	// Normalize whitespace
	arg = normalizeWhitespace(arg)

	// Remove dangerous patterns
	arg = removeDangerousPatterns(arg)

	// Limit argument length
	maxArgLength := 1000
	if len(arg) > maxArgLength {
		arg = arg[:maxArgLength]
	}

	return strings.TrimSpace(arg)
}

// removeControlChars removes control characters except newlines and tabs
func removeControlChars(input string) string {
	// Keep newlines and tabs for legitimate use cases like multi-line scripts
	var result strings.Builder
	for _, r := range input {
		if r == '\n' || r == '\t' || r == '\r' {
			result.WriteRune(r)
		} else if r >= 32 { // Printable ASCII and above
			result.WriteRune(r)
		}
		// Skip control characters (0-31 except newlines/tabs)
	}
	return result.String()
}

// normalizeWhitespace normalizes whitespace characters
func normalizeWhitespace(input string) string {
	// Replace multiple spaces with single space
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(strings.TrimSpace(input), " ")
}

// removeDangerousPatterns removes potentially dangerous patterns
func removeDangerousPatterns(input string) string {
	// Remove shell metacharacters that could be used for injection
	dangerousPatterns := []string{
		";", "&", "|", "`", "$", "(", ")", "{", "}", "[", "]",
		">", "<", "!", "*", "?", "~", "^", "\\",
	}

	result := input
	for _, pattern := range dangerousPatterns {
		result = strings.ReplaceAll(result, pattern, "")
	}

	return result
}
