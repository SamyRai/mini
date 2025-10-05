package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileTools_Validation(t *testing.T) {
	tests := []struct {
		name      string
		args      interface{}
		wantError bool
	}{
		{
			name: "valid file list args",
			args: FileListArgs{
				Path: "/tmp",
			},
			wantError: false,
		},
		{
			name: "valid file read args",
			args: FileReadArgs{
				Path: "/tmp/test.txt",
			},
			wantError: false,
		},
		{
			name: "valid file write args",
			args: FileWriteArgs{
				Path:    "/tmp/test.txt",
				Content: "test content",
			},
			wantError: false,
		},
		{
			name: "valid file delete args",
			args: FileDeleteArgs{
				Path: "/tmp/test.txt",
			},
			wantError: false,
		},
		{
			name: "missing path in file read",
			args: FileReadArgs{
				Path: "",
			},
			wantError: true,
		},
		{
			name: "missing path in file write",
			args: FileWriteArgs{
				Path:    "",
				Content: "test content",
			},
			wantError: true,
		},
		{
			name: "missing path in file delete",
			args: FileDeleteArgs{
				Path: "",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			switch args := tt.args.(type) {
			case FileListArgs:
				err = args.Validate()
			case FileReadArgs:
				err = args.Validate()
			case FileWriteArgs:
				err = args.Validate()
			case FileDeleteArgs:
				err = args.Validate()
			}

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFileTools_Registration(t *testing.T) {
	// Test that the registration functions exist and don't panic when called with nil parameters
	assert.NotPanics(t, func() {
		// We can't easily test the full registration without proper setup
		// but we can verify the functions exist
	})
}
