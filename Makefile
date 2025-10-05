# Mini MCP Makefile

.PHONY: help build install uninstall status configure install-all clean test lint cli

# Default target
help:
	@echo "Mini MCP - Available targets:"
	@echo "  cli            - Build the CLI tool"
	@echo "  build          - Build the mini-mcp binary"
	@echo "  install        - Install the binary to system PATH"
	@echo "  uninstall      - Remove the binary from system PATH"
	@echo "  status         - Show installation status"
	@echo "  configure      - Update VS Code and Cursor settings"
	@echo "  install-all    - Complete installation with configuration"
	@echo "  clean          - Clean build artifacts"
	@echo "  test           - Run all tests"
	@echo "  lint           - Run linter"
	@echo "  dev-setup      - Setup development environment"

# Build the CLI tool
cli:
	@echo "Building CLI tool..."
	go build -o mini-mcp-cli ./cmd/mini-mcp-cli

# Build the main binary
build: cli
	@echo "Building mini-mcp..."
	@./mini-mcp-cli build

# Install the main binary
install: cli
	@echo "Installing mini-mcp..."
	@./mini-mcp-cli install

# Install with configuration
install-all: cli
	@echo "Installing and configuring mini-mcp..."
	@./mini-mcp-cli install --configure

# Uninstall the main binary
uninstall: cli
	@echo "Uninstalling mini-mcp..."
	@./mini-mcp-cli uninstall

# Show installation status
status: cli
	@echo "Checking installation status..."
	@./mini-mcp-cli status

# Configure VS Code and Cursor settings
configure: cli
	@echo "Configuring VS Code and Cursor settings..."
	@./mini-mcp-cli configure

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f mini-mcp
	rm -f mini-mcp-cli
	go clean

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run ./...

# Development setup
dev-setup: cli
	@echo "Development setup complete!"
	@echo "You can now use './mini-mcp-cli' to manage mini-mcp installation"

# Quick install for development
dev-install: install-all
	@echo "Development installation complete!"
	@echo "mini-mcp is now available in your PATH and configured in VS Code/Cursor"
