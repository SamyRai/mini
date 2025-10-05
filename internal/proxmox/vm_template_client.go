package proxmox

import (
	"context"
	"fmt"
	"net/url"

	"mini-mcp/internal/proxmox/types"
)

// CreateVMTemplate creates a template from a VM
func (c *client) CreateVMTemplate(ctx context.Context, nodeName string, vmid int, config types.VMTemplateCreateRequest) error {
	endpoint := c.VMEndpoint(nodeName, fmt.Sprintf("%d", vmid), "/template")

	// Convert config to form data
	formData := url.Values{}
	if config.Description != "" {
		formData.Set("description", config.Description)
	}

	_, err := c.Post(ctx, endpoint, formData)
	return err
}

// DeployFromTemplate deploys a VM from a template
func (c *client) DeployFromTemplate(ctx context.Context, nodeName string, config types.VMTemplateDeployRequest) error {
	endpoint := c.VMEndpoint(nodeName, fmt.Sprintf("%d", config.TemplateID), "/clone")

	// Convert config to form data
	formData := url.Values{}
	formData.Set("newid", fmt.Sprintf("%d", config.NewID))
	if config.Name != "" {
		formData.Set("name", config.Name)
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
	if config.Full {
		formData.Set("full", "1")
	}

	_, err := c.Post(ctx, endpoint, formData)
	return err
}

// GetVMTemplates gets all VM templates
func (c *client) GetVMTemplates(ctx context.Context, nodeName string) ([]types.VMTemplate, error) {
	endpoint := c.NodeEndpoint(nodeName, "/qemu")

	var allVMs []types.VMTemplate
	err := c.GetListAndUnmarshal(ctx, endpoint, nil, &allVMs)
	if err != nil {
		return nil, err
	}

	// Filter for templates only
	var templates []types.VMTemplate
	for _, vm := range allVMs {
		if vm.Template == 1 {
			templates = append(templates, vm)
		}
	}

	return templates, nil
}
