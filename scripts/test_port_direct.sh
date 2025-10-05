#!/bin/bash

echo "=== Testing Port Functionality via MCP Protocol ==="

# Test port_process_tools with list_ports command
echo "1. Testing port_process_tools with list_ports command..."
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/call", "params": {"name": "port_process_tools", "arguments": {"command": "list_ports", "protocol": "all", "output_format": "json", "verbose": true, "timeout": 30}}}' | timeout 10s ./bin/mini-mcp

echo -e "\n2. Testing port_process_tools with find_port command for port 22..."
echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "port_process_tools", "arguments": {"command": "find_port", "port": 22, "output_format": "json", "verbose": true, "timeout": 30}}}' | timeout 10s ./bin/mini-mcp

echo -e "\n3. Testing port_process_tools with network_stats command..."
echo '{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "port_process_tools", "arguments": {"command": "network_stats", "output_format": "json", "verbose": true, "timeout": 30}}}' | timeout 10s ./bin/mini-mcp

echo -e "\n=== Port Functionality Test Complete ==="
