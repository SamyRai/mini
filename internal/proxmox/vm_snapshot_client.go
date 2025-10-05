package proxmox

import (
	"context"
	"fmt"
	"net/url"

	"mini-mcp/internal/proxmox/types"
)

// CreateVMSnapshot creates a snapshot of a VM
func (c *client) CreateVMSnapshot(ctx context.Context, nodeName string, vmid int, config types.VMSnapshotCreateRequest) error {
	endpoint := c.VMEndpoint(nodeName, fmt.Sprintf("%d", vmid), "/snapshot")

	// Convert config to form data
	formData := url.Values{}
	formData.Set("snapname", config.SnapName)
	if config.Description != "" {
		formData.Set("description", config.Description)
	}
	if config.VMState {
		formData.Set("vmstate", "1")
	}

	_, err := c.Post(ctx, endpoint, formData)
	return err
}

// DeleteVMSnapshot deletes a VM snapshot
func (c *client) DeleteVMSnapshot(ctx context.Context, nodeName string, vmid int, snapName string) error {
	endpoint := c.VMEndpoint(nodeName, fmt.Sprintf("%d", vmid), fmt.Sprintf("/snapshot/%s", snapName))

	_, err := c.Delete(ctx, endpoint)
	return err
}

// RollbackVMSnapshot rolls back a VM to a snapshot
func (c *client) RollbackVMSnapshot(ctx context.Context, nodeName string, vmid int, config types.VMSnapshotRollbackRequest) error {
	endpoint := c.VMEndpoint(nodeName, fmt.Sprintf("%d", vmid), fmt.Sprintf("/snapshot/%s/rollback", config.SnapName))

	// Convert config to form data
	formData := url.Values{}
	if config.Start {
		formData.Set("start", "1")
	}

	_, err := c.Post(ctx, endpoint, formData)
	return err
}

// GetVMSnapshots gets all snapshots for a VM
func (c *client) GetVMSnapshots(ctx context.Context, nodeName string, vmid int) ([]types.VMSnapshot, error) {
	endpoint := c.VMEndpoint(nodeName, fmt.Sprintf("%d", vmid), "/snapshot")

	var snapshots []types.VMSnapshot
	err := c.GetListAndUnmarshal(ctx, endpoint, nil, &snapshots)
	return snapshots, err
}
