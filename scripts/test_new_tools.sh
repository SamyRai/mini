#!/bin/bash

# Test script for new MCP tools
echo "=== Testing New MCP Tools ==="

# Create a named pipe for communication
PIPE="/tmp/mcp_tools_test_pipe"
rm -f "$PIPE"
mkfifo "$PIPE"

# Start the MCP server in background
./mini-mcp < "$PIPE" > mcp_tools_output.json &
SERVER_PID=$!

# Wait a moment for server to start
sleep 2

echo "Sending initialization..."
# Send initialization request
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {"roots": {"listChanged": true}, "sampling": {}}}}' > "$PIPE"

# Wait for response
sleep 1

echo "Sending initialized notification..."
# Send initialized notification
echo '{"jsonrpc": "2.0", "method": "notifications/initialized", "params": {}}' > "$PIPE"

# Wait for response
sleep 1

echo "Testing code_analysis tool..."
# Test code analysis tool
echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "code_analysis", "arguments": {"path": ".", "type": "structure", "recursive": true}}}' > "$PIPE"

# Wait for response
sleep 3

echo "Testing project_explorer tool..."
# Test project explorer tool
echo '{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "project_explorer", "arguments": {"path": ".", "command": "overview"}}}' > "$PIPE"

# Wait for response
sleep 3

echo "Testing documentation_generator tool..."
# Test documentation generator tool
echo '{"jsonrpc": "2.0", "id": 4, "method": "tools/call", "params": {"name": "documentation_generator", "arguments": {"path": ".", "command": "summary"}}}' > "$PIPE"

# Wait for response
sleep 3

echo "Testing dev_workflow tool..."
# Test dev workflow tool
echo '{"jsonrpc": "2.0", "id": 5, "method": "tools/call", "params": {"name": "dev_workflow", "arguments": {"path": ".", "command": "test", "args": ["-v", "./..."]}}}' > "$PIPE"

# Wait for response
sleep 5

# Clean up
echo "Cleaning up..."
kill $SERVER_PID 2>/dev/null
rm -f "$PIPE"

# Show results
echo "=== MCP Tools Test Results ==="
if [ -f "mcp_tools_output.json" ]; then
    echo "Output file created. Checking for results..."
    cat mcp_tools_output.json | jq '.result' 2>/dev/null || cat mcp_tools_output.json | grep -E '"result"|"error"' | tail -10
else
    echo "No output file found"
fi

echo ""
echo "=== Test Complete ==="
