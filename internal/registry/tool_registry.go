package registry

import (
	"context"
	"encoding/json"
	"fmt"

	"mini-mcp/internal/shared/logging"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// ToolHandler represents a function that handles a tool call
type ToolHandler[T any] func(ctx context.Context, req *mcp.CallToolRequest, args T) (*mcp.CallToolResult, any, error)

// ToolDefinition represents the definition of a tool
type ToolDefinition[T any] struct {
	Name        string
	Description string
	Handler     ToolHandler[T]
	Validator   func(T) error
}

// ToolRegistry manages tool registration with common patterns
type ToolRegistry struct {
	server *mcp.Server
	logger logging.Logger
}

// NewToolRegistry creates a new tool registry
func NewToolRegistry(server *mcp.Server, logger logging.Logger) *ToolRegistry {
	return &ToolRegistry{
		server: server,
		logger: logger,
	}
}

// RegisterTool registers a tool with common error handling and validation
func RegisterTool[T any](tr *ToolRegistry, def ToolDefinition[T]) {
	// Create a wrapper handler that provides common error handling
	wrapper := func(ctx context.Context, req *mcp.CallToolRequest, args T) (*mcp.CallToolResult, any, error) {
		// Validate arguments if validator is provided
		if def.Validator != nil {
			if err := def.Validator(args); err != nil {
				tr.logger.Error("Tool validation failed", err, map[string]any{
					"tool": def.Name,
				})
				return &mcp.CallToolResult{
					IsError: true,
					Content: []mcp.Content{&mcp.TextContent{Text: err.Error()}},
				}, nil, nil
			}
		}

		// Execute the actual handler
		result, data, err := def.Handler(ctx, req, args)
		if err != nil {
			tr.logger.Error("Tool execution failed", err, map[string]any{
				"tool": def.Name,
			})
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{&mcp.TextContent{Text: err.Error()}},
			}, nil, nil
		}

		// Log successful execution
		tr.logger.Info("Tool executed successfully", map[string]any{
			"tool": def.Name,
		})

		return result, data, nil
	}

	// Register the tool with the MCP server
	mcp.AddTool(tr.server, &mcp.Tool{
		Name:        def.Name,
		Description: def.Description,
	}, wrapper)
}

// RegisterSimpleTool registers a tool with a simple handler function
func RegisterSimpleTool[T any](tr *ToolRegistry, name, description string, handler ToolHandler[T], validator func(T) error) {
	RegisterTool(tr, ToolDefinition[T]{
		Name:        name,
		Description: description,
		Handler:     handler,
		Validator:   validator,
	})
}

// CreateErrorResult creates a standardized error result
func (tr *ToolRegistry) CreateErrorResult(message string, details map[string]any) (*mcp.CallToolResult, any, error) {
	tr.logger.Error("Tool error", fmt.Errorf("%s", message), map[string]any{
		"message": message,
		"details": details,
	})

	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{&mcp.TextContent{Text: message}},
	}, nil, nil
}

// CreateSuccessResult creates a standardized success result
func (tr *ToolRegistry) CreateSuccessResult(data any) (*mcp.CallToolResult, any, error) {
	// Convert data to JSON for consistent formatting
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return tr.CreateErrorResult("Failed to format result", map[string]any{
			"error": err.Error(),
		})
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(jsonData)}},
	}, data, nil
}

// CreateTextResult creates a simple text result
func (tr *ToolRegistry) CreateTextResult(text string) (*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}, text, nil
}
