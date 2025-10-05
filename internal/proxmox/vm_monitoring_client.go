package proxmox

import (
	"context"
	"fmt"

	"mini-mcp/internal/proxmox/types"
)

// GetVMPerformanceData gets performance data for a VM
func (c *client) GetVMPerformanceData(ctx context.Context, nodeName string, vmid int, config types.VMStatistics) (*types.VMPerformanceData, error) {
	endpoint := c.VMEndpoint(nodeName, fmt.Sprintf("%d", vmid), "/rrddata")

	// Add query parameters
	queryParams := make(map[string]string)
	if config.TimeFrame != "" {
		queryParams["timeframe"] = config.TimeFrame
	}
	if config.Start != "" {
		queryParams["start"] = config.Start
	}
	if config.End != "" {
		queryParams["end"] = config.End
	}

	var performanceData []types.VMPerformanceData
	err := c.GetListAndUnmarshal(ctx, endpoint, queryParams, &performanceData)
	if err != nil {
		return nil, err
	}

	if len(performanceData) == 0 {
		return nil, fmt.Errorf("no performance data available")
	}

	return &performanceData[0], nil
}

// GetVMLogs gets logs for a VM
func (c *client) GetVMLogs(ctx context.Context, nodeName string, vmid int, limit int) ([]types.VMLogEntry, error) {
	return nil, fmt.Errorf("VM logs functionality is not supported in this Proxmox version. Please use the Proxmox web interface or check system logs directly")
}

// GetVMEvents gets events for a VM
func (c *client) GetVMEvents(ctx context.Context, nodeName string, vmid int, limit int) ([]types.VMEvent, error) {
	return nil, fmt.Errorf("VM events functionality is not supported in this Proxmox version. Please use the Proxmox web interface or check system logs directly")
}
