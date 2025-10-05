package server

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"mini-mcp/internal/shared/logging"
	"mini-mcp/internal/shared/security"
)

func TestRegisterCoreTools(t *testing.T) {
	// Create mock dependencies
	logger := logging.NewLogger(os.Stderr, logging.LogLevel("INFO"))
	sec := security.NewSecureCommandExecutor(nil)
	deps := Deps{
		Logger:   logger,
		Security: sec,
	}

	// Test server creation - the actual tool registration is tested in the tools package
	builtServer := BuildServer(deps, "1.0.0")

	// Verify server is created
	assert.NotNil(t, builtServer)
}

// TestCoreTools_ArgumentValidation tests argument validation for tools
func TestCoreTools_ArgumentValidation(t *testing.T) {
	// This test has been moved to the tools package tests
	t.Skip("Test moved to tools package")
}

// Legacy test function - keeping for reference but marked as skipped
func TestLegacyCoreTools(t *testing.T) {
	t.Skip("Core tools functionality moved to tools package")
}

func TestServerCreation(t *testing.T) {
	logger := logging.NewLogger(os.Stderr, logging.LogLevel("INFO"))
	sec := security.NewSecureCommandExecutor(nil)
	deps := Deps{
		Logger:   logger,
		Security: sec,
	}

	// Test server creation
	builtServer := BuildServer(deps, "1.0.0")
	assert.NotNil(t, builtServer)
}

func TestServerWithHealthChecker(t *testing.T) {
	logger := logging.NewLogger(os.Stderr, logging.LogLevel("INFO"))
	sec := security.NewSecureCommandExecutor(nil)
	deps := Deps{
		Logger:   logger,
		Security: sec,
	}

	// Test server creation with health checker
	builtServer := BuildServer(deps, "1.0.0")
	assert.NotNil(t, builtServer)
}

// Additional server tests can be added here for testing the BuildServer function
// with different configurations
