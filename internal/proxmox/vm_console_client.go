package proxmox

import (
	"context"
	"fmt"

	"mini-mcp/internal/proxmox/types"
)

// GetVMConsole gets console access information for a VM
func (c *client) GetVMConsole(ctx context.Context, nodeName string, vmid int, config types.VMConsoleRequest) (*types.VMConsole, error) {
	// Determine console type if not specified
	consoleType := config.ConsoleType
	if consoleType == "" {
		// Auto-detect based on VM configuration
		vmConfig, err := c.GetVMConfig(ctx, nodeName, vmid)
		if err != nil {
			return nil, fmt.Errorf("failed to get VM config for console type detection: %v", err)
		}

		// Check for serial console first (better for web terminal), then SPICE, then VNC
		// For modern VMs with QXL VGA, prefer serial console for better web terminal compatibility
		if vmConfig.VGA == "qxl" {
			consoleType = "serial" // Prefer serial for text-based web terminal
		} else {
			// Default to VNC
			consoleType = "vnc"
		}
	}

	var endpoint string
	switch consoleType {
	case "spice":
		endpoint = c.VMEndpoint(nodeName, fmt.Sprintf("%d", vmid), "/spiceproxy")
	case "serial":
		endpoint = c.VMEndpoint(nodeName, fmt.Sprintf("%d", vmid), "/termproxy")
	default: // "vnc"
		endpoint = c.VMEndpoint(nodeName, fmt.Sprintf("%d", vmid), "/vncproxy")
	}

	var console types.VMConsole
	err := c.PostAndUnmarshal(ctx, endpoint, nil, &console)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s console: %v", consoleType, err)
	}

	return &console, nil
}
