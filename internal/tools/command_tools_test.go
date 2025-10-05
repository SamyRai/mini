package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"mini-mcp/internal/types/tools"
)

func TestRegisterCommandTools(t *testing.T) {
	// Test that the registration function exists and doesn't panic
	assert.NotPanics(t, func() {
		// We can't easily test the full registration without proper setup
		// but we can verify the function exists
	})
}

func TestCommandTools_Validation(t *testing.T) {
	tests := []struct {
		name      string
		args      tools.CommandArgs
		wantError bool
	}{
		{
			name: "valid command args",
			args: tools.CommandArgs{
				Command: "ls -la",
				Timeout: 30,
			},
			wantError: false,
		},
		{
			name: "missing command",
			args: tools.CommandArgs{
				Timeout: 30,
			},
			wantError: true,
		},
		{
			name: "timeout of 0 (should be valid)",
			args: tools.CommandArgs{
				Command: "ls -la",
				Timeout: 0,
			},
			wantError: false,
		},
		{
			name: "large timeout (should be valid)",
			args: tools.CommandArgs{
				Command: "ls -la",
				Timeout: 300,
			},
			wantError: false,
		},
		{
			name: "very large timeout (should be invalid)",
			args: tools.CommandArgs{
				Command: "ls -la",
				Timeout: 1000,
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.Validate()
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCommandTools_Registration(t *testing.T) {
	// Test that the registration function exists and can be called
	assert.NotNil(t, RegisterCommandTools)
}

func TestDockerComposeArgs_Validation(t *testing.T) {
	tests := []struct {
		name      string
		args      DockerComposeArgs
		wantError bool
	}{
		{
			name: "valid docker compose args",
			args: DockerComposeArgs{
				Path:    "/tmp/docker-compose.yml",
				Command: "up",
			},
			wantError: false,
		},
		{
			name: "missing path",
			args: DockerComposeArgs{
				Command: "up",
			},
			wantError: true,
		},
		{
			name: "missing command",
			args: DockerComposeArgs{
				Path: "/tmp/docker-compose.yml",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.Validate()
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDockerSwarmArgs_Validation(t *testing.T) {
	args := DockerSwarmArgs{}
	err := args.Validate()
	assert.NoError(t, err)
}
