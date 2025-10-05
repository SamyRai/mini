package config

import (
	"os"

	"gopkg.in/yaml.v3"
	"mini-mcp/internal/proxmox/types"
)

// LoadAuthConfig loads authentication configuration from a YAML file
func LoadAuthConfig(filename string) (*types.AuthConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config types.AuthConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Set defaults
	if config.Proxmox.Timeout == 0 {
		config.Proxmox.Timeout = 30
	}
	if !config.Proxmox.VerifySSL {
		config.Proxmox.VerifySSL = false // Default to false for development
	}

	return &config, nil
}

// GetHost returns the Proxmox host
func GetHost(c *types.AuthConfig) string {
	return c.Proxmox.Host
}

// GetUser returns the Proxmox user
func GetUser(c *types.AuthConfig) string {
	return c.Proxmox.User
}

// GetPassword returns the Proxmox password
func GetPassword(c *types.AuthConfig) string {
	return c.Proxmox.Password
}

// GetTokenName returns the Proxmox token name
func GetTokenName(c *types.AuthConfig) string {
	return c.Proxmox.TokenName
}

// GetTokenValue returns the Proxmox token value
func GetTokenValue(c *types.AuthConfig) string {
	return c.Proxmox.TokenValue
}

// GetVerifySSL returns whether to verify SSL certificates
func GetVerifySSL(c *types.AuthConfig) bool {
	return c.Proxmox.VerifySSL
}

// GetTimeout returns the connection timeout in seconds
func GetTimeout(c *types.AuthConfig) int {
	return c.Proxmox.Timeout
}

// GetNode returns the preferred node (if specified)
func GetNode(c *types.AuthConfig) string {
	return c.Proxmox.Node
}

// IsTokenAuth returns true if token authentication should be used
func IsTokenAuth(c *types.AuthConfig) bool {
	return c.Proxmox.TokenName != "" && c.Proxmox.TokenValue != ""
}
