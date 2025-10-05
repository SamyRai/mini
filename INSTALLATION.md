# Mini MCP Installation Guide

This guide explains how to install the Mini MCP binary system-wide for use with Cursor and other MCP clients.

## Prerequisites

- Go 1.21 or later
- sudo access (for system-wide installation)
- Docker (optional, for Docker-related tools)
- Docker Compose v2 (optional, for compose operations)

## Building the MCP Binary

First, build the MCP binary:

```bash
# Build the MCP binary
go build -o mini-mcp cmd/mini-mcp/main.go
```

## Installation Options

### Option 1: Automatic Installation (Recommended)

Use the provided CLI tool:

```bash
make install-all
```

This will:
- Build the MCP binary
- Copy it to `/usr/local/bin/`
- Set proper permissions
- Configure VS Code and Cursor settings
- Verify the installation

### Option 2: Manual Installation

If you prefer to install manually:

```bash
# Build the MCP binary
go build -o mini-mcp cmd/mini-mcp/main.go

# Copy binary to system-wide location
sudo cp mini-mcp /usr/local/bin/mini-mcp

# Set executable permissions
sudo chmod +x /usr/local/bin/mini-mcp
```

## Verification

Verify that the binary is accessible system-wide:

```bash
which mini-mcp
```

You should see:
```
/usr/local/bin/mini-mcp
```

## Cursor Configuration

Update your Cursor MCP configuration (`~/.cursor/mcp.json`) to use the unified binary:

```json
{
  "mcpServers": {
    "mini-mcp": {
      "command": "mini-mcp",
      "args": ["-mode=server"],
      "env": {
        "ENVIRONMENT": "development",
        "LOG_LEVEL": "INFO"
      },
      "cwd": ".",
      "description": "Unified Mini MCP Infrastructure Management Tool"
    }
  }
}
```

## Testing the Installation

Test that the MCP binary works correctly:

```bash
# Test MCP server mode
echo '{"type": "initialize"}' | mini-mcp -mode=server

# Test CLI mode
echo '{"type": "run", "payload": {"command": "ls -la"}}' | mini-mcp -mode=cli
```

## Docker Configuration (Optional)

If you plan to use Docker-related tools, configure Docker contexts and hosts:

### Docker Contexts

```bash
# List available contexts
docker context ls

# Create a new context for remote Docker
docker context create remote --docker "host=tcp://remote:2376"

# Use a specific context
docker context use remote
```

### Docker Host Configuration

```bash
# Set Docker host environment variable
export DOCKER_HOST=tcp://remote:2376

# Or use SSH-based connection
export DOCKER_HOST=ssh://user@remote
```

### Testing Docker Tools

Test Docker functionality:

```bash
# Test docker swarm info
echo '{"type": "tools/call", "params": {"name": "docker_swarm_info", "arguments": {}}}' | mini-mcp

# Test docker compose with context
echo '{"type": "tools/call", "params": {"name": "docker_compose_up", "arguments": {"path": "/project", "context": "production"}}}' | mini-mcp
```

## Available Modes

The MCP binary supports two modes:

### Server Mode (`-mode=server`)
- MCP protocol server for Cursor and other MCP clients
- Communicates via JSON messages over stdin/stdout
- Used for programmatic access to tools

### CLI Mode (`-mode=cli`)
- Interactive command-line interface
- Direct human interaction
- JSON-based command format

## Uninstallation

To remove the system-wide binary:

```bash
sudo rm /usr/local/bin/mini-mcp
```

## Troubleshooting

### Binary Not Found
If `which mini-mcp` returns nothing:
1. Check if the binary was installed: `ls -la /usr/local/bin/mini-mcp`
2. Reinstall using the CLI: `make install-all`

### Permission Denied
If you get permission errors:
1. Check file permissions: `ls -la /usr/local/bin/mini-mcp`
2. Fix permissions: `sudo chmod +x /usr/local/bin/mini-mcp`

### Build Errors
If the build fails:
1. Ensure Go 1.21+ is installed: `go version`
2. Check dependencies: `go mod tidy`
3. Try building manually: `go build -o mini-mcp cmd/mini-mcp/main.go`
