// Package types defines the data structures used in the MCP protocol.
// This file serves as a central import point for all MCP types.
package types

import (
	// Import subpackages to make them available to users of the types package
	_ "mini-mcp/internal/types/resources"
	_ "mini-mcp/internal/types/tools"
)

// This file intentionally left mostly empty.
// Types have been organized into separate files and packages:
//
// - protocol.go: Core MCP protocol types (Tool, Resource, Metadata)
// - requests.go: Request types (ToolUseRequest, ResourceAccessRequest)
// - responses.go: Response types (ToolUseResponse, ResourceAccessResponse)
// - schema.go: JSON Schema related types
// - tools/: Tool-specific argument types
// - resources/: Resource-specific response types
