package tools

import (
	"mini-mcp/internal/shared/validation"
)

// SSHCommandArgs represents arguments for ssh_command.
// Example:
//
//	{"host": "example.com", "command": "ls -la", "user": "admin", "port": "22", "key_path": "~/.ssh/id_rsa", "timeout": 30}
type SSHCommandArgs struct {
	// Host is the hostname or IP address of the SSH server
	Host string `json:"host"`
	// Command is the command to execute on the remote server
	Command string `json:"command"`
	// User is the SSH username (defaults to "root" if not specified)
	User string `json:"user,omitempty"`
	// Port is the SSH port (defaults to "22" if not specified)
	Port string `json:"port,omitempty"`
	// KeyPath is the path to the SSH private key file
	KeyPath string `json:"key_path,omitempty"`
	// Timeout is the connection timeout in seconds (defaults to 10 if not specified)
	Timeout int `json:"timeout,omitempty"`
}

// Validate checks if the SSH command arguments are valid.
func (args *SSHCommandArgs) Validate() error {
	// Validate host
	if err := validation.StringRequired("host", args.Host); err != nil {
		return err
	}

	// Validate command
	if err := validation.StringRequired("command", args.Command); err != nil {
		return err
	}

	// Validate port if provided
	if args.Port != "" {
		if err := validation.Port("port", args.Port); err != nil {
			return err
		}
	}

	// Validate timeout if provided
	if args.Timeout > 0 {
		if err := validation.Timeout("timeout", args.Timeout); err != nil {
			return err
		}
	}

	return nil
}

// NewSSHCommandArgs creates a new SSHCommandArgs with the given host, command, and optional parameters.
func NewSSHCommandArgs(host, command, user, port, keyPath string) *SSHCommandArgs {
	return &SSHCommandArgs{
		Host:    host,
		Command: command,
		User:    user,
		Port:    port,
		KeyPath: keyPath,
	}
}
