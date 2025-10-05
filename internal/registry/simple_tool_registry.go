package registry

import (
	"context"
	"fmt"

	"mini-mcp/internal/shared/logging"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// SimpleToolHandler represents a function that handles a tool call
type SimpleToolHandler func(ctx context.Context, req *mcp.CallToolRequest, args any) (*mcp.CallToolResult, any, error)

// SimpleToolDefinition represents the definition of a tool
type SimpleToolDefinition struct {
	Name        string
	Description string
	Handler     SimpleToolHandler
	Validator   func(any) error
}

// SimpleToolRegistry provides a simplified tool registration without generics
type SimpleToolRegistry struct {
	server *mcp.Server
	logger logging.Logger
}

// NewSimpleToolRegistry creates a new simple tool registry
func NewSimpleToolRegistry(server *mcp.Server, logger logging.Logger) *SimpleToolRegistry {
	return &SimpleToolRegistry{
		server: server,
		logger: logger,
	}
}

// RegisterTool registers a tool with common error handling and validation
func (str *SimpleToolRegistry) RegisterTool(def SimpleToolDefinition) {
	// Create a wrapper handler that provides common error handling
	wrapper := func(ctx context.Context, req *mcp.CallToolRequest, args any) (*mcp.CallToolResult, any, error) {
		// Validate arguments if validator is provided
		if def.Validator != nil {
			if err := def.Validator(args); err != nil {
				str.logger.Error("Tool validation failed", err, map[string]any{
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
			str.logger.Error("Tool execution failed", err, map[string]any{
				"tool": def.Name,
			})
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{&mcp.TextContent{Text: err.Error()}},
			}, nil, nil
		}

		// Log successful execution
		str.logger.Info("Tool executed successfully", map[string]any{
			"tool": def.Name,
		})

		return result, data, nil
	}

	// Register the tool with the MCP server
	mcp.AddTool(str.server, &mcp.Tool{
		Name:        def.Name,
		Description: def.Description,
	}, wrapper)
}

// RegisterSimpleTool registers a tool with a simple handler function
func (str *SimpleToolRegistry) RegisterSimpleTool(name, description string, handler SimpleToolHandler, validator func(any) error) {
	str.RegisterTool(SimpleToolDefinition{
		Name:        name,
		Description: description,
		Handler:     handler,
		Validator:   validator,
	})
}

// CreateErrorResult creates a standardized error result
func (str *SimpleToolRegistry) CreateErrorResult(message string, details map[string]any) (*mcp.CallToolResult, any, error) {
	str.logger.Error("Tool error", fmt.Errorf("%s", message), map[string]any{
		"message": message,
		"details": details,
	})

	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{&mcp.TextContent{Text: message}},
	}, nil, nil
}

// CreateTextResult creates a simple text result
func (str *SimpleToolRegistry) CreateTextResult(text string) (*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}, text, nil
}
