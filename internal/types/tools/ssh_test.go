package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSSHCommandArgs_Validate(t *testing.T) {
	tests := []struct {
		name     string
		args     SSHCommandArgs
		expected bool
	}{
		{
			name: "valid minimal args",
			args: SSHCommandArgs{
				Host:    "example.com",
				Command: "ls -la",
			},
			expected: true,
		},
		{
			name: "valid with all args",
			args: SSHCommandArgs{
				Host:     "example.com",
				Command:  "ls -la",
				User:     "admin",
				Port:     "22",
				KeyPath:  "~/.ssh/id_rsa",
				Timeout:  30,
			},
			expected: true,
		},
		{
			name: "missing host",
			args: SSHCommandArgs{
				Host:    "",
				Command: "ls -la",
			},
			expected: false,
		},
		{
			name: "missing command",
			args: SSHCommandArgs{
				Host:    "example.com",
				Command: "",
			},
			expected: false,
		},
		{
			name: "invalid port",
			args: SSHCommandArgs{
				Host:    "example.com",
				Command: "ls -la",
				Port:    "99999",
			},
			expected: false,
		},
		{
			name: "invalid timeout",
			args: SSHCommandArgs{
				Host:    "example.com",
				Command: "ls -la",
				Timeout: 500,
			},
			expected: false,
		},
		{
			name: "valid with empty optional fields",
			args: SSHCommandArgs{
				Host:     "example.com",
				Command:  "ls -la",
				User:     "",
				Port:     "",
				KeyPath:  "",
				Timeout:  0,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.Validate()
			if tt.expected {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestNewSSHCommandArgs(t *testing.T) {
	args := NewSSHCommandArgs("example.com", "ls -la", "admin", "22", "~/.ssh/id_rsa")

	assert.Equal(t, "example.com", args.Host)
	assert.Equal(t, "ls -la", args.Command)
	assert.Equal(t, "admin", args.User)
	assert.Equal(t, "22", args.Port)
	assert.Equal(t, "~/.ssh/id_rsa", args.KeyPath)
}

func TestSSHCommandArgs_Validate_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		args     SSHCommandArgs
		expected bool
	}{
		{
			name: "valid port edge cases",
			args: SSHCommandArgs{
				Host:    "example.com",
				Command: "ls -la",
				Port:    "1",
			},
			expected: true,
		},
		{
			name: "valid port maximum",
			args: SSHCommandArgs{
				Host:    "example.com",
				Command: "ls -la",
				Port:    "65535",
			},
			expected: true,
		},
		{
			name: "valid timeout edge cases",
			args: SSHCommandArgs{
				Host:    "example.com",
				Command: "ls -la",
				Timeout: 1,
			},
			expected: true,
		},
		{
			name: "valid timeout maximum",
			args: SSHCommandArgs{
				Host:    "example.com",
				Command: "ls -la",
				Timeout: 300,
			},
			expected: true,
		},
		{
			name: "invalid port just over maximum",
			args: SSHCommandArgs{
				Host:    "example.com",
				Command: "ls -la",
				Port:    "65536",
			},
			expected: false,
		},
		{
			name: "invalid timeout just over maximum",
			args: SSHCommandArgs{
				Host:    "example.com",
				Command: "ls -la",
				Timeout: 301,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.Validate()
			if tt.expected {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
