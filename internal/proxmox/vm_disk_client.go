package proxmox

import (
	"context"
	"fmt"
	"net/url"

	"mini-mcp/internal/proxmox/types"
)

// AddVMDisk adds a disk to a VM
func (c *client) AddVMDisk(ctx context.Context, nodeName string, vmid int, config types.VMDiskAddRequest) error {
	endpoint := c.VMEndpoint(nodeName, fmt.Sprintf("%d", vmid), "/config")

	// Convert config to form data - use the disk ID as the key with storage:size format
	formData := url.Values{}
	diskConfig := fmt.Sprintf("%s:%s", config.Storage, config.Size)
	if config.Format != "" {
		diskConfig += fmt.Sprintf(",format=%s", config.Format)
	}
	if config.Cache != "" {
		diskConfig += fmt.Sprintf(",cache=%s", config.Cache)
	}
	if config.SSD {
		diskConfig += ",ssd=1"
	}

	// Set the disk configuration using the disk ID as key
	formData.Set(config.ID, diskConfig)

	_, err := c.Post(ctx, endpoint, formData)
	return err
}

// RemoveVMDisk removes a disk from a VM
func (c *client) RemoveVMDisk(ctx context.Context, nodeName string, vmid int, diskID string) error {
	endpoint := c.VMEndpoint(nodeName, fmt.Sprintf("%d", vmid), "/config")

	// Convert config to form data - use delete parameter
	formData := url.Values{}
	formData.Set("delete", diskID)

	_, err := c.Put(ctx, endpoint, formData)
	return err
}

// ResizeVMDisk resizes a VM disk
func (c *client) ResizeVMDisk(ctx context.Context, nodeName string, vmid int, config types.VMDiskResizeRequest) error {
	endpoint := c.VMEndpoint(nodeName, fmt.Sprintf("%d", vmid), "/resize")

	// Convert config to form data
	formData := url.Values{}
	formData.Set("disk", config.ID)
	formData.Set("size", config.Size)

	_, err := c.Put(ctx, endpoint, formData)
	return err
}

// MoveVMDisk moves a VM disk to another storage
func (c *client) MoveVMDisk(ctx context.Context, nodeName string, vmid int, config types.VMDiskMoveRequest) error {
	endpoint := c.VMEndpoint(nodeName, fmt.Sprintf("%d", vmid), "/move_disk")

	// Convert config to form data
	formData := url.Values{}
	formData.Set("disk", config.ID)
	formData.Set("storage", config.Storage)
	if config.Delete {
		formData.Set("delete", "1")
	}
	if config.Online {
		formData.Set("online", "1")
	}

	_, err := c.Post(ctx, endpoint, formData)
	return err
}
