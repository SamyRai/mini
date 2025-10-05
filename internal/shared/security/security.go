package security

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"mini-mcp/internal/shared/logging"
)

// CommandSecurity provides secure command execution with allowlisting and sandboxing
type CommandSecurity struct {
	AllowedCommands  map[string]bool
	WorkingDirectory string
	Timeout          time.Duration
	UserPermissions  string
	MaxOutputSize    int64
}

// SecurityConfig holds configuration for security features
type SecurityConfig struct {
	// Command allowlist - only these commands are allowed
	AllowedCommands []string `json:"allowed_commands"`

	// Working directory restrictions
	WorkingDirectory string `json:"working_directory"`

	// Timeout for command execution
	CommandTimeout time.Duration `json:"command_timeout"`

	// Maximum output size in bytes
	MaxOutputSize int64 `json:"max_output_size"`

	// User to run commands as (empty for current user)
	RunAsUser string `json:"run_as_user"`

	// Environment variables to allow
	AllowedEnvVars []string `json:"allowed_env_vars"`

	// Path restrictions
	AllowedPaths []string `json:"allowed_paths"`
	BlockedPaths []string `json:"blocked_paths"`
}

// SecureCommandExecutor handles secure command execution
type SecureCommandExecutor struct {
	config          *SecurityConfig
	allowedCommands map[string]bool
	activeCommands  map[string]*exec.Cmd
	mutex           sync.RWMutex
	validator       CommandValidator
	pathValidator   PathValidator
	sanitizer       InputSanitizer
}

// NewSecureCommandExecutor creates a new secure command executor
func NewSecureCommandExecutor(config *SecurityConfig) *SecureCommandExecutor {
	if config == nil {
		config = DefaultSecurityConfig()
	}
	allowedCommands := make(map[string]bool)
	for _, cmd := range config.AllowedCommands {
		allowedCommands[cmd] = true
	}

	validator := NewCommandValidator(config)
	pathValidator := NewPathValidator(config)
	sanitizer := NewInputSanitizer()

	return &SecureCommandExecutor{
		config:          config,
		allowedCommands: allowedCommands,
		activeCommands:  make(map[string]*exec.Cmd),
		mutex:           sync.RWMutex{},
		validator:       validator,
		pathValidator:   pathValidator,
		sanitizer:       sanitizer,
	}
}

// Cleanup terminates any running commands and cleans up resources
func (e *SecureCommandExecutor) Cleanup() {
	for id, cmd := range e.activeCommands {
		if cmd.Process != nil {
			if err := cmd.Process.Kill(); err != nil {
				// Log error but don't fail the operation
				logging.GetGlobalLogger().Error("Failed to kill process", err, map[string]any{
					"command_id": id,
				})
			}
		}
		delete(e.activeCommands, id)
	}
}

// DefaultSecurityConfig returns a secure default configuration
func DefaultSecurityConfig() *SecurityConfig {
	return &SecurityConfig{
		AllowedCommands: []string{
			"ls", "cat", "head", "tail", "grep", "find", "wc", "sort", "uniq",
			"ps", "top", "df", "du", "free", "uptime", "who", "w",
			"git", "docker", "consul", "terraform",
		},
		WorkingDirectory: "/tmp",
		CommandTimeout:   30 * time.Second,
		MaxOutputSize:    1024 * 1024, // 1MB
		AllowedEnvVars:   []string{"PATH", "HOME", "USER", "PWD"},
		AllowedPaths:     []string{"/tmp", "/var/log", "/proc"},
		BlockedPaths:     []string{"/etc/passwd", "/etc/shadow", "/root", "/home"},
	}
}

// ExecuteCommand safely executes a command with security checks
func (s *SecureCommandExecutor) ExecuteCommand(ctx context.Context, command string) (string, error) {
	// Sanitize input first
	command = s.sanitizer.Sanitize(command)

	// Validate and sanitize command
	if err := s.validator.ValidateCommand(command); err != nil {
		return "", fmt.Errorf("command validation failed: %w", err)
	}

	// Split command into parts
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "", fmt.Errorf("empty command")
	}

	// Check if command is allowed
	if !s.allowedCommands[parts[0]] {
		return "", fmt.Errorf("command '%s' is not allowed", parts[0])
	}

	// Create command with context and timeout
	ctx, cancel := context.WithTimeout(ctx, s.config.CommandTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)

	// Set working directory
	if s.config.WorkingDirectory != "" {
		cmd.Dir = s.config.WorkingDirectory
	}

	// Set allowed environment variables
	cmd.Env = s.filterEnvironment(os.Environ())

	// Execute command
	output, err := cmd.CombinedOutput()

	// Check output size
	if int64(len(output)) > s.config.MaxOutputSize {
		return "", fmt.Errorf("command output exceeds maximum size limit")
	}

	return string(output), err
}

// filterEnvironment filters environment variables to only include allowed ones
func (s *SecureCommandExecutor) filterEnvironment(env []string) []string {
	allowed := make(map[string]bool)
	for _, v := range s.config.AllowedEnvVars {
		allowed[v] = true
	}

	var filtered []string
	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 && allowed[parts[0]] {
			filtered = append(filtered, e)
		}
	}

	return filtered
}

// ValidatePath checks if a path is allowed using the path validator
func (s *SecureCommandExecutor) ValidatePath(path string) error {
	return s.pathValidator.ValidatePath(path)
}

// SanitizeInput sanitizes input using the input sanitizer
func (s *SecureCommandExecutor) SanitizeInput(input string) string {
	return s.sanitizer.Sanitize(input)
}

// IsCommandAllowed checks if a command is allowed using the command validator
func (s *SecureCommandExecutor) IsCommandAllowed(command string) bool {
	return s.validator.IsCommandAllowed(command)
}

// GetPathValidator returns the path validator for external use
func (s *SecureCommandExecutor) GetPathValidator() PathValidator {
	return s.pathValidator
}

// GetCommandValidator returns the command validator for external use
func (s *SecureCommandExecutor) GetCommandValidator() CommandValidator {
	return s.validator
}

// GetInputSanitizer returns the input sanitizer for external use
func (s *SecureCommandExecutor) GetInputSanitizer() InputSanitizer {
	return s.sanitizer
}
