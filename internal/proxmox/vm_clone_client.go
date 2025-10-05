package proxmox

import (
	"context"
	"fmt"
	"net/url"

	"mini-mcp/internal/proxmox/types"
)

// CloneVM clones a VM
func (c *client) CloneVM(ctx context.Context, nodeName string, vmid int, config types.VMCloneRequest) error {
	endpoint := c.VMEndpoint(nodeName, fmt.Sprintf("%d", vmid), "/clone")

	// Convert config to form data
	formData := url.Values{}
	formData.Set("newid", fmt.Sprintf("%d", config.NewID))
	if config.Name != "" {
		formData.Set("name", config.Name)
	}
	if config.Description != "" {
		formData.Set("description", config.Description)
	}
	if config.Full {
		formData.Set("full", "1")
	}
	if config.Pool != "" {
		formData.Set("pool", config.Pool)
	}
	if config.Storage != "" {
		formData.Set("storage", config.Storage)
	}
	if config.Target != "" {
		formData.Set("target", config.Target)
	}

	_, err := c.Post(ctx, endpoint, formData)
	return err
}

// MigrateVM migrates a VM to another node
func (c *client) MigrateVM(ctx context.Context, nodeName string, vmid int, config types.VMMigrateRequest) error {
	endpoint := c.VMEndpoint(nodeName, fmt.Sprintf("%d", vmid), "/migrate")

	// Convert config to form data
	formData := url.Values{}
	formData.Set("target", config.Target)
	if config.Online {
		formData.Set("online", "1")
	}
	if config.WithLocal {
		formData.Set("with-local", "1")
	}
	if config.Force {
		formData.Set("force", "1")
	}

	_, err := c.Post(ctx, endpoint, formData)
	return err
}
