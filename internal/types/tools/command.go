// Package tools contains types for MCP tool arguments and responses.
package tools

import (
	"mini-mcp/internal/shared/validation"
)

// CommandArgs represents arguments for the execute_command tool.
// This tool provides secure command execution with built-in safety measures.
//
// Security Features:
// - Command allowlisting (only pre-approved commands allowed)
// - Input sanitization to prevent injection attacks
// - Timeout controls to prevent hanging processes
// - Output size limits to prevent memory exhaustion
// - Working directory restrictions
//
// Allowed Commands:
// - System utilities: ls, cat, head, tail, grep, find, wc, sort, uniq
// - Process management: ps, top
// - System monitoring: df, du, free, uptime, who, w
// - Version control: git
// - Container management: docker
// - Infrastructure: nomad, consul, terraform
//
// Example:
//
//	{"command": "ls -la /tmp", "timeout": 30}
type CommandArgs struct {
	// Command is the shell command to execute
	// Must be in the allowed command list for security
	Command string `json:"command"`

	// Timeout is the maximum execution time in seconds (optional)
	// Default: 30 seconds, Maximum: 300 seconds (5 minutes)
	Timeout int `json:"timeout,omitempty"`
}

// Validate checks if the command arguments are valid.
func (args *CommandArgs) Validate() error {
	// Validate command
	if err := validation.StringRequired("command", args.Command); err != nil {
		return err
	}

	// Validate timeout if provided
	if args.Timeout > 0 {
		if err := validation.Timeout("timeout", args.Timeout); err != nil {
			return err
		}
	}

	return nil
}

// NewCommandArgs creates a new CommandArgs with the given command.
func NewCommandArgs(command string, timeout ...int) *CommandArgs {
	args := &CommandArgs{
		Command: command,
	}

	if len(timeout) > 0 {
		args.Timeout = timeout[0]
	}

	return args
}
