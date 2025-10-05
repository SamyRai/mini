# Mini MCP - Architecture Documentation

## Overview

Mini MCP is a secure infrastructure management tool that provides both CLI and HTTP server interfaces for executing commands, managing files, and monitoring systems. The project follows a clean, domain-driven architecture with clear separation of concerns.

## Architecture

### Directory Structure

```
mini-mcp/
├── cmd/                          # Application entry points
│   ├── mini-mcp/                # CLI application
│   │   └── main.go             # CLI entry point
│   ├── mcp-server/             # HTTP server application
│   │   └── main.go             # Server entry point
│   └── test-client/            # Test client
│       └── main.go             # Test client entry point
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

- Domain files focus on business logic
- Application files focus on use cases
- Infrastructure files focus on external concerns

### 2. Dependency Inversion

- High-level modules (application) don't depend on low-level modules (infrastructure)
- Both depend on abstractions (interfaces)
- Abstractions are defined in the domain layer

### 3. Clean Architecture

- Dependencies point inward
- Domain layer has no external dependencies
- Infrastructure adapts to domain needs

### 4. File Size Guidelines

- **Domain files**: Max 200 lines
- **Service files**: Max 150 lines
- **Handler files**: Max 100 lines
- **Main files**: Max 50 lines
- **Type files**: Max 100 lines

## Usage

### CLI Mode

```bash
# Build CLI binary
go build -o mini-mcp cmd/mini-mcp/main.go

# Run CLI
./mini-mcp

# Or use go run
go run cmd/mini-mcp/main.go
```

### Server Mode

```bash
# Build server binary
go build -o mcp-server cmd/mcp-server/main.go

# Run server
./mcp-server

# Or use go run
go run cmd/mcp-server/main.go
```

### Entry Point

```bash
# Show help
go run . help

# Run CLI mode
go run . cli

# Run server mode
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

## Migration from Old Architecture

The new architecture addresses several issues from the old structure:

### Problems Solved

1. **Large main.go (464 lines)** → Split into focused entry points
2. **Mixed responsibilities** → Clear layer separation
3. **Inconsistent file sizes** → Enforced size guidelines
4. **No domain boundaries** → Domain-driven organization
5. **Mixed CLI/Server logic** → Separate interfaces

### Migration Strategy

1. ✅ Create new directory structure
2. ✅ Extract domain models and services
3. ✅ Create application services
4. ✅ Refactor CLI and HTTP handlers
5. ✅ Update main files
6. ✅ Add comprehensive tests
7. ✅ Update documentation
8. ✅ Remove legacy and duplicate code
9. ✅ Apply architectural patterns (DDD, Factory, Strategy, Facade, Observer)

## Next Steps

### Phase 1: Complete Domain Implementation

- [x] Implement file repository with actual file operations
- [x] Add system monitoring domain
- [x] Add infrastructure management domains
- [x] Apply Domain-Driven Design patterns
- [x] Implement Factory pattern for handlers
- [x] Implement Strategy pattern for validation
- [x] Implement Facade pattern for service orchestration
- [x] Implement Observer pattern for events

### Phase 2: Add Tests

- [ ] Unit tests for domain services
- [ ] Integration tests for application services
- [ ] End-to-end tests for CLI and HTTP

### Phase 3: Enhance Features

- [ ] Add more file operations (copy, move, search)
- [ ] Add system monitoring capabilities
- [ ] Add infrastructure management features

### Phase 4: Documentation

- [ ] API documentation
- [ ] User guides
- [ ] Development guides

## Contributing

When contributing to this project:

1. **Follow the architecture**: Place code in the appropriate layer
2. **Respect file size limits**: Split large files into smaller ones
3. **Use interfaces**: Define contracts in the domain layer
4. **Write tests**: Ensure new code is well-tested
5. **Update documentation**: Keep docs in sync with code changes

## Conclusion

The new architecture provides a solid foundation for the Mini MCP tool. It follows clean architecture principles, maintains clear separation of concerns, and provides a scalable structure for future development. The domain-driven approach ensures that business logic is well-organized and easily testable, while the infrastructure layer handles external concerns cleanly.
