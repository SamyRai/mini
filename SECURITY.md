# Security Policy

## 🔒 Reporting Security Vulnerabilities

**Do NOT report security vulnerabilities through public GitHub issues.**

If you discover a security vulnerability in Mini MCP, please report it responsibly by emailing us directly at:

**security@damirmukimov.com**

We will acknowledge your report within 24 hours and work with you to understand and address the issue.

## 🛡️ Security Features

Mini MCP includes multiple layers of security to protect against common vulnerabilities:

### Authentication & Authorization
- API key-based authentication with rotation support
- IP whitelisting support for restricted access
- Rate limiting with sliding windows and burst handling
- Request correlation tracking for audit trails
- Enterprise-grade authentication mechanisms

### Command Execution Security
- Command allowlisting (only pre-approved commands allowed)
- Path traversal protection against directory traversal attacks
- Advanced input sanitization preventing injection attacks
- Working directory restrictions for sandboxed execution
- Command timeout enforcement preventing hanging processes
- Output size limits preventing memory exhaustion

### File Operations Security
- Path validation and sanitization
- Directory traversal attack prevention
- System directory access restrictions
- Permission handling and metadata preservation
- Automatic directory creation for safe operations

### Network Security
- No hardcoded credentials in source code
- Environment variable configuration
- Secure defaults (localhost only)
- SSH key authentication for remote operations
- Connection timeout controls

## 🔧 Supported Versions

| Version | Supported          |
|---------|-------------------|
| Latest  | ✅ Full support   |
| All others | ❌ Not supported |

## 🚨 Vulnerability Disclosure Process

1. **Report**: Email security@example.com with vulnerability details
2. **Acknowledge**: We respond within 24 hours
3. **Assess**: We evaluate the vulnerability and impact
4. **Fix**: We develop and test a fix
5. **Release**: We release a security update
6. **Announce**: We publish a security advisory (after reasonable time)

## 📋 What to Include in Reports

When reporting a vulnerability, please include:

- **Description**: Clear description of the vulnerability
- **Impact**: Potential impact and attack scenarios
- **Reproduction**: Steps to reproduce the issue
- **Environment**: Version, OS, configuration details
- **Suggestions**: Any suggested fixes (optional)

## 🎯 Scope

This security policy covers:

- Mini MCP server and CLI
- Official Docker images
- Supported integrations (Proxmox, Docker, etc.)
- Documentation and examples

## ⚖️ Legal Safe Harbor

We consider security research conducted in accordance with this policy to be:

- Authorized in accordance with applicable laws
- Exempt from any restrictions in our Terms of Service
- Eligible for any bug bounty programs (if applicable)

## 📞 Contact

For security-related questions or concerns:

- **Email**: security@damirmukimov.com
- **Response Time**: Within 24 hours for reports
- **Updates**: We provide regular updates on reported issues

## 🔍 Security Best Practices

When using Mini MCP in production:

1. **Environment Configuration**
   - Use environment variables for sensitive configuration
   - Enable SSL verification in production
   - Set appropriate timeout values
   - Configure rate limiting based on your needs

2. **Access Control**
   - Use API keys for authentication
   - Implement IP whitelisting for restricted access
   - Monitor access patterns and audit logs
   - Rotate API keys regularly

3. **Command Security**
   - Review the command allowlist regularly
   - Use specific paths for file operations
   - Monitor command execution logs
   - Implement additional validation for custom use cases

4. **Network Security**
   - Use SSH keys for remote operations
   - Implement proper firewall rules
   - Monitor network connections
   - Use secure protocols (HTTPS, SSH)

---

*This security policy is adapted from industry best practices and may be updated periodically.*
