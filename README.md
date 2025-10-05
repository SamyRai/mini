# Mini MCP - Infrastructure Management Platform

[![Go Version](https://img.shields.io/badge/go-1.25+-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/YOUR_USERNAME/mini-mcp)](https://goreportcard.com/report/github.com/YOUR_USERNAME/mini-mcp)
[![CI Status](https://github.com/YOUR_USERNAME/mini-mcp/workflows/CI/badge.svg)](https://github.com/YOUR_USERNAME/mini-mcp/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/YOUR_USERNAME/mini-mcp.svg)](https://pkg.go.dev/github.com/YOUR_USERNAME/mini-mcp)
[![Release](https://img.shields.io/github/release/YOUR_USERNAME/mini-mcp.svg)](https://github.com/YOUR_USERNAME/mini-mcp/releases)

**🏗️ [Architecture Documentation](docs/README_ARCHITECTURE.md) | [🛠️ Tools Documentation](docs/README_TOOLS.md) | [⚙️ Proxmox Configuration](docs/PROXMOX_CONFIG.md) | [🔒 Type Safety](docs/TYPE_SAFETY_IMPROVEMENTS.md) | [🤖 Agent Guide](docs/AGENT.md)**

A production-ready Model Context Protocol (MCP) server and CLI tool for infrastructure management with comprehensive security, authentication, monitoring, and health check capabilities.

## 📋 Table of Contents

- [✨ Features](#-features)
- [🚀 Quick Start](#-quick-start)
- [📖 Usage](#-usage)
- [🐳 Docker Management](#-docker-management)
- [⚙️ Configuration](#️-configuration)
- [🧪 Testing](#-testing)
- [🔧 Development](#-development)
- [📚 Documentation](#-documentation)
- [🤝 Contributing](#-contributing)
- [📄 License](#-license)

## 🚀 Features

### 🔒 Production-Ready Security & Safety
- **Multi-Layer Security**: Command, path, and input validation with dedicated security layers
- **Sandboxed Execution**: Restricted working directories and environment variables
- **Advanced Input Sanitization**: Comprehensive sanitization preventing injection attacks
- **Path Traversal Protection**: Sophisticated path validation against directory traversal
- **Timeout Enforcement**: Configurable command execution timeouts with graceful cancellation
- **Security Error Tracking**: Detailed error reporting with stack traces and suggestions

### 🔐 Enterprise Authentication & Authorization
- **API Key Management**: Secure API key generation and validation with rotation support
- **Advanced Rate Limiting**: Configurable rate limiting with sliding window and burst handling
- **IP Whitelisting**: Restrict access to specific IP addresses and CIDR ranges
- **Request Correlation**: Unique request IDs for tracking and debugging across services
- **Audit Logging**: Comprehensive audit trails for security events

### 📊 Production Monitoring & Observability
- **Structured Logging**: JSON-formatted logs with contextual metadata and correlation IDs
- **Advanced Metrics Collection**: Request counts, response times, error rates, and performance percentiles (P95, P99)
- **Comprehensive Health Checks**: System resources, filesystem, network, security, and dependency monitoring
- **Performance Monitoring**: Real-time performance tracking with alerting thresholds
- **Observability Tools**: Built-in metrics endpoint for monitoring dashboards and alerting

### 🛠️ Production-Ready Core Tools
- **Command Execution**: Run shell commands securely with comprehensive security validation and performance tracking
- **File Operations**: Read, write, list, and delete files with path validation and detailed error reporting
- **System Monitoring**: Get system information with resource usage and health status
- **Health Checks**: Comprehensive health monitoring with dependency tracking and alerting
- **Performance Metrics**: Real-time metrics collection with detailed observability data
- **Graceful Operations**: Proper resource cleanup and graceful shutdown capabilities

### 🚀 Infrastructure Management
- **Docker Management**: Enhanced compose operations with context, host, and path support
- **Docker Swarm**: Cluster information and node management with remote context support
- **Docker Contexts**: Support for multiple Docker contexts and remote hosts
- **Git Operations**: Repository cloning and management
- **SSH Operations**: Remote command execution
- **Documentation**: Fetch command and tool documentation

## 🏗️ Architecture

The project follows clean architecture principles with clear separation of concerns and production-ready patterns:

```
mini-mcp/
├── cmd/
│   └── mini-mcp/         # MCP server and CLI executable with graceful shutdown
├── internal/
│   ├── application/      # Application services with dependency injection
│   ├── domain/          # Domain logic with enhanced security & error handling
│   ├── handlers/        # Business logic handlers with metrics tracking
│   ├── health/          # Comprehensive health checks and monitoring
│   ├── server/          # MCP protocol server with observability
│   ├── shared/          # Shared utilities (auth, config, logging, security, errors)
│   │   ├── auth/        # Authentication and authorization
│   │   ├── config/      # Environment-specific configuration management
│   │   ├── errors/      # Structured error handling with stack traces
│   │   ├── logging/     # Structured logging with performance metrics
│   │   ├── security/    # Multi-layer security validation and sandboxing
│   │   └── validation/  # Input validation and sanitization
│   └── types/           # Type definitions and schemas
└── pkg/                 # Public packages (if any)
```

## 🚀 Quick Start

### Prerequisites
- **Go 1.25 or later** - Required for building and running
- **Docker** (optional) - For Docker-related tools and operations
- **Docker Compose v2** (optional) - For compose operations
- **sudo access** (for system-wide installation)

### Installation

#### Option 1: Quick Installation (Recommended)

```bash
# Clone the repository
git clone <repository-url>
cd mini-mcp

# Build and install system-wide with one command
make install-all
```

This automatically:
- Builds the MCP binary with optimizations
- Installs it to `/usr/local/bin/`
- Configures VS Code and Cursor settings
- Verifies the installation

#### Option 2: Manual Installation

1. **Clone and build**
   ```bash
   git clone <repository-url>
   cd mini-mcp

   # Build the MCP binary with optimizations
   go build -ldflags="-s -w" -o mini-mcp cmd/mini-mcp/main.go
   ```

2. **Install system-wide (optional)**
   ```bash
   # Copy to system PATH
   sudo cp mini-mcp /usr/local/bin/mini-mcp
   sudo chmod +x /usr/local/bin/mini-mcp
   ```

#### Option 3: Development Installation

For development, you can run directly without installation:
```bash
# Run as MCP server
go run cmd/mini-mcp/main.go -mode=server

# Run as CLI tool
go run cmd/mini-mcp/main.go -mode=cli
```

#### Verification

After installation, verify everything works correctly:

```bash
# Check binary is in PATH
which mini-mcp

# Test MCP server mode
echo '{"type": "initialize"}' | mini-mcp -mode=server

# Test CLI mode
echo '{"type": "run", "payload": {"command": "ls -la"}}' | mini-mcp -mode=cli
```

#### Docker Configuration (Optional)

If using Docker tools, configure Docker contexts:

```bash
# List available contexts
docker context ls

# Create context for remote Docker
docker context create remote --docker "host=tcp://remote:2376"

# Set Docker host environment variable
export DOCKER_HOST=tcp://remote:2376
```

#### Uninstallation

To remove the system-wide binary:
```bash
sudo rm /usr/local/bin/mini-mcp
```

#### Troubleshooting

**Binary Not Found**
If `which mini-mcp` returns nothing:
1. Check if the binary was installed: `ls -la /usr/local/bin/mini-mcp`
2. Reinstall using: `make install-all`

**Permission Denied**
If you get permission errors:
1. Check file permissions: `ls -la /usr/local/bin/mini-mcp`
2. Fix permissions: `sudo chmod +x /usr/local/bin/mini-mcp`

**Build Errors**
If the build fails:
1. Ensure Go 1.25+ is installed: `go version`
2. Check dependencies: `go mod tidy`
3. Try building manually: `go build -o mini-mcp cmd/mini-mcp/main.go`

## 📖 Usage

### MCP Server Mode (for Cursor and other MCP clients)

The MCP binary can run as an MCP protocol server:

```bash
# Run in server mode
./mini-mcp -mode=server

# Or use the system-wide installation
mini-mcp -mode=server
```

**Cursor Configuration** (`~/.cursor/mcp.json`):
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
      "description": "Mini MCP Infrastructure Management Tool"
    }
  }
}
```

### CLI Mode (Interactive)

The MCP binary can also run as an interactive CLI tool:

```bash
# Run in CLI mode
./mini-mcp -mode=cli

# Or use the system-wide installation
mini-mcp -mode=cli
```

**Available Production-Ready MCP Tools**:
```json
{"tool": "run", "arguments": {"command": "ls -la"}}
{"tool": "ls", "arguments": {"path": "/tmp"}}
{"tool": "cat", "arguments": {"path": "/etc/hosts"}}
{"tool": "write", "arguments": {"path": "test.txt", "content": "Hello World"}}
{"tool": "rm", "arguments": {"path": "test.txt"}}
{"tool": "system", "arguments": {"metric": "processes"}}
{"tool": "metrics", "arguments": {}}
{"tool": "ssh", "arguments": {"host": "server.example.com", "command": "uptime"}}
{"tool": "docker_compose", "arguments": {"path": "/project", "command": "up", "detached": true}}
{"tool": "docker_swarm", "arguments": {"context": "production"}}
{"tool": "port_process_tools", "arguments": {"command": "list_ports"}}
```

**New Metrics Tool**:
```json
{"tool": "metrics", "arguments": {}}
```
Returns comprehensive application metrics including:
- Request counts and response times
- Error rates and performance percentiles
- System resource usage and health status
- Security validation statistics

## 🐳 Docker Management

Mini MCP provides comprehensive Docker management capabilities with support for multiple contexts, remote hosts, and custom configurations.

### Docker Compose Operations

**Docker Compose Operations:**

*Start services:*
```json
{
  "tool": "docker_compose",
  "arguments": {
    "path": "/project",
    "command": "up",
    "detached": true,
    "context": "production",
    "host": "tcp://remote:2376",
    "compose_file": "/custom/compose.yml"
  }
}
```

*Stop services:*
```json
{
  "tool": "docker_compose",
  "arguments": {
    "path": "/project",
    "command": "down",
    "remove_volumes": true,
    "context": "production",
    "host": "tcp://remote:2376"
  }
}
```

**Docker Swarm Operations:**

*Get swarm cluster information:*
```json
{
  "tool": "docker_swarm",
  "arguments": {
    "context": "production",
    "host": "tcp://remote:2376",
    "docker_path": "/usr/local/bin/docker"
  }
}
```

**SSH Operations:**

*Execute remote command:*
```json
{
  "tool": "ssh",
  "arguments": {
    "host": "server.example.com",
    "command": "uptime",
    "user": "admin",
    "port": "22"
  }
}
```

**Port/Process Management:**

*List network ports:*
```json
{
  "tool": "port_process_tools",
  "arguments": {
    "command": "list_ports"
  }
}
```

*Find processes using port 8080:*
```json
{
  "tool": "port_process_tools",
  "arguments": {
    "command": "find_port",
    "port": 8080
  }
}
```

### Docker Configuration Options

- **`context`** - Docker context selection (`--context`)
- **`host`** - Docker host URL (`-H` or `DOCKER_HOST` environment variable)
- **`docker_path`** - Custom docker binary path
- **`compose_file`** - Custom docker-compose.yml file path (for compose commands)

### Remote Docker Operations

Mini MCP supports remote Docker operations through:
- **Docker contexts** for switching between different Docker environments
- **Remote hosts** for connecting to remote Docker daemons
- **Custom binary paths** for containerized environments
- **Swarm cluster management** across multiple nodes

## ⚙️ Production-Ready Configuration

The application supports comprehensive configuration management with environment-specific defaults, file-based configuration, and validation:

```bash
# Environment (development, staging, production)
export ENVIRONMENT=development

# Logging
export LOG_LEVEL=DEBUG

# Configuration file (optional)
export CONFIG_FILE=/path/to/config.json

# Security (enhanced with multi-layer validation)
export SECURITY_WORKING_DIR=/tmp
export SECURITY_COMMAND_TIMEOUT=30s
export SECURITY_MAX_OUTPUT_SIZE=1048576
export SECURITY_ALLOWED_COMMANDS=ls,cat,head,tail,grep,find,wc,sort,uniq,ps,top,df,du,free,uptime,who,w,git,docker,nomad,consul,terraform
export SECURITY_ALLOWED_PATHS=/tmp,/var/log,/proc
export SECURITY_BLOCKED_PATHS=/etc/passwd,/etc/shadow,/root,/home

# Authentication (enterprise-grade)
export AUTH_RATE_LIMITING=100ms
export AUTH_MAX_REQUESTS=1000
export AUTH_WINDOW_SIZE=1h
export AUTH_IP_WHITELIST=127.0.0.1,::1
export AUTH_API_KEYS=key1:value1,key2:value2

# Performance (monitoring and limits)
export PERF_MAX_CONCURRENT_REQUESTS=100
export PERF_REQUEST_TIMEOUT=30s
export PERF_CACHE_ENABLED=false
export PERF_CACHE_TTL=5m

# Health monitoring (production checks)
export HEALTH_CHECK_INTERVAL=30s
export HEALTH_METRICS_ENABLED=true
```

### Environment-Specific Configurations

**Development**: Relaxed security, detailed logging, extended timeouts
**Staging**: Balanced security and performance, moderate logging
**Production**: Maximum security, optimized performance, minimal logging

## 🧪 Production-Ready Testing

### Comprehensive Test Suite

The application includes a comprehensive test suite with unit tests, integration tests, and security validation:

```bash
# Run all tests with coverage
go test -cover ./...

# Run specific test packages
go test ./internal/shared/security/ -v
go test ./internal/domain/file/ -v
go test ./internal/shared/errors/ -v

# Run tests with race detection
go test -race ./...

# Generate test coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### MCP Protocol Testing

Test the MCP server protocol:

```bash
# Test initialization
echo '{"type": "initialize"}' | mini-mcp -mode=server

# Test tools listing
echo '{"type": "tools/list"}' | mini-mcp -mode=server

# Test command execution with security validation
echo '{"type": "tools/call", "data": {"name": "run", "arguments": {"command": "ls -la"}}}' | mini-mcp -mode=server

# Test new metrics tool
echo '{"type": "tools/call", "data": {"name": "metrics", "arguments": {}}}' | mini-mcp -mode=server
```

### Security Testing

Test security features and validation:

```bash
# Test dangerous command blocking
echo '{"type": "tools/call", "data": {"name": "run", "arguments": {"command": "rm -rf /"}}}' | mini-mcp -mode=server

# Test path traversal protection
echo '{"type": "tools/call", "data": {"name": "ls", "arguments": {"path": "../../../etc"}}}' | mini-mcp -mode=server
```

### Performance Testing

Test performance and monitoring:

```bash
# Load testing with multiple concurrent requests
for i in {1..10}; do
  echo '{"type": "tools/call", "data": {"name": "run", "arguments": {"command": "sleep 1"}}}' | mini-mcp -mode=server &
done
wait

# Monitor metrics during load
echo '{"type": "tools/call", "data": {"name": "metrics", "arguments": {}}}' | mini-mcp -mode=server
```

## 🔧 Production-Ready Development

### Building

```bash
# Build the production-ready MCP binary with optimizations
go build -ldflags="-s -w" -o mini-mcp cmd/mini-mcp/main.go

# Build with race detection (development)
go build -race -o mini-mcp cmd/mini-mcp/main.go

# Build and install system-wide
make install-all

# Cross-compile for different platforms
GOOS=linux GOARCH=amd64 go build -o mini-mcp-linux cmd/mini-mcp/main.go
GOOS=darwin GOARCH=arm64 go build -o mini-mcp-darwin cmd/mini-mcp/main.go
```

### Running Tests

```bash
# Run all tests with coverage and race detection
go test -race -cover ./...

# Run specific test packages
go test ./internal/shared/security/ -v -race -cover
go test ./internal/domain/file/ -v -race -cover

# Generate test coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Benchmark tests
go test -bench=. ./internal/shared/security/
```

### Development Workflow

1. **Write tests first** - Follow TDD practices
2. **Run tests frequently** - Use `go test -race` for development
3. **Check coverage** - Maintain high test coverage (>80%)
4. **Security review** - All security-related changes need review
5. **Performance testing** - Test with realistic loads

## 📚 Production-Ready Documentation

- [Architecture Documentation](docs/README_ARCHITECTURE.md) - Detailed architecture overview and design principles
- [Tools Documentation](docs/README_TOOLS.md) - Complete tool reference with security guidelines and usage examples
- [Proxmox Configuration Guide](docs/PROXMOX_CONFIG.md) - Proxmox server setup and authentication
- [Type Safety Documentation](docs/TYPE_SAFETY_IMPROVEMENTS.md) - Go 1.25 generics and type-safe implementations
- [Agent Development Guide](docs/AGENT.md) - Comprehensive context for AI agents and developers
- [Installer Documentation](internal/installer/README.md) - Detailed installer API and configuration
- [Security Guide](SECURITY.md) - Security hardening and best practices (if available)
- [Monitoring Guide](MONITORING.md) - Observability and alerting setup (if available)
- [Operations Guide](OPERATIONS.md) - Production operations and maintenance (if available)

## 🚀 Production-Ready Features Summary

This refactoring has transformed Mini MCP into a production-ready infrastructure management platform with:

### ✅ **Security Enhancements**
- Multi-layer security validation (command, path, input)
- Advanced input sanitization and path traversal protection
- Comprehensive security error tracking with stack traces
- Enterprise-grade authentication and authorization

### ✅ **Observability & Monitoring**
- Structured logging with correlation IDs and contextual metadata
- Advanced metrics collection (P95/P99 percentiles, error rates)
- Comprehensive health checks (system, filesystem, network, security)
- Built-in metrics endpoint for monitoring dashboards

### ✅ **Performance & Reliability**
- Graceful shutdown with proper resource cleanup
- Performance monitoring with alerting thresholds
- Dependency injection and proper error handling
- Comprehensive test coverage with race detection

### ✅ **Operational Excellence**
- Environment-specific configuration management
- Production-ready build optimizations
- Cross-platform compilation support
- Detailed documentation and operational guides

## 🤝 Contributing to Production-Ready Code

1. **Security First**: All changes must maintain or improve security posture
2. **Test Coverage**: Maintain >80% test coverage with race detection
3. **Performance Impact**: Consider performance implications of changes
4. **Documentation**: Update docs for any new features or configuration
5. **Code Review**: All security and performance changes require review
6. **Backwards Compatibility**: Maintain API compatibility where possible

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.
