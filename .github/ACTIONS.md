# GitHub Actions Workflows

This document describes the comprehensive GitHub Actions setup for the Mini MCP project, following modern production-ready practices for open source projects.

## Overview

The project includes multiple workflows designed to ensure code quality, security, and reliability:

- **CI/CD Pipeline**: Automated testing, building, and deployment
- **Security Scanning**: Comprehensive security analysis and vulnerability detection
- **Code Quality**: Linting, formatting, and performance checks
- **Dependency Management**: Automated dependency updates and security scanning
- **Maintenance**: Regular maintenance tasks and cleanup

## Workflows

### 1. CI Workflow (`.github/workflows/ci.yml`)

**Purpose**: Continuous Integration for testing, building, and basic security checks.

**Triggers**:
- Push to `main` and `develop` branches
- Pull requests to `main` and `develop` branches
- Manual dispatch with options

**Features**:
- Multi-platform testing (Ubuntu, Windows, macOS)
- Multi-version Go testing (1.24, 1.25)
- Race condition testing
- Coverage reporting
- Security scanning with Gosec and Trivy
- Static analysis with staticcheck
- Build artifacts for multiple platforms

**Manual Options**:
- `skip_tests`: Skip test execution

### 2. Security Workflow (`.github/workflows/security.yml`)

**Purpose**: Comprehensive security scanning and vulnerability detection.

**Triggers**:
- Push to `main` and `develop` branches
- Pull requests to `main` and `develop` branches
- Weekly scheduled runs (Monday 3 AM UTC)
- Manual dispatch with scan type selection

**Features**:
- Trivy vulnerability scanning
- Gosec security analysis
- Semgrep static analysis
- OWASP Dependency Check
- License compliance checking
- Container security scanning
- govulncheck for Go vulnerabilities

**Manual Options**:
- `scan_type`: Choose between 'all', 'dependencies', 'code', 'container'

### 3. CodeQL Workflow (`.github/workflows/codeql.yml`)

**Purpose**: Advanced static analysis using GitHub's CodeQL engine.

**Triggers**:
- Push to `main` and `develop` branches
- Pull requests to `main` and `develop` branches
- Weekly scheduled runs (Monday 2 AM UTC)
- Manual dispatch

**Features**:
- Security and quality queries
- Extended security analysis
- SARIF result upload to GitHub Security tab

### 4. Release Workflow (`.github/workflows/release.yml`)

**Purpose**: Automated release management and Docker image building.

**Triggers**:
- Git tags starting with 'v*'
- Manual dispatch with options

**Features**:
- GoReleaser integration
- Multi-platform builds
- Docker image building and scanning
- Container vulnerability scanning
- Artifact upload to GitHub Releases

**Manual Options**:
- `tag`: Tag to release
- `skip_tests`: Skip test execution
- `skip_docker`: Skip Docker build

### 5. Security Audit Workflow (`.github/workflows/security-audit.yml`)

**Purpose**: Comprehensive security auditing with multiple tools.

**Triggers**:
- Push to `main` and `develop` branches
- Pull requests to `main` and `develop` branches
- Weekly scheduled runs (Monday 4 AM UTC)
- Manual dispatch

**Features**:
- Trivy vulnerability scanning
- Gosec security analysis
- Semgrep static analysis
- OWASP Dependency Check
- Supply chain security analysis
- License compliance checking

### 6. Quality Workflow (`.github/workflows/quality.yml`)

**Purpose**: Code quality analysis and performance testing.

**Triggers**:
- Push to `main` and `develop` branches
- Pull requests to `main` and `develop` branches
- Manual dispatch

**Features**:
- golangci-lint analysis
- Static analysis tools
- Performance benchmarking
- Code complexity analysis
- Coverage reporting

### 7. Maintenance Workflow (`.github/workflows/maintenance.yml`)

**Purpose**: Regular maintenance tasks and dependency management.

**Triggers**:
- Weekly scheduled runs (Monday 6 AM UTC)
- Manual dispatch

**Features**:
- Dependency updates checking
- Security vulnerability scanning
- Code cleanup
- Dead code detection

### 8. Dependabot Auto-merge (`.github/workflows/dependabot-auto-merge.yml`)

**Purpose**: Automated merging of Dependabot pull requests after validation.

**Triggers**:
- Dependabot pull requests

**Features**:
- Automated testing
- Security vulnerability checking
- Auto-merge after successful validation

## Security Features

### Permissions
All workflows follow the principle of least privilege:
- `contents: read` - Read repository contents
- `pull-requests: read` - Read pull request information
- `security-events: write` - Upload security scan results
- `packages: write` - Upload packages to GitHub Packages
- `id-token: write` - Required for OIDC token generation

### Security Scanning Tools
- **Trivy**: Vulnerability scanning for dependencies and containers
- **Gosec**: Go security analysis
- **Semgrep**: Static analysis with security rules
- **CodeQL**: GitHub's advanced static analysis
- **govulncheck**: Go vulnerability database checking
- **OWASP Dependency Check**: Comprehensive dependency analysis

### Caching Strategy
- Go module caching with `cache-dependency-path: '**/go.sum'`
- Docker layer caching with GitHub Actions cache
- Build artifact caching for improved performance

## Dependabot Configuration

The project uses Dependabot for automated dependency updates:

- **Go modules**: Weekly updates on Mondays
- **GitHub Actions**: Weekly updates on Mondays
- **Docker**: Weekly updates on Mondays
- **Grouping**: Dependencies are grouped by type (minor/patch vs major)
- **Auto-merge**: Enabled for validated updates

## Best Practices Implemented

### 1. Security
- All workflows use the latest action versions
- Proper permission scoping
- Security scanning on every PR
- SARIF result upload to GitHub Security tab
- Dependency vulnerability checking

### 2. Performance
- Optimized caching strategies
- Parallel job execution
- Artifact retention policies
- Build matrix optimization

### 3. Reliability
- Comprehensive testing across platforms
- Race condition detection
- Timeout configurations
- Error handling and reporting

### 4. Maintainability
- Clear workflow documentation
- Modular workflow design
- Consistent naming conventions
- Regular maintenance tasks

## Manual Workflow Execution

All workflows support manual execution with the following options:

### CI Workflow
```bash
# Skip tests
gh workflow run ci.yml -f skip_tests=true
```

### Security Workflow
```bash
# Run specific scan type
gh workflow run security.yml -f scan_type=dependencies
```

### Release Workflow
```bash
# Release with custom tag
gh workflow run release.yml -f tag=v1.0.0 -f skip_tests=false
```

### Quality Workflow
```bash
# Run specific quality check
gh workflow run quality.yml -f quality_type=performance
```

### Maintenance Workflow
```bash
# Run specific maintenance task
gh workflow run maintenance.yml -f maintenance_type=security
```

## Monitoring and Alerts

The workflows are designed to provide comprehensive monitoring:

1. **GitHub Security Tab**: All security scan results are uploaded
2. **Codecov**: Coverage reports and trends
3. **GitHub Actions**: Workflow run history and logs
4. **Dependabot**: Automated dependency update tracking

## Troubleshooting

### Common Issues

1. **Workflow Failures**: Check the Actions tab for detailed logs
2. **Security Alerts**: Review the Security tab for vulnerability details
3. **Dependency Issues**: Check Dependabot pull requests
4. **Build Failures**: Review build logs for platform-specific issues

### Debugging

1. **Enable Debug Logging**: Set `ACTIONS_STEP_DEBUG=true` in repository secrets
2. **Check Permissions**: Ensure workflows have required permissions
3. **Review Caching**: Clear caches if builds are inconsistent
4. **Update Actions**: Keep action versions up to date

## Contributing

When contributing to the workflows:

1. Follow the existing patterns and conventions
2. Test changes in a fork before submitting
3. Update documentation for any new features
4. Ensure security best practices are maintained
5. Consider the impact on workflow performance

## Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Go Security Best Practices](https://golang.org/doc/security)
- [OWASP Security Guidelines](https://owasp.org/)
- [Dependabot Documentation](https://docs.github.com/en/code-security/dependabot)
