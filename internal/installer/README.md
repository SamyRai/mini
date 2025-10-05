# Mini MCP Installer

This package provides functionality to build, install, and configure the mini-mcp tool on your system.

## Features

- **Build**: Compile the mini-mcp binary
- **Install**: Install the binary to system PATH
- **Uninstall**: Remove the binary from system PATH
- **Configure**: Update VS Code and Cursor settings automatically
- **Status**: Check installation status
- **Cross-platform**: Works on Windows, macOS, and Linux

## Installation Paths

The installer automatically determines the appropriate installation path based on your operating system:

- **Windows**: `%USERPROFILE%\bin\mini-mcp.exe`
- **macOS/Linux**: `/usr/local/bin/mini-mcp`
- **Other**: `$HOME/bin/mini-mcp`

## Usage

### Command Line Tool

```bash
# Build the CLI tool
go build -o mini-mcp-cli ./cmd/mini-mcp-cli

# Install mini-mcp
./mini-mcp-cli install

# Install and configure VS Code/Cursor
./mini-mcp-cli install --configure

# Check status
./mini-mcp-cli status

# Uninstall
./mini-mcp-cli uninstall
```

### Makefile

```bash
# Quick installation
make install-all

# Build only
make build

# Check status
make status

# Uninstall
make uninstall
```

### Shell Script (Legacy - Use CLI instead)

```bash
# Full installation
make install-all

# Build CLI only
make cli

# Check status
make status

# Uninstall
make uninstall
```

## API Usage

```go
package main

import (
    "mini-mcp/internal/installer"
    "mini-mcp/internal/shared/logging"
)

func main() {
    // Set up logging
    logging.SetGlobalLogger(logging.NewLogger(os.Stdout, logging.INFO))
    
    // Create installer
    inst := installer.NewInstaller("/path/to/project")
    
    // Build and install
    if err := inst.Build(); err != nil {
        log.Fatal(err)
    }
    
    if err := inst.Install(); err != nil {
        log.Fatal(err)
    }
    
    // Configure editors
    inst.UpdateVSCodeSettings()
    inst.UpdateCursorSettings()
    
    // Check status
    status := inst.GetStatus()
    fmt.Printf("Installed: %v\n", status["is_installed"])
}
```

## Configuration

The installer automatically updates VS Code and Cursor settings with the following MCP server configuration:

```json
{
  "mcpServers": {
    "mini-mcp": {
      "command": "/usr/local/bin/mini-mcp",
      "args": [],
      "env": {
        "ENVIRONMENT": "production",
        "LOG_LEVEL": "INFO",
        "PORT": ":8080"
      },
      "cwd": "/path/to/project",
      "description": "Enhanced Mini MCP Infrastructure Management Tool"
    }
  }
}
```

## Requirements

- Go 1.21 or later
- Write permissions to installation directory
- Write permissions to VS Code/Cursor settings directory

## Error Handling

The installer provides comprehensive error handling and logging:

- **Build errors**: Clear error messages for compilation failures
- **Permission errors**: Helpful messages for permission issues
- **Path errors**: Validation of project and installation paths
- **Configuration errors**: Graceful handling of settings file issues

## Logging

The installer uses structured logging with the following levels:

- **DEBUG**: Detailed debugging information
- **INFO**: General information about operations
- **WARN**: Warning messages for non-critical issues
- **ERROR**: Error messages for critical failures

## Examples

### Development Setup

```bash
# Clone the repository
git clone https://github.com/your-org/mini-mcp.git
cd mini-mcp

# Quick development setup
make dev-setup

# Install with configuration
make install-all
```

### Production Installation

```bash
# Download and run CLI
wget https://github.com/your-org/mini-mcp/releases/latest/download/mini-mcp-cli
chmod +x mini-mcp-cli
./mini-mcp-cli install --configure
```

### Custom Installation

```bash
# Install to custom location (requires installer package modification)
./mini-mcp-cli install --project-root /custom/path

# Configure only
./mini-mcp-cli configure
```

## Troubleshooting

### Permission Denied

If you get permission denied errors:

```bash
# On macOS/Linux, use sudo for system-wide installation
sudo ./mini-mcp-cli install

# Or install to user directory (requires installer package modification)
export INSTALL_PATH=$HOME/bin
./mini-mcp-cli install
```

### VS Code/Cursor Settings Not Updated

If settings are not updated automatically:

1. Check if the settings file exists
2. Verify write permissions
3. Manually update the settings file
4. Restart VS Code/Cursor

### Binary Not Found in PATH

If the binary is not found after installation:

1. Check the installation path: `./mini-mcp-cli status`
2. Add the installation directory to your PATH
3. Restart your terminal
4. Verify with `which mini-mcp`

## Contributing

When contributing to the installer:

1. Test on all supported platforms
2. Update documentation for new features
3. Add appropriate error handling
4. Include logging for debugging
5. Update tests for new functionality
