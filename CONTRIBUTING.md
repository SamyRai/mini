# Contributing to Mini MCP

Thank you for your interest in contributing to Mini MCP! We welcome contributions from the community and appreciate your help in making this project better.

## 🚀 Quick Start

1. **Fork** the repository on GitHub
2. **Clone** your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/mini-mcp.git
   cd mini-mcp
   ```

3. **Set up** the development environment:
   ```bash
   make dev-setup
   ```

4. **Create a branch** for your feature:
   ```bash
   git checkout -b feature/amazing-feature
   ```

5. **Make your changes** and ensure tests pass:
   ```bash
   go test -race -cover ./...
   ```

6. **Commit** your changes with a clear message:
   ```bash
   git commit -m "feat: add amazing feature

   - Implement new functionality
   - Add comprehensive tests
   - Update documentation"
   ```

7. **Push** to your fork:
   ```bash
   git push origin feature/amazing-feature
   ```

8. **Open a Pull Request** on GitHub with a detailed description

## 🧪 Development Guidelines

### Code Standards

- Follow Go best practices and conventions
- Maintain test coverage above 80%
- Use `go fmt` to format your code
- Run `go vet` to check for potential issues
- Ensure all tests pass with `go test -race ./...`
- Follow the project's architecture patterns
- Use Go 1.25 generics where appropriate
- Avoid reflection in favor of type-safe approaches

### Testing

- Write tests for new functionality
- Use table-driven tests where appropriate
- Test edge cases and error conditions
- Run tests with race detection: `go test -race ./...`
- Test security features thoroughly
- Include integration tests for new tools
- Maintain test coverage reports

### Documentation

- Update documentation for new features
- Keep examples current and working
- Document breaking changes clearly
- Follow the existing documentation structure
- Include usage examples for new tools
- Update architecture docs for structural changes

## 🔒 Security

If you discover a security vulnerability, please follow our security policy:

1. **Do NOT** open a public issue
2. **Email** security reports to: security@damirmukimov.com
3. **Include** detailed steps to reproduce the issue
4. **Allow** time for the fix before public disclosure

### Security Considerations for Contributions

- All security-related changes require review
- Test security features thoroughly
- Follow the principle of least privilege
- Validate all inputs and sanitize outputs
- Use the project's security validation patterns

## 📝 Code of Conduct

We are committed to providing a welcoming and inclusive environment for all contributors.

### Our Standards

- Use welcoming and inclusive language
- Be respectful of differing viewpoints and experiences
- Give and gracefully accept constructive feedback
- Focus on what is best for the community

### Unacceptable Behavior

- Harassment, intimidation, or discrimination
- Offensive comments or personal attacks
- Publishing others' private information
- Other conduct that could reasonably be considered inappropriate

## 📄 License

By contributing, you agree that your contributions will be licensed under the same MIT License that covers the project.

## 🏗️ Architecture Guidelines

When contributing, please follow the project's clean architecture principles:

- **Domain Layer**: Business logic and domain models
- **Application Layer**: Use cases and orchestration
- **Infrastructure Layer**: External concerns (CLI, HTTP, persistence)
- **Shared Layer**: Cross-cutting concerns (auth, config, logging, security)

### File Organization

- Keep files under 250 lines when possible
- Split large files into smaller, focused ones
- Use interfaces for dependencies
- Follow the Single Responsibility Principle

## 🧪 Testing Requirements

- Write tests for all new functionality
- Maintain >80% test coverage
- Use table-driven tests for multiple scenarios
- Test security features thoroughly
- Include integration tests for new tools

## 📚 Documentation Requirements

- Update relevant documentation files
- Include usage examples for new features
- Update architecture docs for structural changes
- Keep examples current and working

## 🙏 Acknowledgments

Thank you to all our contributors! Your work helps make Mini MCP better for everyone.
