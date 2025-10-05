package tools

import (
	"context"

	"mini-mcp/internal/registry"
	"mini-mcp/internal/types/tools"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterInfrastructureTools registers infrastructure-related tools
func RegisterInfrastructureTools(server *mcp.Server, toolRegistry *registry.TypeSafeToolRegistry, executor *registry.CommandExecutor) {
	// ssh - Execute remote commands over SSH
	sshBuilder := registry.NewToolBuilder[tools.SSHCommandArgs](toolRegistry, "ssh", "Execute a remote command over SSH with security validation")

	sshBuilder.
		WithHandler(func(ctx context.Context, req *mcp.CallToolRequest, args tools.SSHCommandArgs) (*mcp.CallToolResult, any, error) {
			output, err := executor.ExecuteSSHCommand(ctx, args.Host, args.Command, args.User, args.Port, args.KeyPath, args.Timeout)
			if err != nil {
				errorResult, _, _ := toolRegistry.CreateErrorResult(err.Error(), map[string]any{
					"host":    args.Host,
					"command": args.Command,
					"user":    args.User,
				})
				return errorResult, nil, nil
			}

			successResult, _, _ := toolRegistry.CreateTextResult(output)
			return successResult, nil, nil
		}).
		WithValidator(func(args tools.SSHCommandArgs) error {
			return args.Validate()
		}).
		Register()

	// docker_compose - Docker Compose operations
	dockerComposeBuilder := registry.NewToolBuilder[DockerComposeArgs](toolRegistry, "docker_compose", "Execute Docker Compose operations with security validation")

	dockerComposeBuilder.
		WithHandler(func(ctx context.Context, req *mcp.CallToolRequest, args DockerComposeArgs) (*mcp.CallToolResult, any, error) {
			output, err := executor.ExecuteDockerCompose(ctx, args.Path, args.Command, args.Detached, args.RemoveVolumes)
			if err != nil {
				errorResult, _, _ := toolRegistry.CreateErrorResult(err.Error(), map[string]any{
					"path":    args.Path,
					"command": args.Command,
				})
				return errorResult, nil, nil
			}

			successResult, _, _ := toolRegistry.CreateTextResult(output)
			return successResult, nil, nil
		}).
		WithValidator(func(args DockerComposeArgs) error {
			if args.Path == "" {
				return registry.NewValidationError("missing_path", "path is required")
			}
			if args.Command == "" {
				return registry.NewValidationError("missing_command", "command is required")
			}
			return nil
		}).
		Register()

	// docker_swarm - Docker Swarm operations
	dockerSwarmBuilder := registry.NewToolBuilder[DockerSwarmArgs](toolRegistry, "docker_swarm", "Get Docker Swarm cluster information")

	dockerSwarmBuilder.
		WithHandler(func(ctx context.Context, req *mcp.CallToolRequest, args DockerSwarmArgs) (*mcp.CallToolResult, any, error) {
			output, err := executor.ExecuteSystemCommand(ctx, "docker", "swarm", "info")
			if err != nil {
				errorResult, _, _ := toolRegistry.CreateErrorResult(err.Error(), map[string]any{})
				return errorResult, nil, nil
			}

			successResult, _, _ := toolRegistry.CreateTextResult(output)
			return successResult, nil, nil
		}).
		WithValidator(func(args DockerSwarmArgs) error {
			return nil
		}).
		Register()
}

// DockerComposeArgs represents arguments for Docker Compose operations
type DockerComposeArgs struct {
	Path          string `json:"path" jsonschema:"Path to docker-compose.yml file"`
	Command       string `json:"command" jsonschema:"Docker Compose command (up, down, ps, logs, etc.)"`
	Detached      bool   `json:"detached" jsonschema:"Run in detached mode"`
	RemoveVolumes bool   `json:"remove_volumes" jsonschema:"Remove volumes when stopping"`
}

// Validate validates DockerComposeArgs
func (args DockerComposeArgs) Validate() error {
	if args.Path == "" {
		return registry.NewValidationError("missing_path", "path is required")
	}
	if args.Command == "" {
		return registry.NewValidationError("missing_command", "command is required")
	}
	return nil
}

// DockerSwarmArgs represents arguments for Docker Swarm operations
type DockerSwarmArgs struct{}

// Validate validates DockerSwarmArgs
func (args DockerSwarmArgs) Validate() error {
	return nil
}
