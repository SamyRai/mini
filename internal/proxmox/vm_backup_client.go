package proxmox

import (
	"context"
	"fmt"
	"net/url"

	"mini-mcp/internal/proxmox/types"
)

// CreateVMBackup creates a backup of a VM
func (c *client) CreateVMBackup(ctx context.Context, nodeName string, vmid int, config types.VMBackupCreateRequest) error {
	endpoint := c.VMDumpEndpoint(nodeName)

	// Convert config to form data
	formData := url.Values{}
	formData.Set("vmid", fmt.Sprintf("%d", vmid))
	formData.Set("storage", config.Storage)
	if config.Mode != "" {
		formData.Set("mode", config.Mode)
	}
	if config.Compress != "" {
		formData.Set("compress", config.Compress)
	}
	if config.Remove {
		formData.Set("remove", "1")
	}
	if config.Notes != "" {
		formData.Set("notes", config.Notes)
	}
	if config.MailTo != "" {
		formData.Set("mailto", config.MailTo)
	}
	// Note: mailpolicy is not supported by the API

	_, err := c.Post(ctx, endpoint, formData)
	return err
}

// GetVMBackups gets all backups for a storage
func (c *client) GetVMBackups(ctx context.Context, nodeName string, storageName string) ([]types.VMBackup, error) {
	endpoint := c.StorageContentEndpoint(nodeName, storageName)

	var allContent []types.VMBackup
	err := c.GetListAndUnmarshal(ctx, endpoint, nil, &allContent)
	if err != nil {
		return nil, err
	}

	// Filter for backup files only
	var backups []types.VMBackup
	for _, item := range allContent {
		if item.Content == "backup" {
			backups = append(backups, item)
		}
	}

	return backups, nil
}

// RestoreVMBackup restores a VM from backup
func (c *client) RestoreVMBackup(ctx context.Context, nodeName string, vmid int, config types.VMBackupRestoreRequest) error {
	endpoint := c.NodeEndpoint(nodeName, "/qemu")

	// Convert config to form data
	formData := url.Values{}
	formData.Set("vmid", fmt.Sprintf("%d", vmid))
	formData.Set("storage", config.Storage)
	formData.Set("archive", config.Backup)
	if config.Force {
		formData.Set("force", "1")
	}
	if config.Pool != "" {
		formData.Set("pool", config.Pool)
	}

	_, err := c.Post(ctx, endpoint, formData)
	return err
}
