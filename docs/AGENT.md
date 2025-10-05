# Mini MCP Agent Development Guide

**📚 [← Back to Main README](../README.md) | [🏗️ Architecture Documentation](README_ARCHITECTURE.md) | [🛠️ Tools Documentation](README_TOOLS.md) | [⚙️ Proxmox Configuration](PROXMOX_CONFIG.md) | [🔒 Type Safety](TYPE_SAFETY_IMPROVEMENTS.md)**

This guide provides comprehensive context for AI agents and developers working with the Mini MCP infrastructure management platform.

## 🤖 Agent Context Overview

Mini MCP is a production-ready Model Context Protocol (MCP) server that provides secure infrastructure management capabilities through a unified binary supporting both server and CLI modes. It's designed for AI agents to safely execute infrastructure operations with comprehensive security, monitoring, and type safety.

## 🏗️ Project Architecture

### **Core Components**
- **Unified Binary**: Single `mini-mcp` binary supporting both MCP server and CLI modes
- **Clean Architecture**: Domain-driven design with clear separation of concerns
- **Type Safety**: Full Go 1.25 generics implementation with compile-time guarantees
- **Security-First**: Multi-layer security with command allowlisting and path validation
- **Production Monitoring**: Built-in metrics, health checks, and observability

### **Directory Structure**
```
mini-mcp/
├── README.md                    # Main project documentation
├── docs/                        # Documentation directory
│   ├── AGENT.md                # This file - Agent development guide
│   ├── README_ARCHITECTURE.md # Architecture documentation
│   ├── README_TOOLS.md         # Complete tool reference
│   ├── PROXMOX_CONFIG.md       # Proxmox configuration
│   └── TYPE_SAFETY_IMPROVEMENTS.md # Type safety documentation
├── cmd/                         # Application entry points
│   ├── mini-mcp/               # Main MCP server and CLI application
│   └── mini-mcp-cli/           # CLI installer tool
├── internal/                    # Internal application code
│   ├── application/            # Application services (use cases)
│   ├── domain/                 # Business logic and domain models
│   ├── handlers/               # Business logic handlers
│   ├── health/                 # Health checks and monitoring
│   ├── installer/              # Installation and configuration
│   ├── proxmox/                # Proxmox integration
│   ├── registry/               # Tool registry and execution
│   ├── server/                 # MCP protocol server
│   ├── shared/                 # Shared utilities
│   │   ├── auth/               # Authentication and authorization
│   │   ├── config/             # Configuration management
│   │   ├── errors/             # Structured error handling
│   │   ├── logging/            # Structured logging
│   │   ├── security/           # Security validation
│   │   └── validation/         # Input validation
│   ├── tools/                  # Tool implementations
│   └── types/                  # Type definitions
└── scripts/                     # Test and utility scripts
```

## 🛠️ Available MCP Tools

### **Core Infrastructure Tools**
- **`execute_command`** - Secure command execution with allowlisting
- **`file_operations`** - File system operations with path validation
- **`system_monitoring`** - System information and resource monitoring
- **`health_check`** - Comprehensive health diagnostics
- **`get_metrics`** - Application metrics and performance data

### **Infrastructure Management Tools**
- **`docker_compose`** - Docker Compose operations with context support
- **`docker_swarm`** - Docker Swarm cluster information
- **`ssh`** - Remote command execution over SSH
- **`port_process_tools`** - Network port and process management
- **`git`** - Git repository operations
- **`docs`** - Documentation fetching

### **Security Features**
- **Command Allowlisting**: Only pre-approved commands can be executed
- **Path Validation**: Prevents directory traversal attacks
- **Input Sanitization**: Comprehensive sanitization preventing injection
- **Timeout Controls**: Prevents hanging processes
- **Working Directory Restrictions**: Isolates command execution

## 🚀 Quick Start for Agents

### **1. Installation**
```bash
# Clone and build
git clone <repository-url>
cd mini-mcp
make install-all
```

### **2. MCP Server Mode (Recommended for Agents)**
```bash
# Run as MCP server for AI agent integration
mini-mcp -mode=server
```

### **3. Cursor Configuration**
Add to `~/.cursor/mcp.json`:
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

## 🔧 Development Context

### **Go Version Requirements**
- **Go 1.25 or later** - Required for full generics support and type safety
- **Docker** (optional) - For Docker-related tools
- **Docker Compose v2** (optional) - For compose operations

### **Key Development Patterns**

#### **Type-Safe Validation**
```go
// Example: Type-safe command validation
func ValidateCommand(cmd string) error {
    if !isAllowedCommand(cmd) {
        return SecurityError{
            Code: "COMMAND_BLOCKED",
            Message: "Command not in allowlist",
            Details: map[string]interface{}{
                "command": cmd,
                "allowed_commands": getAllowedCommands(),
            },
        }
    }
    return nil
}
```

#### **Structured Error Handling**
```go
// Example: Structured error responses
type TypedError[T any] struct {
    Code      ErrorCode `json:"code"`
    Message   string    `json:"message"`
    Details   T         `json:"details,omitempty"`
    Timestamp time.Time `json:"timestamp"`
    Retryable bool      `json:"retryable"`
}
```

#### **Security-First Design**
```go
// Example: Multi-layer security validation
func ExecuteCommand(cmd string, timeout time.Duration) (*CommandResult, error) {
    // 1. Command allowlist validation
    if err := validateCommand(cmd); err != nil {
        return nil, err
    }
    
    // 2. Input sanitization
    sanitizedCmd := sanitizeInput(cmd)
    
    // 3. Path validation
    if err := validateWorkingDirectory(); err != nil {
        return nil, err
    }
    
    // 4. Execute with timeout
    return executeWithTimeout(sanitizedCmd, timeout)
}
```

## 📊 Monitoring and Observability

### **Built-in Metrics**
- Request counts and response times
- Error rates and performance percentiles (P95, P99)
- System resource usage and health status
- Security validation statistics

### **Health Checks**
- Disk space monitoring with threshold alerts
- Memory usage analysis with percentage calculations
- Service availability checks
- Overall system health assessment

### **Structured Logging**
```json
{
  "timestamp": "2024-01-01T12:00:00Z",
  "level": "INFO",
  "tool": "execute_command",
  "duration": "150ms",
  "success": true,
  "user_agent": "cursor-mcp-client",
  "request_id": "req-12345"
}
```

## 🔒 Security Considerations for Agents

### **Command Execution Security**
- **Allowlisted Commands Only**: Only pre-approved commands can be executed
- **Timeout Enforcement**: Commands timeout after 30 seconds (configurable, max 300s)
- **Output Size Limits**: Prevents memory exhaustion from large outputs
- **Working Directory Isolation**: Commands run in restricted directories

### **File Operations Security**
- **Path Validation**: Prevents directory traversal attacks (`../` patterns blocked)
- **System Directory Restrictions**: Cannot access `/etc`, `/var`, `/usr`, etc.
- **Permission Handling**: Respects file system permissions
- **Automatic Directory Creation**: Safe directory creation for write operations

### **Network Security**
- **SSH Key Authentication**: Supports private key authentication
- **Connection Timeouts**: Configurable SSH connection timeouts
- **Host Validation**: Validates target hosts before connection

## 🎯 Agent Usage Patterns

### **Infrastructure Management**
```json
// Start Docker services
{
  "tool": "docker_compose",
  "arguments": {
    "path": "/project",
    "command": "up",
    "detached": true,
    "context": "production"
  }
}

// Monitor system health
{
  "tool": "health_check",
  "arguments": {
    "service": "all"
  }
}

// Execute secure command
{
  "tool": "execute_command",
  "arguments": {
    "command": "df -h",
    "timeout": 15
  }
}
```

### **File Operations**
```json
// Read configuration
{
  "tool": "file_operations",
  "arguments": {
    "operation": "read",
    "path": "config.json"
  }
}

// Write log file
{
  "tool": "file_operations",
  "arguments": {
    "operation": "write",
    "path": "logs/app.log",
    "content": "Application started"
  }
}
```

### **Remote Operations**
```json
// SSH to remote server
{
  "tool": "ssh",
  "arguments": {
    "host": "server.example.com",
    "command": "uptime",
    "user": "admin",
    "key_path": "~/.ssh/id_rsa"
  }
}
```

## 🧪 Testing and Development

### **Running Tests**
```bash
# Run all tests with coverage
go test -race -cover ./...

# Run specific test packages
go test ./internal/shared/security/ -v -race -cover

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### **MCP Protocol Testing**
```bash
# Test initialization
echo '{"type": "initialize"}' | mini-mcp -mode=server

# Test tools listing
echo '{"type": "tools/list"}' | mini-mcp -mode=server

# Test command execution
echo '{"type": "tools/call", "data": {"name": "run", "arguments": {"command": "ls -la"}}}' | mini-mcp -mode=server
```

### **Security Testing**
```bash
# Test dangerous command blocking
echo '{"type": "tools/call", "data": {"name": "run", "arguments": {"command": "rm -rf /"}}}' | mini-mcp -mode=server

# Test path traversal protection
echo '{"type": "tools/call", "data": {"name": "ls", "arguments": {"path": "../../../etc"}}}' | mini-mcp -mode=server
```

## 🔧 Configuration for Agents

### **Environment Variables**
```bash
# Environment
export ENVIRONMENT=development

# Logging
export LOG_LEVEL=DEBUG

# Security
export SECURITY_WORKING_DIR=/tmp
export SECURITY_COMMAND_TIMEOUT=30s
export SECURITY_MAX_OUTPUT_SIZE=1048576
export SECURITY_ALLOWED_COMMANDS=ls,cat,head,tail,grep,find,wc,sort,uniq,ps,top,df,du,free,uptime,who,w,git,docker,nomad,consul,terraform
export SECURITY_ALLOWED_PATHS=/tmp,/var/log,/proc
export SECURITY_BLOCKED_PATHS=/etc/passwd,/etc/shadow,/root,/home

# Authentication
export AUTH_RATE_LIMITING=100ms
export AUTH_MAX_REQUESTS=1000
export AUTH_WINDOW_SIZE=1h
export AUTH_IP_WHITELIST=127.0.0.1,::1
export AUTH_API_KEYS=key1:value1,key2:value2

# Performance
export PERF_MAX_CONCURRENT_REQUESTS=100
export PERF_REQUEST_TIMEOUT=30s
export PERF_CACHE_ENABLED=false
export PERF_CACHE_TTL=5m
```

### **Proxmox Configuration**
Create `proxmox-auth.yaml`:
```yaml
proxmox:
  host: "your-proxmox-server.com"
  token_name: "user@pam!token-name"
  token_value: "your-token-value"
  verify_ssl: false
  timeout: 30
  node: "your-node-name"
```

## 🚨 Error Handling for Agents

### **Common Error Types**
- **Validation Errors**: Invalid input parameters or values
- **Security Errors**: Blocked commands or restricted paths
- **System Errors**: File not found, permission denied, etc.
- **Timeout Errors**: Command execution timeouts
- **Internal Errors**: Application-level errors

### **Error Response Format**
```json
{
  "code": "ERROR_CODE",
  "message": "Human-readable error message",
  "details": {
    "field": "Additional context information"
  },
  "retryable": false,
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### **Agent Error Handling Best Practices**
1. **Always check error responses** before processing results
2. **Handle retryable errors** with exponential backoff
3. **Log security errors** for audit purposes
4. **Validate responses** before using data
5. **Implement circuit breakers** for failing services

## 📈 Performance Considerations

### **Command Execution Limits**
- Default timeout: 30 seconds
- Maximum timeout: 300 seconds (5 minutes)
- Output size limits to prevent memory issues
- Working directory isolation

### **File Operations Limits**
- Automatic directory creation for write operations
- Metadata preservation during operations
- Efficient directory listing with size calculations
- Safe file deletion with existence checks

### **System Monitoring Limits**
- Process listing limited to top 20 processes
- Efficient data collection using system commands
- Structured JSON responses for easy parsing
- Automatic unit conversions and formatting

## 🔄 Development Workflow

### **Building**
```bash
# Production build with optimizations
go build -ldflags="-s -w" -o mini-mcp cmd/mini-mcp/main.go

# Development build with race detection
go build -race -o mini-mcp cmd/mini-mcp/main.go

# Cross-compile for different platforms
GOOS=linux GOARCH=amd64 go build -o mini-mcp-linux cmd/mini-mcp/main.go
GOOS=darwin GOARCH=arm64 go build -o mini-mcp-darwin cmd/mini-mcp/main.go
```

### **Installation**
```bash
# Quick installation
make install-all

# Manual installation
sudo cp mini-mcp /usr/local/bin/mini-mcp
sudo chmod +x /usr/local/bin/mini-mcp
```

### **Development Mode**
```bash
# Run directly without installation
go run cmd/mini-mcp/main.go -mode=server
go run cmd/mini-mcp/main.go -mode=cli
```

## 🎯 Agent Integration Best Practices

### **1. Security First**
- Always validate input parameters
- Use appropriate timeouts for operations
- Avoid accessing sensitive system directories
- Monitor tool usage and access patterns

### **2. Performance Optimization**
- Use appropriate timeouts for long-running operations
- Limit process listings to necessary information
- Cache frequently accessed data when possible
- Monitor tool performance and resource usage

### **3. Error Handling**
- Handle errors gracefully with proper error messages
- Implement retry logic for transient failures
- Log all operations for audit purposes
- Validate responses before processing

### **4. Monitoring and Observability**
- Use built-in metrics for performance monitoring
- Implement health checks for system status
- Use structured logging for debugging
- Monitor security events and access patterns

## 📚 Additional Resources

- **[Architecture Documentation](README_ARCHITECTURE.md)** - Detailed architecture overview
- **[Tools Documentation](README_TOOLS.md)** - Complete tool reference with examples
- **[Proxmox Configuration](PROXMOX_CONFIG.md)** - Proxmox server setup guide
- **[Type Safety Documentation](TYPE_SAFETY_IMPROVEMENTS.md)** - Go 1.25 generics implementation
- **[Main README](../README.md)** - Project overview and quick start

## 🤝 Contributing

When contributing to Mini MCP:

1. **Follow the architecture** - Place code in appropriate layers
2. **Respect file size limits** - Split large files into smaller ones
3. **Use interfaces** - Define contracts in the domain layer
4. **Write tests** - Ensure new code is well-tested
5. **Update documentation** - Keep docs in sync with code changes
6. **Security review** - All security-related changes need review
7. **Performance testing** - Test with realistic loads

## 🚀 Production Deployment

### **Environment-Specific Configurations**
- **Development**: Relaxed security, detailed logging, extended timeouts
- **Staging**: Balanced security and performance, moderate logging
- **Production**: Maximum security, optimized performance, minimal logging

### **Monitoring Setup**
- Enable structured logging with correlation IDs
- Configure metrics collection and alerting
- Set up health check endpoints
- Implement security monitoring and audit logging

This guide provides comprehensive context for AI agents and developers working with Mini MCP. The platform is designed to be secure, reliable, and production-ready while providing powerful infrastructure management capabilities through a clean, type-safe API.
