package registry

import (
	"context"
	"encoding/json"
	"fmt"

	"mini-mcp/internal/shared/logging"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// ===== TYPE DEFINITIONS =====

// TypeSafeToolHandler represents a type-safe function that handles a tool call
type TypeSafeToolHandler[T any] func(ctx context.Context, req *mcp.CallToolRequest, args T) (*mcp.CallToolResult, any, error)

// TypeSafeToolDefinition represents a type-safe tool definition
type TypeSafeToolDefinition[T any] struct {
	Name        string
	Description string
	Handler     TypeSafeToolHandler[T]
	Validator   func(T) error
}

// ===== FACTORY PATTERN: Tool Registry Factory =====

// TypeSafeToolRegistry provides type-safe tool registration with proper design patterns
type TypeSafeToolRegistry struct {
	server     *mcp.Server
	logger     logging.Logger
	validators map[string]ValidationStrategy
}

// NewTypeSafeToolRegistry creates a new type-safe tool registry (Factory Pattern)
func NewTypeSafeToolRegistry(server *mcp.Server, logger logging.Logger) *TypeSafeToolRegistry {
	return &TypeSafeToolRegistry{
		server:     server,
		logger:     logger,
		validators: make(map[string]ValidationStrategy),
	}
}

// ===== BUILDER PATTERN: Fluent Tool Configuration =====

// ToolBuilder provides a fluent interface for building tools (Builder Pattern)
type ToolBuilder[T any] struct {
	definition TypeSafeToolDefinition[T]
	registry   *TypeSafeToolRegistry
}

// NewToolBuilder creates a new tool builder (Factory Pattern)
func NewToolBuilder[T any](tsr *TypeSafeToolRegistry, name, description string) *ToolBuilder[T] {
	return &ToolBuilder[T]{
		definition: TypeSafeToolDefinition[T]{
			Name:        name,
			Description: description,
		},
		registry: tsr,
	}
}

// WithHandler sets the tool handler
func (tb *ToolBuilder[T]) WithHandler(handler TypeSafeToolHandler[T]) *ToolBuilder[T] {
	tb.definition.Handler = handler
	return tb
}

// WithValidator sets the tool validator
func (tb *ToolBuilder[T]) WithValidator(validator func(T) error) *ToolBuilder[T] {
	tb.definition.Validator = validator
	return tb
}

// Register registers the tool (Template Method Pattern)
func (tb *ToolBuilder[T]) Register() error {
	return RegisterTypeSafeTool(tb.registry, tb.definition)
}

// ===== STRATEGY PATTERN: Validation Strategies =====

// ValidationStrategy defines the interface for validation strategies
type ValidationStrategy interface {
	Validate(args any) error
}

// TypeValidationStrategy implements validation for a specific type
type TypeValidationStrategy[T any] struct {
	validator func(T) error
}

// Validate validates arguments using the strategy
func (tvs *TypeValidationStrategy[T]) Validate(args any) error {
	typedArgs, ok := args.(T)
	if !ok {
		return NewValidationError("type_mismatch", "invalid argument type")
	}
	return tvs.validator(typedArgs)
}

// ===== TEMPLATE METHOD PATTERN: Tool Registration Workflow =====

// RegisterTypeSafeTool registers a type-safe tool with common error handling and validation (Template Method)
func RegisterTypeSafeTool[T any](tsr *TypeSafeToolRegistry, def TypeSafeToolDefinition[T]) error {
	// Step 1: Create wrapper handler with cross-cutting concerns (Decorator Pattern)
	wrapper := createTypeSafeWrapper(tsr, def)

	// Step 2: Register the tool with the MCP server
	mcp.AddTool(tsr.server, &mcp.Tool{
		Name:        def.Name,
		Description: def.Description,
	}, wrapper)

	// Step 3: Register validation strategy
	if def.Validator != nil {
		tsr.validators[def.Name] = &TypeValidationStrategy[T]{validator: def.Validator}
	}

	return nil
}

// createTypeSafeWrapper creates a wrapper with cross-cutting concerns (Decorator Pattern)
func createTypeSafeWrapper[T any](tsr *TypeSafeToolRegistry, def TypeSafeToolDefinition[T]) func(ctx context.Context, req *mcp.CallToolRequest, args any) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, args any) (*mcp.CallToolResult, any, error) {
		// Type-safe argument handling
		typedArgs, ok := args.(T)
		if !ok {
			return createErrorResult(tsr, "type_mismatch", "invalid argument type")
		}

		// Validate arguments using strategy pattern
		if def.Validator != nil {
			if err := def.Validator(typedArgs); err != nil {
				tsr.logger.Error("Tool validation failed", err, map[string]any{
					"tool": def.Name,
				})
				return createErrorResult(tsr, "validation_failed", err.Error())
			}
		}

		// Execute the actual handler
		result, data, err := def.Handler(ctx, req, typedArgs)
		if err != nil {
			tsr.logger.Error("Tool execution failed", err, map[string]any{
				"tool": def.Name,
			})
			return createErrorResult(tsr, "execution_failed", err.Error())
		}

		// Log successful execution
		tsr.logger.Info("Tool executed successfully", map[string]any{
			"tool": def.Name,
		})

		return result, data, nil
	}
}

// RegisterTypeSafeSimpleTool registers a type-safe tool with a simple handler function
func RegisterTypeSafeSimpleTool[T any](tsr *TypeSafeToolRegistry, name, description string, handler TypeSafeToolHandler[T], validator func(T) error) error {
	return RegisterTypeSafeTool(tsr, TypeSafeToolDefinition[T]{
		Name:        name,
		Description: description,
		Handler:     handler,
		Validator:   validator,
	})
}

// ===== UTILITY FUNCTIONS =====

// createErrorResult creates a standardized error result
func createErrorResult(tsr *TypeSafeToolRegistry, code, message string) (*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{&mcp.TextContent{Text: message}},
	}, nil, nil
}

// CreateErrorResult creates a standardized error result (public method)
func (tsr *TypeSafeToolRegistry) CreateErrorResult(message string, details map[string]any) (*mcp.CallToolResult, any, error) {
	tsr.logger.Error("Tool error", fmt.Errorf("%s", message), map[string]any{
		"message": message,
		"details": details,
	})

	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{&mcp.TextContent{Text: message}},
	}, nil, nil
}

// CreateSuccessResult creates a standardized success result
func (tsr *TypeSafeToolRegistry) CreateSuccessResult(data any) (*mcp.CallToolResult, any, error) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return createErrorResult(tsr, "formatting_failed", "Failed to format result")
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(jsonData)}},
	}, data, nil
}

// CreateTextResult creates a simple text result
func (tsr *TypeSafeToolRegistry) CreateTextResult(text string) (*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}, text, nil
}

// ValidationError represents a validation error
type ValidationError struct {
	Code    string
	Message string
}

func (ve ValidationError) Error() string {
	return ve.Message
}

// NewValidationError creates a new validation error
func NewValidationError(code, message string) ValidationError {
	return ValidationError{
		Code:    code,
		Message: message,
	}
}
