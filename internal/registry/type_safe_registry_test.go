package registry

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"mini-mcp/internal/shared/logging"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestNewTypeSafeToolRegistry(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "1.0.0"}, nil)
	logger := logging.NewLogger(os.Stderr, logging.LogLevel("DEBUG"))

	registry := NewTypeSafeToolRegistry(server, logger)
	assert.NotNil(t, registry)
	assert.Equal(t, server, registry.server)
	assert.Equal(t, logger, registry.logger)
}

func TestToolBuilder(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "1.0.0"}, nil)
	registry := NewTypeSafeToolRegistry(server, nil)

	builder := NewToolBuilder[string](registry, "test_tool", "Test tool description")
	assert.NotNil(t, builder)
	assert.Equal(t, "test_tool", builder.definition.Name)
	assert.Equal(t, "Test tool description", builder.definition.Description)
}

func TestToolBuilder_WithHandler(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "1.0.0"}, nil)
	registry := NewTypeSafeToolRegistry(server, nil)

	handler := func(ctx context.Context, req *mcp.CallToolRequest, args string) (*mcp.CallToolResult, any, error) {
		return &mcp.CallToolResult{}, "result", nil
	}

	builder := NewToolBuilder[string](registry, "test_tool", "Test tool description").
		WithHandler(handler)

	assert.NotNil(t, builder.definition.Handler)
}

func TestToolBuilder_WithValidator(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "1.0.0"}, nil)
	registry := NewTypeSafeToolRegistry(server, nil)

	validator := func(args string) error {
		return nil
	}

	builder := NewToolBuilder[string](registry, "test_tool", "Test tool description").
		WithValidator(validator)

	assert.NotNil(t, builder.definition.Validator)
}

func TestValidationError_Error(t *testing.T) {
	err := ValidationError{
		Code:    "test_code",
		Message: "test message",
	}

	assert.Equal(t, "test message", err.Error())
}

func TestValidationError_Fields(t *testing.T) {
	err := ValidationError{
		Code:    "test_code",
		Message: "test message",
	}
	assert.Equal(t, "test_code", err.Code)
	assert.Equal(t, "test message", err.Message)
}

func TestTypeSafeToolRegistry_CreateErrorResult(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "1.0.0"}, nil)
	logger := logging.NewLogger(os.Stderr, logging.LogLevel("INFO"))
	registry := NewTypeSafeToolRegistry(server, logger)

	result, data, err := registry.CreateErrorResult("test error", map[string]any{"key": "value"})

	assert.NoError(t, err) // CreateErrorResult doesn't return an error, it returns a result
	assert.Nil(t, data)
	assert.NotNil(t, result)
	assert.True(t, result.IsError)
}

func TestTypeSafeToolRegistry_CreateTextResult(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "1.0.0"}, nil)
	logger := logging.NewLogger(os.Stderr, logging.LogLevel("INFO"))
	registry := NewTypeSafeToolRegistry(server, logger)

	result, data, err := registry.CreateTextResult("test text")

	assert.NoError(t, err)
	assert.Equal(t, "test text", data)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
}

func TestTypeSafeToolRegistry_CreateSuccessResult(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "1.0.0"}, nil)
	logger := logging.NewLogger(os.Stderr, logging.LogLevel("INFO"))
	registry := NewTypeSafeToolRegistry(server, logger)

	testData := map[string]string{"key": "value"}
	result, data, err := registry.CreateSuccessResult(testData)

	assert.NoError(t, err)
	assert.Equal(t, testData, data)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
}

func TestToolBuilder_Register(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "1.0.0"}, nil)
	logger := logging.NewLogger(os.Stderr, logging.LogLevel("INFO"))
	registry := NewTypeSafeToolRegistry(server, logger)

	builder := NewToolBuilder[string](registry, "test_tool", "Test tool")

	// Register should not panic and should succeed since we don't validate handlers
	err := builder.Register()
	assert.NoError(t, err)
}
