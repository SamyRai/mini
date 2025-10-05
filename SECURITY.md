# Security Policy

## 🔒 Reporting Security Vulnerabilities

**Do NOT report security vulnerabilities through public GitHub issues.**

If you discover a security vulnerability in Mini MCP, please report it responsibly by emailing us directly at:

**security@example.com**

We will acknowledge your report within 24 hours and work with you to understand and address the issue.

## 🛡️ Security Features

Mini MCP includes multiple layers of security to protect against common vulnerabilities:

### Authentication & Authorization
- API key-based authentication
- IP whitelisting support
- Rate limiting with sliding windows
- Request correlation tracking

### Command Execution Security
- Command allowlisting (only safe commands allowed)
- Path traversal protection
- Input sanitization
- Working directory restrictions
- Command timeout enforcement

### Network Security
- No hardcoded credentials in source code
- Environment variable configuration
- Secure defaults (localhost only)

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

- **Email**: security@example.com
- **Response Time**: Within 24 hours for reports
- **Updates**: We provide regular updates on reported issues

---

*This security policy is adapted from industry best practices and may be updated periodically.*
