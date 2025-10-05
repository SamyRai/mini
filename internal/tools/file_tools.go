package tools

import (
	"context"

	"mini-mcp/internal/handlers/core"
	"mini-mcp/internal/registry"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// ===== TYPE-SAFE ARGUMENT STRUCTURES =====

// FileListArgs represents arguments for the ls command
type FileListArgs struct {
	Path string `json:"path" jsonschema:"Directory path to list"`
}

// FileReadArgs represents arguments for the cat command
type FileReadArgs struct {
	Path string `json:"path" jsonschema:"File path to read"`
}

// FileWriteArgs represents arguments for the write command
type FileWriteArgs struct {
	Path    string `json:"path" jsonschema:"File path to write"`
	Content string `json:"content" jsonschema:"Content to write"`
}

// FileDeleteArgs represents arguments for the rm command
type FileDeleteArgs struct {
	Path string `json:"path" jsonschema:"File or directory path to remove"`
}

// ===== VALIDATION METHODS (STRATEGY PATTERN) =====

// Validate validates FileListArgs
func (args FileListArgs) Validate() error {
	// Basic validation - path will be validated by security layer
	return nil
}

// Validate validates FileReadArgs
func (args FileReadArgs) Validate() error {
	if args.Path == "" {
		return registry.NewValidationError("missing_path", "path is required")
	}
	return nil
}

// Validate validates FileWriteArgs
func (args FileWriteArgs) Validate() error {
	if args.Path == "" {
		return registry.NewValidationError("missing_path", "path is required")
	}
	return nil
}

// Validate validates FileDeleteArgs
func (args FileDeleteArgs) Validate() error {
	if args.Path == "" {
		return registry.NewValidationError("missing_path", "path is required")
	}
	return nil
}

// ===== TOOL REGISTRATION USING DESIGN PATTERNS =====

// RegisterFileTools registers file-related tools using proper design patterns
func RegisterFileTools(server *mcp.Server, toolRegistry *registry.TypeSafeToolRegistry, fileHandler *core.FileHandlerImpl) {
	// ls - List directory contents (Builder Pattern)
	lsBuilder := registry.NewToolBuilder[FileListArgs](toolRegistry, "ls", "List directory contents with security validation")

	lsBuilder.
		WithHandler(func(ctx context.Context, req *mcp.CallToolRequest, args FileListArgs) (*mcp.CallToolResult, any, error) {
			if args.Path == "" {
				args.Path = "."
			}

			output, err := fileHandler.ListDirectory(ctx, map[string]any{
				"path": args.Path,
			})
			if err != nil {
				errorResult, _, _ := toolRegistry.CreateErrorResult(err.Error(), map[string]any{
					"path": args.Path,
				})
				return errorResult, nil, nil
			}

			successResult, _, _ := toolRegistry.CreateTextResult(output)
			return successResult, nil, nil
		}).
		WithValidator(func(args FileListArgs) error {
			return args.Validate()
		})

	if err := lsBuilder.Register(); err != nil {
		// Log error but continue - tool registration failure should not crash the server
		return
	}

	// cat - Read file contents (Builder Pattern)
	catBuilder := registry.NewToolBuilder[FileReadArgs](toolRegistry, "cat", "Read file contents with security validation")

	catBuilder.
		WithHandler(func(ctx context.Context, req *mcp.CallToolRequest, args FileReadArgs) (*mcp.CallToolResult, any, error) {
			output, err := fileHandler.ReadFile(ctx, map[string]any{
				"path": args.Path,
			})
			if err != nil {
				errorResult, _, _ := toolRegistry.CreateErrorResult(err.Error(), map[string]any{
					"path": args.Path,
				})
				return errorResult, nil, nil
			}

			successResult, _, _ := toolRegistry.CreateTextResult(output)
			return successResult, nil, nil
		}).
		WithValidator(func(args FileReadArgs) error {
			return args.Validate()
		})

	if err := catBuilder.Register(); err != nil {
		// Log error but continue - tool registration failure should not crash the server
		return
	}

	// write - Write content to file (Builder Pattern)
	writeBuilder := registry.NewToolBuilder[FileWriteArgs](toolRegistry, "write", "Write content to file with security validation")

	writeBuilder.
		WithHandler(func(ctx context.Context, req *mcp.CallToolRequest, args FileWriteArgs) (*mcp.CallToolResult, any, error) {
			_, err := fileHandler.WriteFile(ctx, map[string]any{
				"path":    args.Path,
				"content": args.Content,
			})
			if err != nil {
				errorResult, _, _ := toolRegistry.CreateErrorResult(err.Error(), map[string]any{
					"path": args.Path,
				})
				return errorResult, nil, nil
			}

			successResult, _, _ := toolRegistry.CreateTextResult("File written successfully")
			return successResult, nil, nil
		}).
		WithValidator(func(args FileWriteArgs) error {
			return args.Validate()
		})

	if err := writeBuilder.Register(); err != nil {
		// Log error but continue - tool registration failure should not crash the server
		return
	}

	// rm - Remove file or directory (Builder Pattern)
	rmBuilder := registry.NewToolBuilder[FileDeleteArgs](toolRegistry, "rm", "Remove file or directory with security validation")

	rmBuilder.
		WithHandler(func(ctx context.Context, req *mcp.CallToolRequest, args FileDeleteArgs) (*mcp.CallToolResult, any, error) {
			_, err := fileHandler.DeleteFile(ctx, map[string]any{
				"path": args.Path,
			})
			if err != nil {
				errorResult, _, _ := toolRegistry.CreateErrorResult(err.Error(), map[string]any{
					"path": args.Path,
				})
				return errorResult, nil, nil
			}

			successResult, _, _ := toolRegistry.CreateTextResult("File/directory removed successfully")
			return successResult, nil, nil
		}).
		WithValidator(func(args FileDeleteArgs) error {
			return args.Validate()
		})

	if err := rmBuilder.Register(); err != nil {
		// Log error but continue - tool registration failure should not crash the server
		return
	}
}
