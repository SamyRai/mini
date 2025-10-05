package tools

import (
	"context"

	"mini-mcp/internal/handlers/core"
	"mini-mcp/internal/registry"
	"mini-mcp/internal/types/tools"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterCommandTools registers command-related tools using proper design patterns
func RegisterCommandTools(server *mcp.Server, toolRegistry *registry.TypeSafeToolRegistry, commandHandler *core.CommandHandlerImpl) {
	// Use Builder Pattern for fluent tool configuration
	builder := registry.NewToolBuilder[tools.CommandArgs](toolRegistry, "run", "Execute a shell command securely with allowlisting and timeout controls")

	// Configure the tool using Builder Pattern
	err := builder.
		WithHandler(func(ctx context.Context, req *mcp.CallToolRequest, args tools.CommandArgs) (*mcp.CallToolResult, any, error) {
			output, err := commandHandler.ExecuteCommand(ctx, map[string]any{
				"command": args.Command,
				"timeout": float64(args.Timeout),
			})
			if err != nil {
				errorResult, _, _ := toolRegistry.CreateErrorResult(err.Error(), map[string]any{
					"command": args.Command,
					"timeout": args.Timeout,
				})
				return errorResult, nil, nil
			}

			successResult, _, _ := toolRegistry.CreateTextResult(output)
			return successResult, nil, nil
		}).
		WithValidator(func(args tools.CommandArgs) error {
			return args.Validate()
		}).
		Register()

	if err != nil {
		// Handle registration error
		panic(err)
	}
}