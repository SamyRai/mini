#!/bin/bash

# MCP Protocol Test Script
echo "=== MCP Server Test with Proxmox Integration ==="

# Create a named pipe for communication
PIPE="/tmp/mcp_test_pipe"
rm -f "$PIPE"
mkfifo "$PIPE"

# Start the MCP server in background
docker run --rm -i -v "$(pwd)/proxmox-auth.yaml:/app/proxmox-auth.yaml:ro" mini-mcp:latest < "$PIPE" > mcp_output.json &
SERVER_PID=$!

# Wait a moment for server to start
sleep 2

# Send initialization request
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {"roots": {"listChanged": true}, "sampling": {}}}}' > "$PIPE"

# Wait for response
sleep 1

# Send initialized notification
echo '{"jsonrpc": "2.0", "method": "notifications/initialized", "params": {}}' > "$PIPE"

# Wait for response
sleep 1

# Test Proxmox status tool
echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "proxmox_status", "arguments": {}}}' > "$PIPE"

# Wait for response
sleep 3

# Clean up
kill $SERVER_PID 2>/dev/null
rm -f "$PIPE"

# Show results
echo "=== MCP Server Output ==="
cat mcp_output.json 2>/dev/null | grep -E '"result"|"error"' | tail -5

echo ""
echo "=== Test Complete ==="
