#!/bin/bash

echo "=== Testing Mini MCP Port Functionality via MCP Protocol ==="

# Test 1: List ports
echo "1. Testing list_ports command..."
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/call", "params": {"name": "port_process_tools", "arguments": {"command": "list_ports", "protocol": "all", "output_format": "json", "verbose": true, "timeout": 30}}}' | ./bin/mini-mcp

echo -e "\n2. Testing find_port command for port 22..."
echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "port_process_tools", "arguments": {"command": "find_port", "port": 22, "output_format": "json", "verbose": true, "timeout": 30}}}' | ./bin/mini-mcp

echo -e "\n3. Testing network_stats command..."
echo '{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "port_process_tools", "arguments": {"command": "network_stats", "output_format": "json", "verbose": true, "timeout": 30}}}' | ./bin/mini-mcp

echo -e "\n=== Port Functionality Test Complete ==="
