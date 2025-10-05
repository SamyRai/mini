#!/bin/bash

# Test MCP server with proper initialization sequence
echo "Testing MCP server with Proxmox integration..."

# Initialize the session
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {"roots": {"listChanged": true}, "sampling": {}}}}' | \
docker run --rm -i -v "$(pwd)/proxmox-auth.yaml:/app/proxmox-auth.yaml:ro" mini-mcp:latest | \
grep -E '"result"|"error"' | head -1

echo ""
echo "Testing Proxmox status tool..."

# Test Proxmox status tool
echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "proxmox_status", "arguments": {}}}' | \
docker run --rm -i -v "$(pwd)/proxmox-auth.yaml:/app/proxmox-auth.yaml:ro" mini-mcp:latest | \
grep -E '"result"|"error"' | head -1
