package server

import (
	"mini-mcp/internal/health"
	"mini-mcp/internal/handlers/core"
	"mini-mcp/internal/registry"
	"mini-mcp/internal/tools"
	"mini-mcp/internal/shared/logging"
	"mini-mcp/internal/shared/security"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Deps encapsulates external dependencies required to build the MCP server.
type Deps struct {
    Logger        logging.Logger
    Security      *security.SecureCommandExecutor
    HealthChecker *health.HealthChecker
}

// BuildServer constructs and returns a configured MCP server instance.
// version is the semantic version of the implementation, used for handshake metadata.
func BuildServer(deps Deps, version string) *mcp.Server {
    server := mcp.NewServer(&mcp.Implementation{Name: "mini-mcp", Version: version}, nil)

    // Initialize health checker if not provided
    if deps.HealthChecker == nil {
        deps.HealthChecker = health.CreateDefaultHealthChecker(version)
    }

    // Create tool registry and command executor
    toolRegistry := registry.NewTypeSafeToolRegistry(server, deps.Logger)
    executor := registry.NewCommandExecutor(deps.Security, deps.Logger)

    // Create handlers
    commandHandler := core.NewCommandHandler(nil, deps.Logger) // Will be properly injected
    fileHandler := core.NewFileHandler(nil, deps.Logger)       // Will be properly injected
    systemHandler := core.NewSystemHandler(nil, deps.Logger)   // Will be properly injected

    // Register tools by category (following SRP)
    tools.RegisterCommandTools(server, toolRegistry, commandHandler.(*core.CommandHandlerImpl))
    tools.RegisterFileTools(server, toolRegistry, fileHandler.(*core.FileHandlerImpl))
    tools.RegisterSystemTools(server, toolRegistry, systemHandler.(*core.SystemHandlerImpl), deps.HealthChecker, deps.Logger)
    tools.RegisterInfrastructureTools(server, toolRegistry, executor)
    tools.RegisterPortProcessTools(server, toolRegistry, executor)

    // Register resources
    registerResources(server)

    return server
}
