# Mini MCP - Architecture Documentation

**📚 [← Back to Main README](../README.md) | [🛠️ Tools Documentation](README_TOOLS.md) | [⚙️ Proxmox Configuration](PROXMOX_CONFIG.md) | [🤖 Agent Guide](AGENT.md)**

## Overview

Mini MCP is a production-ready infrastructure management platform that provides both MCP server and CLI interfaces for secure command execution, file operations, and system monitoring. The project follows clean architecture principles with domain-driven design, comprehensive type safety, and enterprise-grade security.

## Architecture

### Directory Structure

```
mini-mcp/
├── cmd/                          # Application entry points
│   ├── mini-mcp/                # MCP server and CLI application
│   │   └── main.go             # Main entry point supporting both server and CLI modes
│   └── mini-mcp-cli/           # CLI installer tool
│       ├── main.go             # CLI installer entry point
│       └── cmd/                # CLI subcommands
│           ├── build.go        # Build command
│           ├── configure.go    # Configuration command
│           ├── install.go      # Installation command
│           ├── status.go       # Status command
│           ├── uninstall.go    # Uninstallation command
│           └── utils.go        # CLI utilities
├── internal/                    # Internal application code
│   ├── domain/                 # Business logic and domain models
│   │   ├── command/            # Command execution domain
│   │   │   ├── types.go        # Command domain types
│   │   │   ├── repository.go   # Command repository interface
│   │   │   ├── service.go      # Command domain service
│   │   │   └── errors.go       # Command domain errors
│   │   ├── file/               # File operations domain
│   │   │   ├── types.go        # File domain types
│   │   │   ├── repository.go   # File repository interface
│   │   │   ├── service.go      # File domain service
│   │   │   └── errors.go       # File domain errors
│   │   ├── system/             # System monitoring domain
│   │   └── infrastructure/     # Infrastructure management domain
│   │       ├── docker/         # Docker operations
│   │       ├── ssh/            # SSH operations
│   │       └── git/            # Git operations
│   ├── application/            # Application services (use cases)
│   │   ├── command/            # Command application services
│   │   │   └── service.go      # Command application service
│   │   ├── file/               # File application services
│   │   │   └── service.go      # File application service
│   │   ├── system/             # System application services
│   │   └── infrastructure/     # Infrastructure application services
│   ├── infrastructure/         # External concerns
│   │   ├── cli/                # CLI interface
│   │   │   └── handlers/       # CLI handlers
│   │   │       ├── command_handler.go
│   │   │       └── file_handler.go
│   │   ├── http/               # HTTP interface
│   │   │   ├── handlers/       # HTTP handlers
│   │   │   │   └── handlers.go
│   │   │   ├── middleware/     # HTTP middleware
│   │   │   └── routes/         # HTTP routes
│   │   └── persistence/        # Data persistence
│   └── shared/                 # Shared utilities
│       ├── auth/               # Authentication
│       ├── config/             # Configuration
│       ├── errors/             # Error handling
│       ├── logging/            # Logging
│       ├── security/           # Security
│       └── validation/         # Validation
├── main.go                     # Simple entry point
└── README_ARCHITECTURE.md      # This file
```

## Architecture Layers

### 1. Domain Layer (`internal/domain/`)

The domain layer contains the core business logic and domain models. Each domain is self-contained with its own types, services, and repositories.

#### Command Domain

- **Purpose**: Handles command execution with security and validation
- **Key Components**:
  - `Command`: Domain entity representing a command to execute
  - `Repository`: Interface for command execution operations
  - `Service`: Business logic for command execution
  - `Status`: Enum for command execution status

#### File Domain

- **Purpose**: Handles file system operations with security
- **Key Components**:
  - `FileInfo`: Domain entity for file information
  - `Repository`: Interface for file operations
  - `Service`: Business logic for file operations
  - `Operation`: Enum for file operation types

#### Docker Domain

- **Purpose**: Handles Docker operations with enhanced configuration support
- **Key Components**:
  - `DockerConfig`: Configuration for Docker operations (context, host, path)
  - `DockerComposeUpArgs`: Arguments for compose up operations
  - `DockerComposeDownArgs`: Arguments for compose down operations
  - `DockerSwarmInfoArgs`: Arguments for swarm information
  - `Repository`: Interface for Docker operations
  - `Service`: Business logic for Docker operations

### 2. Application Layer (`internal/application/`)

The application layer contains use cases and orchestrates between domains. It provides a clean API for external interfaces.

#### Command Application Service

- **Purpose**: Coordinates command execution use cases
- **Responsibilities**:
  - Execute commands with validation
  - Handle command timeouts
  - Provide error handling

#### File Application Service

- **Purpose**: Coordinates file operation use cases
- **Responsibilities**:
  - Read, write, list, and delete files
  - Validate file operations
  - Handle file metadata

#### Docker Application Service

- **Purpose**: Coordinates Docker operation use cases
- **Responsibilities**:
  - Execute docker compose operations (up/down)
  - Get Docker Swarm cluster information
  - Handle Docker context and host configuration
  - Manage custom docker binary paths
  - Support remote Docker operations

### 3. Infrastructure Layer (`internal/infrastructure/`)

The infrastructure layer handles external concerns like CLI, HTTP, and persistence.

#### CLI Infrastructure

- **Purpose**: Provides command-line interface
- **Components**:
  - `CommandHandler`: Handles CLI command execution
  - `FileHandler`: Handles CLI file operations

#### HTTP Infrastructure

- **Purpose**: Provides HTTP API interface
- **Components**:
  - `Handlers`: HTTP request handlers
  - `Middleware`: HTTP middleware (auth, logging)
  - `Routes`: HTTP route definitions

### 4. Shared Layer (`internal/shared/`)

The shared layer contains cross-cutting concerns used throughout the application.

#### Components

- **Auth**: Authentication and authorization
- **Config**: Configuration management
- **Errors**: Error handling utilities
- **Logging**: Structured logging
- **Security**: Security utilities
- **Validation**: Input validation

## Design Principles

### 1. Single Responsibility Principle (SRP)

Each file and package has a single, well-defined responsibility:

- Domain files focus on business logic and domain models
- Application files focus on use cases and orchestration
- Infrastructure files focus on external concerns (CLI, HTTP, persistence)
- Shared files focus on cross-cutting concerns (auth, config, logging, security)

### 2. Dependency Inversion

- High-level modules (application) don't depend on low-level modules (infrastructure)
- Both depend on abstractions (interfaces)
- Abstractions are defined in the domain layer
- Infrastructure adapts to domain needs, not vice versa

### 3. Clean Architecture

- Dependencies point inward toward the domain
- Domain layer has no external dependencies
- Infrastructure adapts to domain needs
- Clear separation between business logic and external concerns

### 4. Type Safety

- Full Go 1.25 generics implementation
- Compile-time type checking throughout
- No unsafe type assertions in critical paths
- Type-safe validation and error handling

### 5. Security-First Design

- Multi-layer security validation
- Command allowlisting and path validation
- Input sanitization and output limits
- Comprehensive audit logging

### 6. File Size Guidelines

- **Domain files**: Max 200 lines
- **Service files**: Max 150 lines
- **Handler files**: Max 100 lines
- **Main files**: Max 50 lines
- **Type files**: Max 100 lines

## Usage

### MCP Server Mode

```bash
# Build unified MCP binary
go build -o mini-mcp cmd/mini-mcp/main.go

# Run as MCP server (for Cursor and other MCP clients)
./mini-mcp -mode=server

# Or use go run
go run cmd/mini-mcp/main.go -mode=server
```

### CLI Mode

```bash
# Run as CLI tool
./mini-mcp -mode=cli

# Or use go run
go run cmd/mini-mcp/main.go -mode=cli
```

### Entry Point

```bash
# Show help
go run . help

# Run CLI mode
go run . cli

# Run server mode (default)
go run . server
```

## Benefits of New Architecture

### 1. Maintainability

- Clear separation of concerns
- Easy to locate and modify code
- Reduced coupling between components

### 2. Testability

- Business logic is isolated in domain layer
- Easy to mock dependencies
- Clear interfaces for testing

### 3. Scalability

- Easy to add new domains
- Clear boundaries for team collaboration
- Modular design supports growth

### 4. Flexibility

- Easy to change interfaces (CLI/HTTP)
- Business logic is reusable
- Clear extension points

### 5. Readability

- Smaller, focused files
- Clear naming conventions
- Consistent structure

## Architecture Evolution

The current architecture evolved from earlier implementations to provide better separation of concerns and maintainability:

### Design Decisions

1. **Unified Binary**: Single binary supporting both MCP server and CLI modes
2. **Clean Architecture**: Clear separation between domain, application, and infrastructure layers
3. **Domain-Driven Design**: Business logic organized by domain (command, file, system, infrastructure)
4. **Shared Components**: Common utilities for auth, config, logging, security, and validation
5. **Type Safety**: Comprehensive use of Go generics and type-safe interfaces

### Architecture Benefits

1. **Unified Interface**: Single binary reduces complexity and maintenance overhead
2. **Clean Separation**: Domain logic isolated from infrastructure concerns
3. **Type Safety**: Compile-time guarantees with Go 1.25 generics
4. **Testability**: Clear interfaces enable comprehensive testing
5. **Maintainability**: Well-organized code with consistent patterns

## Current Implementation Status

### ✅ Completed Features

- **Domain Implementation**: All domain services implemented with proper interfaces
- **Type Safety**: Full Go 1.25 generics implementation with type-safe validation
- **Security**: Multi-layer security with command allowlisting and path validation
- **Infrastructure Tools**: Docker, SSH, Git, and port management capabilities
- **Monitoring**: Comprehensive health checks and metrics collection
- **Documentation**: Complete documentation suite across multiple files

### 🔄 Active Development Areas

- **Testing**: Unit and integration tests being added incrementally
- **Performance**: Ongoing optimization and monitoring improvements
- **Security**: Continuous security enhancements and hardening
- **Features**: Additional infrastructure management capabilities

## Contributing

When contributing to this project:

1. **Follow the architecture**: Place code in the appropriate layer
2. **Respect file size limits**: Split large files into smaller ones
3. **Use interfaces**: Define contracts in the domain layer
4. **Write tests**: Ensure new code is well-tested
5. **Update documentation**: Keep docs in sync with code changes

## Conclusion

The current architecture provides a robust, production-ready foundation for the Mini MCP infrastructure management platform. Built on clean architecture principles with comprehensive type safety and security measures, it offers:

- **Unified Interface**: Single binary supporting both MCP server and CLI modes
- **Domain-Driven Design**: Well-organized business logic with clear boundaries
- **Type Safety**: Full compile-time type checking with Go 1.25 generics
- **Security-First**: Multi-layer security with comprehensive validation
- **Production Monitoring**: Built-in metrics, health checks, and observability
- **Infrastructure Management**: Comprehensive Docker, SSH, Git, and system management tools

The architecture successfully balances maintainability, testability, and operational requirements while providing a solid foundation for continued development and feature expansion.
