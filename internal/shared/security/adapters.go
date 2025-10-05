package security

import (
	"context"
	"os"
	"path/filepath"
)

// CommandSecurityAdapter adapts SecureCommandExecutor to domain command security interface
type CommandSecurityAdapter struct {
	executor *SecureCommandExecutor
}

// NewCommandSecurityAdapter creates a new command security adapter
func NewCommandSecurityAdapter(executor *SecureCommandExecutor) *CommandSecurityAdapter {
	return &CommandSecurityAdapter{
		executor: executor,
	}
}

// ExecuteCommand executes a command
func (a *CommandSecurityAdapter) ExecuteCommand(ctx context.Context, command string) (string, error) {
	return a.executor.ExecuteCommand(ctx, command)
}

// SanitizeInput sanitizes command input
func (a *CommandSecurityAdapter) SanitizeInput(input string) string {
	return a.executor.SanitizeInput(input)
}

// IsCommandAllowed checks if a command is allowed
func (a *CommandSecurityAdapter) IsCommandAllowed(command string) bool {
	// For now, allow all commands that pass basic validation
	// In a real implementation, this would check against an allowlist
	return command != "" && len(command) < 1000
}

// FileSecurityAdapter adapts SecureCommandExecutor to domain file security interface
type FileSecurityAdapter struct {
	executor *SecureCommandExecutor
}

// NewFileSecurityAdapter creates a new file security adapter
func NewFileSecurityAdapter(executor *SecureCommandExecutor) *FileSecurityAdapter {
	return &FileSecurityAdapter{
		executor: executor,
	}
}

// IsPathAllowed checks if a path is allowed
func (a *FileSecurityAdapter) IsPathAllowed(path string) bool {
	// For now, allow all paths that don't contain dangerous patterns
	// In a real implementation, this would check against a path allowlist
	return path != "" && !containsDangerousPatterns(path)
}

// SanitizePath sanitizes a file path
func (a *FileSecurityAdapter) SanitizePath(path string) string {
	// Basic path sanitization
	cleanPath := filepath.Clean(path)
	if filepath.IsAbs(cleanPath) {
		// For absolute paths, ensure they're not accessing system directories
		if isSystemPath(cleanPath) {
			return ""
		}
	}
	return cleanPath
}

// CreateDirectory creates a directory
func (a *FileSecurityAdapter) CreateDirectory(path string) error {
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, 0755)
}

// ReadFile reads a file and returns its content
func (a *FileSecurityAdapter) ReadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// WriteFile writes content to a file
func (a *FileSecurityAdapter) WriteFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// ListDirectory lists directory contents
func (a *FileSecurityAdapter) ListDirectory(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		files = append(files, entry.Name())
	}
	return files, nil
}

// DeleteFile deletes a file or directory
func (a *FileSecurityAdapter) DeleteFile(path string) error {
	return os.RemoveAll(path)
}

// containsDangerousPatterns checks if a path contains dangerous patterns
func containsDangerousPatterns(path string) bool {
	dangerousPatterns := []string{
		"..", "~", "/etc", "/var", "/usr", "/bin", "/sbin", "/sys", "/proc",
	}

	for _, pattern := range dangerousPatterns {
		if contains(path, pattern) {
			return true
		}
	}
	return false
}

// isSystemPath checks if a path is a system path
func isSystemPath(path string) bool {
	systemPaths := []string{
		"/etc", "/var", "/usr", "/bin", "/sbin", "/sys", "/proc", "/dev",
	}

	for _, sysPath := range systemPaths {
		if startsWith(path, sysPath) {
			return true
		}
	}
	return false
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

// startsWith checks if a string starts with a prefix
func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

// containsSubstring checks if a string contains a substring (simplified)
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
