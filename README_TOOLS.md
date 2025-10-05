# Mini MCP Tools Documentation

This document provides comprehensive documentation for all available MCP tools in the enhanced mini-mcp server.

## Overview

The mini-mcp server provides a comprehensive set of tools for infrastructure management, system monitoring, file operations, and command execution. All tools are designed with security, reliability, and production-readiness in mind.

## Consistent Patterns

All tools follow consistent patterns for:
- **Tool Registration**: Multi-line `mcp.Tool{}` structs with consistent formatting
- **Argument Types**: Standardized validation using direct validation functions
- **Error Handling**: Consistent error response format
- **Documentation**: Comprehensive descriptions with parameter explanations

### Tool Registration Pattern

```go
mcp.AddTool(server, &mcp.Tool{
    Name: "tool_name", 
    Description: "Comprehensive description with parameter explanations. Use 'param1' for description, 'param2' for description.",
}, func(ctx context.Context, req *mcp.CallToolRequest, args toolArgs) (*mcp.CallToolResult, any, error) {
    // Implementation
})
```

### Argument Type Pattern

```go
type ToolArgs struct {
    // RequiredParam is a required parameter with description
    RequiredParam string `json:"required_param"`
    // OptionalParam is an optional parameter with description
    OptionalParam string `json:"optional_param,omitempty"`
}

func (args *ToolArgs) Validate() error {
    // Validate required fields
    if err := validation.StringRequired("required_param", args.RequiredParam); err != nil {
        return err
    }
    
    // Validate optional fields if provided
    if args.OptionalParam != "" {
        // Additional validation
    }
    
    return nil
}
```

## Available Tools

### 1. execute_command

**Purpose**: Execute shell commands securely with allowlisting and sandboxing.

**Description**: This tool provides safe command execution with built-in security measures including command allowlisting, input sanitization, timeout controls, and output size limits. Ideal for system administration tasks, file operations, and process management.

**Security Features**:

- Command allowlisting (only pre-approved commands allowed)
- Input sanitization to prevent injection attacks
- Timeout controls to prevent hanging processes
- Output size limits to prevent memory exhaustion
- Working directory restrictions for isolation

**Allowed Commands**:

- System utilities: `ls`, `cat`, `head`, `tail`, `grep`, `find`, `wc`, `sort`, `uniq`
- Process management: `ps`, `top`
- System monitoring: `df`, `du`, `free`, `uptime`, `who`, `w`
- Version control: `git`
- Container management: `docker`
- Infrastructure: `consul`, `terraform`

**Parameters**:

- `command` (string, required): The command to execute (must be in allowlist)
- `timeout` (number, optional): Timeout in seconds (default: 30, max: 300)

**Examples**:

```json
{
  "command": "ls -la /tmp",
  "timeout": 10
}
```

```json
{
  "command": "ps aux --no-headers | head -20",
  "timeout": 20
}
```

```json
{
  "command": "df -h",
  "timeout": 15
}
```

### 2. file_operations

**Purpose**: Perform comprehensive file system operations with advanced path validation and security measures.

**Description**: Supports reading, writing, listing, and deleting files and directories with detailed metadata, permission handling, and automatic directory creation. Includes protection against directory traversal attacks and restricted access to system directories.

**Security Features**:

- Path validation and sanitization
- Directory traversal attack prevention
- Restricted access to system directories (`/etc`, `/var`, `/usr`)
- Automatic directory creation for write operations
- Permission handling and metadata preservation

**Operations**:

- `read`: Read file content with detailed metadata (size, permissions, timestamps)
- `write`: Create or overwrite files with automatic directory creation
- `list`: List directory contents with file information and total size
- `delete`: Remove files or directories with safety checks

**Path Restrictions**:

- Cannot access system directories: `/etc`, `/var`, `/usr`, `/bin`, `/sbin`
- Cannot use directory traversal: `../`, `../`
- Supports relative and absolute paths
- Automatic path normalization

**Parameters**:

- `operation` (string, required): The file operation to perform (`read`, `write`, `list`, `delete`)
- `path` (string, required): The file or directory path
- `content` (string, optional): Content to write (required for write operations)

**Examples**:

```json
{
  "operation": "read",
  "path": "config.json"
}
```

```json
{
  "operation": "write",
  "path": "logs/app.log",
  "content": "Application started at 2024-01-01 12:00:00"
}
```

```json
{
  "operation": "list",
  "path": "data/"
}
```

```json
{
  "operation": "delete",
  "path": "temp/cache.dat"
}
```

### 3. system_monitoring

**Purpose**: Get comprehensive system information and monitoring data for infrastructure management and performance analysis.

**Description**: Provides detailed metrics about processes, disk usage, memory consumption, network interfaces, and system uptime with formatted output and calculated statistics.

**Monitoring Capabilities**:

- `processes`: Detailed process information with CPU/memory usage, status, and command details
- `disk_usage`: Filesystem usage statistics with capacity, used space, and availability
- `memory_usage`: RAM consumption with total, used, free, and percentage calculations
- `network`: Network interface information including IP addresses, status, and configuration
- `uptime`: System uptime, load averages, boot time, and user session information

**Data Format**:

- All metrics return structured JSON with detailed information
- Numeric values include both raw and formatted representations
- Timestamps are provided in ISO 8601 format
- Percentages and ratios are calculated automatically

**Performance Considerations**:

- Process listing is limited to top 20 processes for performance
- Disk usage includes all mounted filesystems
- Memory calculations include swap and cache information
- Network data includes all active interfaces

**Parameters**:

- `metric` (string, required): The system metric to retrieve

**Examples**:

```json
{
  "metric": "processes"
}
```

```json
{
  "metric": "disk_usage"
}
```

```json
{
  "metric": "memory_usage"
}
```

```json
{
  "metric": "network"
}
```

```json
{
  "metric": "uptime"
}
```

### 4. health_check

**Purpose**: Check system health and service status with comprehensive diagnostics.

**Description**: Monitors disk space, memory usage, service availability, and overall system health. Provides detailed health reports with status indicators, performance metrics, and actionable recommendations.

**Health Checks**:

- Disk space monitoring with threshold alerts
- Memory usage analysis with percentage calculations
- Service availability checks
- Overall system health assessment

**Parameters**:

- `service` (string, optional): Service name to check (default: "all")

**Examples**:

```json
{
  "service": "all"
}
```

```json
{
  "service": "disk"
}
```

```json
{
  "service": "memory"
}
```

### 5. get_metrics

**Purpose**: Get application metrics and performance data for monitoring and analytics.

**Description**: Provides detailed insights into application performance, error rates, response times, connection statistics, and log analysis. Essential for performance monitoring and troubleshooting.

**Metrics Types**:

- `all`: Comprehensive metrics overview
- `logs`: Application log analysis
- `response_times`: API performance metrics
- `error_rates`: Error statistics and trends
- `connections`: Connection metrics and status

**Parameters**:

- `type` (string, optional): Type of metrics to retrieve

**Examples**:

```json
{
  "type": "all"
}
```

```json
{
  "type": "logs"
}
```

```json
{
  "type": "response_times"
}
```

### 6. docker_compose

**Purpose**: Execute Docker Compose operations with enhanced configuration support.

**Description**: Provides comprehensive docker compose functionality with support for multiple contexts, remote hosts, custom docker binaries, and custom compose files. Supports all standard docker-compose commands including up, down, ps, logs, restart, stop, and start. Ideal for managing containerized applications across different environments.

**Supported Commands**:

- `up`: Start services (with optional detached mode)
- `down`: Stop services (with optional volume removal)
- `ps`: List services
- `logs`: View service logs
- `restart`: Restart services
- `stop`: Stop services
- `start`: Start services

**Enhanced Features**:

- **Docker Context Support**: Switch between different Docker contexts
- **Remote Host Support**: Connect to remote Docker daemons
- **Custom Docker Binary**: Use custom docker binary paths
- **Custom Compose Files**: Specify custom docker-compose.yml files
- **Volume Management**: Remove volumes when stopping services

**Parameters**:

- `path` (string, required): Path to the docker-compose.yml file directory
- `command` (string, required): Docker Compose command to execute (up, down, ps, logs, restart, stop, start)
- `detached` (boolean, optional): Run in detached mode (for up command, default: false)
- `remove_volumes` (boolean, optional): Remove volumes when stopping (for down command, default: false)
- `context` (string, optional): Docker context to use
- `host` (string, optional): Docker host URL (e.g., "tcp://remote:2376")
- `docker_path` (string, optional): Path to docker binary
- `compose_file` (string, optional): Path to custom compose file

**Examples**:

```json
{
  "path": "/project",
  "command": "up",
  "detached": true
}
```

```json
{
  "path": "/project",
  "command": "down",
  "remove_volumes": true,
  "context": "production",
  "host": "tcp://remote:2376"
}
```

```json
{
  "path": "/project",
  "command": "ps"
}
```

```json
{
  "path": "/project",
  "command": "logs",
  "context": "production"
}
```

### 7. docker_swarm

**Purpose**: Get Docker Swarm cluster information with enhanced configuration support.

**Description**: Provides comprehensive Docker Swarm cluster information including node status, manager information, and cluster health. Supports remote contexts and custom docker configurations for managing distributed Docker environments.

**Enhanced Features**:

- **Docker Context Support**: Switch between different Docker contexts
- **Remote Host Support**: Connect to remote Docker daemons
- **Custom Docker Binary**: Use custom docker binary paths
- **Cluster Information**: Node status, manager information, cluster health
- **Non-destructive**: Read-only operations for safe monitoring

**Parameters**:

- `context` (string, optional): Docker context to use
- `host` (string, optional): Docker host URL (e.g., "tcp://remote:2376")
- `docker_path` (string, optional): Path to docker binary

**Examples**:

```json
{}
```

```json
{
  "context": "production",
  "host": "tcp://remote:2376"
}
```

```json
{
  "context": "production",
  "docker_path": "/usr/local/bin/docker"
}
```

### 8. docker_swarm

**Purpose**: Get Docker Swarm cluster information with enhanced configuration support.

**Description**: Provides comprehensive Docker Swarm cluster information including node status, manager information, and cluster health. Supports remote contexts and custom docker configurations for managing distributed Docker environments.

**Enhanced Features**:

- **Docker Context Support**: Switch between different Docker contexts
- **Remote Host Support**: Connect to remote Docker daemons
- **Custom Docker Binary**: Use custom docker binary paths
- **Cluster Information**: Node status, manager information, cluster health
- **Non-destructive**: Read-only operations for safe monitoring

**Parameters**:

- `context` (string, optional): Docker context to use
- `host` (string, optional): Docker host URL (e.g., "tcp://remote:2376")
- `docker_path` (string, optional): Path to docker binary

**Examples**:

```json
{}
```

```json
{
  "context": "production",
  "host": "tcp://remote:2376"
}
```

```json
{
  "context": "production",
  "docker_path": "/usr/local/bin/docker"
}
```

### 9. ssh

**Purpose**: Execute remote commands over SSH with security validation and comprehensive connection options.

**Description**: Provides secure SSH command execution with support for custom ports, private keys, users, timeouts, and connection validation. Essential for remote system administration, configuration management, and distributed system operations.

**Security Features**:

- Host validation and connection security
- Private key authentication support
- Configurable connection timeouts
- Input sanitization for commands
- Restricted to approved command patterns

**Connection Options**:

- **Host**: Target hostname or IP address
- **Port**: SSH port (default: 22, range: 1-65535)
- **User**: SSH username (default: root)
- **KeyPath**: Path to private key file for authentication
- **Timeout**: Connection timeout in seconds (default: 10, max: 300)

**Parameters**:

- `host` (string, required): Target hostname or IP address
- `command` (string, required): Command to execute on remote host
- `user` (string, optional): SSH username (default: "root")
- `port` (string, optional): SSH port number (default: "22")
- `key_path` (string, optional): Path to SSH private key file
- `timeout` (number, optional): Connection timeout in seconds (default: 10)

**Examples**:

```json
{
  "host": "server.example.com",
  "command": "uptime"
}
```

```json
{
  "host": "192.168.1.100",
  "command": "df -h",
  "user": "admin",
  "port": "2222",
  "key_path": "~/.ssh/id_rsa",
  "timeout": 30
}
```

```json
{
  "host": "remote-server",
  "command": "systemctl status nginx",
  "user": "deploy",
  "timeout": 15
}
```

### 10. port_process_tools

**Purpose**: Investigate and manage network ports and processes for debugging, monitoring, and system maintenance.

**Description**: Comprehensive tool for network port and process management including listing ports, finding processes using specific ports, killing processes, cleaning up occupied ports, and getting detailed port and process information. Essential for troubleshooting network issues and managing system resources.

**Operations**:

- `list_ports`: List network ports with optional filtering by state
- `list_processes`: List running processes with optional user filtering
- `kill_process`: Terminate a process by PID
- `find_port`: Find processes using a specific port
- `clean_ports`: Attempt to clean up occupied ports
- `port_info`: Get detailed information about a specific port
- `process_info`: Get detailed information about a specific process
- `network_stats`: Get network interface statistics

**Parameters**:

- `command` (string, required): Operation to perform
- `port` (number, optional): Port number for port-specific operations
- `process_id` (number, optional): Process ID for process-specific operations
- `process_name` (string, optional): Process name to search for
- `user` (string, optional): User to filter processes by
- `state` (string, optional): Port state to filter by (LISTEN, ESTABLISHED, etc.)

**Examples**:

```json
{
  "command": "list_ports"
}
```

```json
{
  "command": "list_ports",
  "state": "LISTEN"
}
```

```json
{
  "command": "find_port",
  "port": 8080
}
```

```json
{
  "command": "process_info",
  "process_id": 1234
}
```

```json
{
  "command": "kill_process",
  "process_id": 5678
}
```

```json
{
  "command": "list_processes",
  "user": "www-data"
}
```

```json
{
  "command": "network_stats"
}
```

## Error Handling

All tools provide comprehensive error handling with structured error responses:

### Error Types

- **Validation Errors**: Invalid input parameters or values
- **Security Errors**: Blocked commands or restricted paths
- **System Errors**: File not found, permission denied, etc.
- **Timeout Errors**: Command execution timeouts
- **Internal Errors**: Application-level errors

### Error Response Format

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

## Security Features

### Command Execution Security

- **Allowlisting**: Only pre-approved commands can be executed
- **Input Sanitization**: Prevents command injection attacks
- **Timeout Controls**: Prevents hanging processes
- **Output Limits**: Prevents memory exhaustion
- **Working Directory Restrictions**: Isolates command execution

### File Operations Security

- **Path Validation**: Ensures safe file paths
- **Directory Traversal Prevention**: Blocks `../` patterns
- **System Directory Restrictions**: Prevents access to sensitive directories
- **Permission Handling**: Respects file system permissions
- **Automatic Directory Creation**: Safe directory creation for writes

### System Monitoring Security

- **Read-Only Operations**: All monitoring operations are read-only
- **Limited Data Exposure**: Only necessary system information is exposed
- **Performance Optimization**: Limits data collection to prevent system impact

## Performance Considerations

### Command Execution

- Default timeout: 30 seconds
- Maximum timeout: 300 seconds (5 minutes)
- Output size limits to prevent memory issues
- Working directory isolation

### File Operations

- Automatic directory creation for write operations
- Metadata preservation during operations
- Efficient directory listing with size calculations
- Safe file deletion with existence checks

### System Monitoring

- Process listing limited to top 20 processes
- Efficient data collection using system commands
- Structured JSON responses for easy parsing
- Automatic unit conversions and formatting

## Usage Examples

### System Administration

```json
// Check system health
{
  "tool": "health_check",
  "params": {"service": "all"}
}

// Monitor disk usage
{
  "tool": "system_monitoring",
  "params": {"metric": "disk_usage"}
}

// List running processes
{
  "tool": "system_monitoring",
  "params": {"metric": "processes"}
}
```

### File Management

```json
// Read configuration file
{
  "tool": "file_operations",
  "params": {
    "operation": "read",
    "path": "config.json"
  }
}

// Write log file
{
  "tool": "file_operations",
  "params": {
    "operation": "write",
    "path": "logs/app.log",
    "content": "Application started at 2024-01-01 12:00:00"
  }
}

// List directory contents
{
  "tool": "file_operations",
  "params": {
    "operation": "list",
    "path": "data/"
  }
}
```

### Command Execution

```json
// Check disk space
{
  "tool": "execute_command",
  "params": {
    "command": "df -h",
    "timeout": 15
  }
}

// Get system processes
{
  "tool": "execute_command",
  "params": {
    "command": "ps aux --no-headers | head -10",
    "timeout": 20
  }
}

// Check git status
{
  "tool": "execute_command",
  "params": {
    "command": "git status",
    "timeout": 30
  }
}
```

### Docker Management

```json
// Start services with docker compose
{
  "tool": "docker_compose",
  "params": {
    "path": "/project",
    "command": "up",
    "detached": true,
    "context": "production"
  }
}

// Stop services with volume cleanup
{
  "tool": "docker_compose",
  "params": {
    "path": "/project",
    "command": "down",
    "remove_volumes": true,
    "context": "production"
  }
}

// Get swarm cluster information
{
  "tool": "docker_swarm",
  "params": {
    "context": "production",
    "host": "tcp://remote:2376"
  }
}

// List docker compose services
{
  "tool": "docker_compose",
  "params": {
    "path": "/project",
    "command": "ps"
  }
}

// Deploy to remote swarm cluster
{
  "tool": "docker_compose",
  "params": {
    "path": "/project",
    "command": "up",
    "detached": true,
    "context": "glpx-proxy",
    "compose_file": "/custom/compose.yml"
  }
}
```

## Best Practices

### Security

1. Always validate input parameters before use
2. Use appropriate timeouts for command execution
3. Avoid accessing sensitive system directories
4. Monitor tool usage and access patterns

### Performance

1. Use appropriate timeouts for long-running operations
2. Limit process listings to necessary information
3. Cache frequently accessed data when possible
4. Monitor tool performance and resource usage

### Reliability

1. Handle errors gracefully with proper error messages
2. Implement retry logic for transient failures
3. Log all operations for audit purposes
4. Validate responses before processing

## Configuration

The MCP server can be configured through environment variables:

- `ENVIRONMENT`: Set to "development" or "production"
- `LOG_LEVEL`: Set to "DEBUG", "INFO", "WARNING", or "ERROR"
- `SECURITY_WORKING_DIR`: Working directory for command execution
- `SECURITY_COMMAND_TIMEOUT`: Default command timeout
- `SECURITY_MAX_OUTPUT_SIZE`: Maximum command output size
- `SECURITY_ALLOWED_COMMANDS`: Comma-separated list of allowed commands

## Troubleshooting

### Common Issues

1. **Command Blocked**: Check if the command is in the allowlist
2. **Path Access Denied**: Verify the path is not in a restricted directory
3. **Timeout Errors**: Increase timeout value or optimize command
4. **Permission Errors**: Check file system permissions

### Debug Information

Enable debug logging by setting `LOG_LEVEL=DEBUG` to get detailed information about:

- Command execution details
- File operation steps
- System monitoring data collection
- Error context and stack traces

## Support

For issues, questions, or feature requests, please refer to the project documentation or create an issue in the project repository.
