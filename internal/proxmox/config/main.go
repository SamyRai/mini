package config

import (
	"os"

	"mini-mcp/internal/proxmox/types"

	"gopkg.in/yaml.v3"
)

// LoadMainConfig loads the main configuration from a YAML file
func LoadMainConfig(filename string) (*types.MainConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Expand environment variables in the YAML content
	expandedData := os.ExpandEnv(string(data))

	var config types.MainConfig
	if err := yaml.Unmarshal([]byte(expandedData), &config); err != nil {
		return nil, err
	}

	// Set defaults
	if config.Proxmox.Storage.Local == "" {
		config.Proxmox.Storage.Local = "local"
	}
	if config.Proxmox.Storage.ISOStorage == "" {
		config.Proxmox.Storage.ISOStorage = "local"
	}
	if config.Proxmox.Storage.SnippetsStore == "" {
		config.Proxmox.Storage.SnippetsStore = "local"
	}
	if config.Proxmox.Network.Bridge == "" {
		config.Proxmox.Network.Bridge = "vmbr0"
	}
	if config.Proxmox.Network.Model == "" {
		config.Proxmox.Network.Model = "virtio"
	}
	if config.Proxmox.VMDefaults.ScsiController == "" {
		config.Proxmox.VMDefaults.ScsiController = "virtio-scsi-pci"
	}
	if config.Proxmox.VMDefaults.VGA == "" {
		config.Proxmox.VMDefaults.VGA = "std"
	}
	if config.Proxmox.VMDefaults.Agent == "" {
		config.Proxmox.VMDefaults.Agent = "enabled=1"
	}
	if config.Proxmox.VMDefaults.BootOrder == "" {
		config.Proxmox.VMDefaults.BootOrder = "order=scsi0;ide2"
	}
	if config.Proxmox.VMDefaults.BootDisk == "" {
		config.Proxmox.VMDefaults.BootDisk = "scsi0"
	}

	return &config, nil
}

// GetVMByID returns a VM configuration by ID
func GetVMByID(c *types.MainConfig, id int) *types.VMConfigItem {
	for _, vm := range c.VMs {
		if vm.ID == id {
			return &vm
		}
	}
	return nil
}

// GetFirstVM returns the first VM configuration
func GetFirstVM(c *types.MainConfig) *types.VMConfigItem {
	if len(c.VMs) > 0 {
		return &c.VMs[0]
	}
	return nil
}
