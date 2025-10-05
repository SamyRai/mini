# Proxmox Configuration Guide

**📚 [← Back to Main README](../README.md) | [🏗️ Architecture Documentation](README_ARCHITECTURE.md) | [🛠️ Tools Documentation](README_TOOLS.md) | [🔒 Type Safety](TYPE_SAFETY_IMPROVEMENTS.md) | [🤖 Agent Guide](AGENT.md)**

This guide explains how to configure Proxmox credentials for the mini-mcp server.

## Configuration Methods

The mini-mcp server supports multiple ways to configure Proxmox credentials:

### 1. Configuration File (Recommended)

Create a `proxmox-auth.yaml` file in the project root:

```yaml
proxmox:
  host: "your-proxmox-server.com"
  
  # API Token Authentication (Recommended)
  token_name: "user@pam!token-name"
  token_value: "your-token-value"

  # Alternative: Username/Password Authentication
  user: "user@pam"
  password: "your-password"

  # SSL Configuration
  verify_ssl: false  # Set to true in production

  # Connection timeout settings
  timeout: 30

  # Default node (will be auto-detected if not specified)
  node: "your-node-name"
```

### 2. Environment Variables

You can override any configuration using environment variables:

```bash
export PROXMOX_HOST="your-proxmox-server.com"
export PROXMOX_USER="user@pam"
export PROXMOX_PASSWORD="your-password"
export PROXMOX_TOKEN_NAME="user@pam!token-name"
export PROXMOX_TOKEN_VALUE="your-token-value"
export PROXMOX_VERIFY_SSL="true"
export PROXMOX_TIMEOUT="30"
export PROXMOX_NODE="your-node-name"
```

### 3. Home Directory Configuration

The system will also look for configuration in:
- `~/.config/mini-mcp/proxmox-auth.yaml`

## Authentication Methods

### API Token Authentication (Recommended)

1. Log into your Proxmox web interface
2. Go to **Datacenter** → **Permissions** → **API Tokens**
3. Create a new token with appropriate permissions
4. Use the token in your configuration:

```yaml
proxmox:
  host: "your-proxmox-server.com"
  token_name: "user@pam!token-name"
  token_value: "your-token-value"
```

### Username/Password Authentication

```yaml
proxmox:
  host: "your-proxmox-server.com"
  user: "user@pam"
  password: "your-password"
```

## Security Notes

- **Never commit `proxmox-auth.yaml` to version control**
- Use API tokens instead of passwords when possible
- Enable SSL verification in production (`verify_ssl: true`)
- Consider using environment variables for sensitive data
- The configuration file is automatically ignored by git

## Configuration Priority

1. Environment variables (highest priority)
2. Configuration file (`proxmox-auth.yaml`)
3. Home directory configuration (`~/.config/mini-mcp/proxmox-auth.yaml`)

## Example Usage

Once configured, you can use the Proxmox tools:

```bash
# Get Proxmox status
mini-mcp-cli proxmox_status

# Or through MCP protocol
# The tool will automatically use the configured credentials
```

## Troubleshooting

### Common Issues

1. **"Proxmox host is required"**
   - Ensure `PROXMOX_HOST` is set or configured in `proxmox-auth.yaml`

2. **"Authentication failed"**
   - Check your credentials
   - Verify the user has appropriate permissions
   - For token auth, ensure both `token_name` and `token_value` are set

3. **"SSL verification failed"**
   - Set `verify_ssl: false` for development
   - Ensure SSL certificates are valid for production

4. **"Connection timeout"**
   - Check network connectivity to Proxmox server
   - Increase timeout value if needed
   - Verify firewall settings

### Testing Configuration

You can test your configuration by running:

```bash
go run ./cmd/mini-mcp
```

Then use the `proxmox_status` tool to verify connectivity.
