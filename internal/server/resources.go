package server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerResources registers static resources and templates.
func registerResources(server *mcp.Server) {
	// Add basic system information resource
	server.AddResource(&mcp.Resource{
		Name:        "system_info",
		Description: "Basic system information",
		MIMEType:    "application/json",
		URI:         "system://info",
	}, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		info := map[string]interface{}{
			"hostname": getHostname(),
			"working_directory": getWorkingDirectory(),
			"environment": getEnvironment(),
		}
		
		jsonData, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal system info: %w", err)
		}
		
		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     string(jsonData),
			}},
		}, nil
	})
}

// getHostname returns the system hostname
func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

// getWorkingDirectory returns the current working directory
func getWorkingDirectory() string {
	wd, err := os.Getwd()
	if err != nil {
		return "unknown"
	}
	return wd
}

// getEnvironment returns basic environment information
func getEnvironment() map[string]string {
	env := make(map[string]string)
	for _, e := range os.Environ() {
		parts := splitEnvVar(e)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}
	return env
}

// splitEnvVar splits an environment variable into key and value
func splitEnvVar(envVar string) []string {
	parts := make([]string, 2)
	for i, char := range envVar {
		if char == '=' {
			parts[0] = envVar[:i]
			parts[1] = envVar[i+1:]
			break
		}
	}
	return parts
}