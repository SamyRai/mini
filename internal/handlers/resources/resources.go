package handlers

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"mini-mcp/internal/types/resources"
)

// This file contains handlers for accessing resources.

// AccessResource retrieves the content of a resource by URI.
// It returns the resource data or an error if the resource is not found.
func AccessResource(uri string) (any, error) {
	switch uri {
	case "system/info":
		return getSystemInfo()
	case "docker/info":
		return getDockerInfo()
	case "docs/commands":
		return getCommandDocs()
	default:
		return nil, fmt.Errorf("unknown resource: %s", uri)
	}
}

// getSystemInfo returns basic system information.
func getSystemInfo() (*resources.SystemInfo, error) {
	// Create a new SystemInfo instance
	info := &resources.SystemInfo{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
		CPUs: runtime.NumCPU(),
	}

	// Get hostname
	if hostname, err := exec.Command("hostname").Output(); err == nil {
		info.Hostname = strings.TrimSpace(string(hostname))
	}

	// Get memory info if on Linux
	if runtime.GOOS == "linux" {
		if memInfo, err := exec.Command("free", "-m").Output(); err == nil {
			info.Memory = strings.TrimSpace(string(memInfo))
		}
	}

	return info, nil
}

// getDockerInfo returns Docker system information.
func getDockerInfo() (any, error) {
	// Check if Docker is installed
	if _, err := exec.LookPath("docker"); err != nil {
		return resources.NewDockerError("Docker not found on system"), nil
	}

	// Get Docker info
	output, err := exec.Command("docker", "info", "--format", "{{json .}}").Output()
	if err != nil {
		return resources.NewDockerError(fmt.Sprintf("Failed to get Docker info: %v", err)), nil
	}

	// Parse JSON output
	var info any
	if err := json.Unmarshal(output, &info); err != nil {
		return resources.NewDockerRawOutput(string(output)), nil
	}

	// In a real implementation, we would create a DockerInfo instance and populate its fields
	// from the parsed JSON. For now, we'll just return the raw parsed JSON.
	// Example:
	// dockerInfo := resources.NewDockerInfo()
	// dockerInfo.Version = ...
	// return dockerInfo, nil

	return info, nil
}

// getCommandDocs returns documentation for common commands.
func getCommandDocs() (resources.CommandDocCollection, error) {
	// Create a new CommandDocCollection
	docs := resources.NewCommandDocCollection()

	// Add documentation for common commands
	docs.Add("docker", resources.CommandDoc{
		Description: "Docker container management",
		Usage:       "docker [OPTIONS] COMMAND",
		URL:         "https://docs.docker.com/engine/reference/commandline/cli/",
	})

	docs.Add("git", resources.CommandDoc{
		Description: "Distributed version control system",
		Usage:       "git [--version] [--help] [-C <path>] [-c <name>=<value>] [--exec-path[=<path>]] [--html-path] [--man-path] [--info-path] [-p | --paginate | -P | --no-pager] [--no-replace-objects] [--bare] [--git-dir=<path>] [--work-tree=<path>] [--namespace=<name>] <command> [<args>]",
		URL:         "https://git-scm.com/docs",
	})

	docs.Add("kubectl", resources.CommandDoc{
		Description: "Kubernetes command line tool",
		Usage:       "kubectl [command] [TYPE] [NAME] [flags]",
		URL:         "https://kubernetes.io/docs/reference/kubectl/",
	})

	return docs, nil
}
