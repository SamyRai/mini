package server

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"mini-mcp/internal/shared/logging"
	"mini-mcp/internal/shared/security"
)

// CommandExecutor provides common command execution patterns
type CommandExecutor struct {
	security *security.SecureCommandExecutor
	logger   logging.Logger
}

// NewCommandExecutor creates a new command executor
func NewCommandExecutor(security *security.SecureCommandExecutor, logger logging.Logger) *CommandExecutor {
	return &CommandExecutor{
		security: security,
		logger:   logger,
	}
}

// ExecuteCommand executes a command with common security and logging
func (ce *CommandExecutor) ExecuteCommand(ctx context.Context, command string, timeout int) (string, error) {
	// Validate command security
	if !ce.security.IsCommandAllowed(command) {
		return "", fmt.Errorf("command not allowed: %s", command)
	}

	// Set default timeout if not provided
	if timeout <= 0 {
		timeout = 30
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	// Log command execution
	ce.logger.Info("Executing command", map[string]any{
		"command": command,
		"timeout": timeout,
	})

	// Execute command
	start := time.Now()
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	output, err := cmd.Output()
	duration := time.Since(start)

	// Log execution result
	if err != nil {
		ce.logger.Error("Command execution failed", err, map[string]any{
			"command":  command,
			"error":    err.Error(),
			"duration": duration.String(),
		})
		return "", fmt.Errorf("command failed: %w", err)
	}

	ce.logger.Info("Command executed successfully", map[string]any{
		"command":       command,
		"duration":      duration.String(),
		"output_length": len(output),
	})

	return string(output), nil
}

// ExecuteSSHCommand executes a command over SSH with common patterns
func (ce *CommandExecutor) ExecuteSSHCommand(ctx context.Context, host, command, user, port, keyPath string, timeout int) (string, error) {
	// Validate SSH command security
	if !ce.security.IsCommandAllowed(command) {
		return "", fmt.Errorf("SSH command not allowed: %s", command)
	}

	// Validate SSH key path if provided
	if keyPath != "" {
		if err := ce.security.ValidatePath(keyPath); err != nil {
			return "", fmt.Errorf("SSH key path validation failed: %w", err)
		}
	}

	// Build SSH command
	sshCmd := []string{"ssh"}

	if user != "" {
		sshCmd = append(sshCmd, "-l", user)
	}
	if port != "" {
		sshCmd = append(sshCmd, "-p", port)
	}
	if keyPath != "" {
		sshCmd = append(sshCmd, "-i", keyPath)
	}

	sshCmd = append(sshCmd, host, command)

	// Set default timeout if not provided
	if timeout <= 0 {
		timeout = 30
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	// Log SSH command execution
	ce.logger.Info("Executing SSH command", map[string]any{
		"host":    host,
		"command": command,
		"user":    user,
		"timeout": timeout,
	})

	// Execute SSH command
	start := time.Now()
	cmd := exec.CommandContext(ctx, sshCmd[0], sshCmd[1:]...)
	output, err := cmd.Output()
	duration := time.Since(start)

	// Log execution result
	if err != nil {
		ce.logger.Error("SSH command execution failed", err, map[string]any{
			"host":     host,
			"command":  command,
			"error":    err.Error(),
			"duration": duration.String(),
		})
		return "", fmt.Errorf("SSH command failed: %w", err)
	}

	ce.logger.Info("SSH command executed successfully", map[string]any{
		"host":          host,
		"command":       command,
		"duration":      duration.String(),
		"output_length": len(output),
	})

	return string(output), nil
}

// ExecuteDockerCompose executes Docker Compose commands with common patterns
func (ce *CommandExecutor) ExecuteDockerCompose(ctx context.Context, path, command string, detached, removeVolumes bool) (string, error) {
	// Validate Docker Compose path
	if err := ce.security.ValidatePath(path); err != nil {
		return "", fmt.Errorf("docker compose path validation failed: %w", err)
	}

	// Validate Docker Compose command is allowed
	allowedCommands := []string{"up", "down", "ps", "logs", "restart", "stop", "start"}
	commandAllowed := false
	for _, cmd := range allowedCommands {
		if command == cmd {
			commandAllowed = true
			break
		}
	}
	if !commandAllowed {
		return "", fmt.Errorf("docker compose command not allowed: %s", command)
	}

	// Build docker-compose command
	args := []string{"docker-compose", "-f", path}

	switch command {
	case "up":
		args = append(args, "up")
		if detached {
			args = append(args, "-d")
		}
	case "down":
		args = append(args, "down")
		if removeVolumes {
			args = append(args, "-v")
		}
	case "ps", "logs", "restart", "stop", "start":
		args = append(args, command)
	default:
		return "", fmt.Errorf("unsupported command: %s", command)
	}

	// Log Docker Compose command execution
	ce.logger.Info("Executing Docker Compose command", map[string]any{
		"path":           path,
		"command":        command,
		"detached":       detached,
		"remove_volumes": removeVolumes,
	})

	// Execute command
	start := time.Now()
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	output, err := cmd.Output()
	duration := time.Since(start)

	// Log execution result
	if err != nil {
		ce.logger.Error("Docker Compose command execution failed", err, map[string]any{
			"path":     path,
			"command":  command,
			"error":    err.Error(),
			"duration": duration.String(),
		})
		return "", fmt.Errorf("docker-compose command failed: %w", err)
	}

	ce.logger.Info("Docker Compose command executed successfully", map[string]any{
		"path":          path,
		"command":       command,
		"duration":      duration.String(),
		"output_length": len(output),
	})

	return string(output), nil
}

// ExecuteSystemCommand executes system commands with common patterns
func (ce *CommandExecutor) ExecuteSystemCommand(ctx context.Context, command string, args ...string) (string, error) {
	// Validate command security
	if !ce.security.IsCommandAllowed(command) {
		return "", fmt.Errorf("system command not allowed: %s", command)
	}

	// Log system command execution
	ce.logger.Info("Executing system command", map[string]any{
		"command": command,
		"args":    args,
	})

	// Execute command
	start := time.Now()
	cmd := exec.CommandContext(ctx, command, args...)
	output, err := cmd.Output()
	duration := time.Since(start)

	// Log execution result
	if err != nil {
		ce.logger.Error("System command execution failed", err, map[string]any{
			"command":  command,
			"args":     args,
			"error":    err.Error(),
			"duration": duration.String(),
		})
		return "", fmt.Errorf("system command failed: %w", err)
	}

	ce.logger.Info("System command executed successfully", map[string]any{
		"command":       command,
		"args":          args,
		"duration":      duration.String(),
		"output_length": len(output),
	})

	return string(output), nil
}

// ParsePortNumbers parses port numbers from netstat/ss output
func (ce *CommandExecutor) ParsePortNumbers(output string) []int {
	ports := make([]int, 0)
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		// Skip header lines
		if strings.Contains(line, "LISTEN") || strings.Contains(line, "Local Address") {
			continue
		}

		// Extract port from lines like "*:8080 *:* LISTEN"
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			address := parts[1] // Should be something like "*:8080"
			if strings.Contains(address, ":") {
				portStr := strings.Split(address, ":")[1]
				if port, err := strconv.Atoi(portStr); err == nil && port > 0 {
					ports = append(ports, port)
				}
			}
		}
	}

	return ports
}

// KillProcessGracefully kills a process by PID with graceful shutdown
func (ce *CommandExecutor) KillProcessGracefully(ctx context.Context, pid int) bool {
	// First try SIGTERM (graceful)
	cmd := exec.CommandContext(ctx, "kill", "-TERM", strconv.Itoa(pid))
	err := cmd.Run()

	if err != nil {
		// If SIGTERM fails, try SIGKILL (forceful)
		cmd = exec.CommandContext(ctx, "kill", "-KILL", strconv.Itoa(pid))
		err = cmd.Run()
	}

	success := err == nil
	if success {
		ce.logger.Info("Process killed successfully", map[string]any{
			"pid": pid,
		})
	} else {
		ce.logger.Error("Failed to kill process", err, map[string]any{
			"pid":   pid,
			"error": err.Error(),
		})
	}

	return success
}
