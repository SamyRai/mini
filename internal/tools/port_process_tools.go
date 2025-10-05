package tools

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"mini-mcp/internal/registry"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterPortProcessTools registers port/process-related tools
func RegisterPortProcessTools(server *mcp.Server, toolRegistry *registry.TypeSafeToolRegistry, executor *registry.CommandExecutor) {
	// port_process_tools - Investigate and manage ports and processes
	portProcessBuilder := registry.NewToolBuilder[PortProcessArgs](toolRegistry, "port_process_tools", "Investigate and manage network ports and processes. List ports, find processes using ports, kill processes, clean up occupied ports, and get detailed port/process information. Essential for debugging and system maintenance.")

	portProcessBuilder.
		WithHandler(func(ctx context.Context, req *mcp.CallToolRequest, args PortProcessArgs) (*mcp.CallToolResult, any, error) {
			output, err := executePortProcessCommand(executor, args.Command, args.Port, args.ProcessID, args.ProcessName, args.User, args.State)
			if err != nil {
				errorResult, _, _ := toolRegistry.CreateErrorResult(err.Error(), map[string]any{
					"command": args.Command,
					"port":    args.Port,
					"pid":     args.ProcessID,
				})
				return errorResult, nil, nil
			}

			successResult, _, _ := toolRegistry.CreateTextResult(output)
			return successResult, nil, nil
		}).
		WithValidator(func(args PortProcessArgs) error {
			if args.Command == "" {
				return registry.NewValidationError("missing_command", "command is required")
			}
			return nil
		}).
		Register()
}

// PortProcessArgs represents arguments for port/process operations
type PortProcessArgs struct {
	Command     string `json:"command" jsonschema:"Operation (list_ports, list_processes, kill_process, find_port, clean_ports, port_info, process_info, network_stats)"`
	Port        int    `json:"port,omitempty" jsonschema:"Port number to investigate"`
	ProcessID   int    `json:"process_id,omitempty" jsonschema:"Process ID to investigate"`
	ProcessName string `json:"process_name,omitempty" jsonschema:"Process name to search for"`
	User        string `json:"user,omitempty" jsonschema:"User to filter by"`
	State       string `json:"state,omitempty" jsonschema:"Port state to filter by (LISTEN, ESTABLISHED, etc.)"`
}

// executePortProcessCommand executes port/process related commands
func executePortProcessCommand(executor *registry.CommandExecutor, command string, port int, processID int, processName string, user string, state string) (string, error) {
	switch command {
	case "list_ports":
		return listPorts(executor, state)
	case "list_processes":
		return listProcesses(executor, user)
	case "kill_process":
		return killProcess(executor, processID)
	case "find_port":
		return findPort(executor, port)
	case "clean_ports":
		return cleanPorts(executor)
	case "port_info":
		return getPortInfo(executor, port)
	case "process_info":
		return getProcessInfo(executor, processID)
	case "network_stats":
		return getNetworkStats(executor)
	default:
		return "", fmt.Errorf("unsupported command: %s", command)
	}
}

// listPorts lists network ports
func listPorts(executor *registry.CommandExecutor, state string) (string, error) {
	output, err := executor.ExecuteSystemCommand(context.Background(), "netstat", "-tuln")
	if err != nil {
		return "", err
	}
	
	lines := strings.Split(output, "\n")
	var result strings.Builder
	
	for _, line := range lines {
		if state == "" || strings.Contains(line, state) {
			result.WriteString(line + "\n")
		}
	}
	
	return result.String(), nil
}

// listProcesses lists running processes
func listProcesses(executor *registry.CommandExecutor, user string) (string, error) {
	output, err := executor.ExecuteSystemCommand(context.Background(), "ps", "aux")
	if err != nil {
		return "", err
	}
	
	lines := strings.Split(output, "\n")
	var result strings.Builder
	
	for _, line := range lines {
		if user == "" || strings.Contains(line, user) {
			result.WriteString(line + "\n")
		}
	}
	
	return result.String(), nil
}

// killProcess kills a process by ID
func killProcess(executor *registry.CommandExecutor, pid int) (string, error) {
	if pid <= 0 {
		return "", fmt.Errorf("invalid process ID: %d", pid)
	}
	
	success := executor.KillProcessGracefully(context.Background(), pid)
	if !success {
		return "", fmt.Errorf("failed to kill process %d", pid)
	}
	
	return fmt.Sprintf("Process %d killed successfully", pid), nil
}

// findPort finds processes using a specific port
func findPort(executor *registry.CommandExecutor, port int) (string, error) {
	if port <= 0 {
		return "", fmt.Errorf("invalid port: %d", port)
	}
	
	output, err := executor.ExecuteSystemCommand(context.Background(), "lsof", "-i", ":"+strconv.Itoa(port))
	if err != nil {
		return "", fmt.Errorf("failed to find processes on port %d: %w", port, err)
	}
	
	return output, nil
}

// cleanPorts attempts to clean up occupied ports
func cleanPorts(executor *registry.CommandExecutor) (string, error) {
	result := "Port cleanup completed:\n"

	// Get list of used ports
	ports, err := getUsedPorts(executor)
	if err != nil {
		return "", fmt.Errorf("failed to get used ports: %w", err)
	}

	if len(ports) == 0 {
		return "No ports to clean up", nil
	}

	// Attempt to clean up each port
	cleaned := 0
	for _, port := range ports {
		if cleanupPort(executor, port) {
			result += fmt.Sprintf("✓ Cleaned up port %d\n", port)
			cleaned++
		} else {
			result += fmt.Sprintf("✗ Failed to clean up port %d\n", port)
		}
	}

	result += fmt.Sprintf("\nSummary: %d/%d ports cleaned up", cleaned, len(ports))
	return result, nil
}

// getUsedPorts gets a list of currently used ports
func getUsedPorts(executor *registry.CommandExecutor) ([]int, error) {
	// Try ss first, fallback to netstat
	output, err := executor.ExecuteSystemCommand(context.Background(), "ss", "-tuln")
	if err != nil {
		// Fallback to netstat if ss is not available
		output, err = executor.ExecuteSystemCommand(context.Background(), "netstat", "-tuln")
		if err != nil {
			return nil, fmt.Errorf("failed to get used ports: %w", err)
		}
	}

	// Parse output to extract port numbers
	ports := executor.ParsePortNumbers(output)
	return ports, nil
}

// cleanupPort attempts to clean up a specific port
func cleanupPort(executor *registry.CommandExecutor, port int) bool {
	// Get processes using the port
	processes, err := getProcessesUsingPort(executor, port)
	if err != nil || len(processes) == 0 {
		return false
	}

	// Kill processes (be careful - this is a dangerous operation)
	for _, pid := range processes {
		if executor.KillProcessGracefully(context.Background(), pid) {
			// Process killed - logged via structured logging
		}
	}

	return true
}

// getProcessesUsingPort gets PIDs of processes using a specific port
func getProcessesUsingPort(executor *registry.CommandExecutor, port int) ([]int, error) {
	// Use lsof to find processes using the port
	output, err := executor.ExecuteSystemCommand(context.Background(), "lsof", "-ti", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, fmt.Errorf("failed to find processes using port %d: %w", port, err)
	}

	// Parse PIDs from output
	pids := make([]int, 0)
	lines := strings.Split(strings.TrimSpace(output), "\n")

	for _, line := range lines {
		if pid, err := strconv.Atoi(line); err == nil {
			pids = append(pids, pid)
		}
	}

	return pids, nil
}

// getPortInfo gets detailed information about a port
func getPortInfo(executor *registry.CommandExecutor, port int) (string, error) {
	if port <= 0 {
		return "", fmt.Errorf("invalid port: %d", port)
	}
	
	output, err := executor.ExecuteSystemCommand(context.Background(), "lsof", "-i", ":"+strconv.Itoa(port))
	if err != nil {
		return "", fmt.Errorf("failed to get info for port %d: %w", port, err)
	}
	
	return output, nil
}

// getProcessInfo gets detailed information about a process
func getProcessInfo(executor *registry.CommandExecutor, pid int) (string, error) {
	if pid <= 0 {
		return "", fmt.Errorf("invalid process ID: %d", pid)
	}
	
	output, err := executor.ExecuteSystemCommand(context.Background(), "ps", "-p", strconv.Itoa(pid), "-o", "pid,ppid,cmd,user,time")
	if err != nil {
		return "", fmt.Errorf("failed to get info for process %d: %w", pid, err)
	}
	
	return output, nil
}

// getNetworkStats gets network statistics
func getNetworkStats(executor *registry.CommandExecutor) (string, error) {
	output, err := executor.ExecuteSystemCommand(context.Background(), "netstat", "-i")
	if err != nil {
		return "", err
	}
	
	return output, nil
}
