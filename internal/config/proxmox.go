package config

import (
	"fmt"
	"os"
	"path/filepath"

	"mini-mcp/internal/proxmox/types"
	"gopkg.in/yaml.v3"
)

// ProxmoxConfig holds Proxmox configuration
type ProxmoxConfig struct {
	Host      string `yaml:"host" json:"host"`
	User      string `yaml:"user" json:"user"`
	Password  string `yaml:"password" json:"password"`
	TokenName string `yaml:"token_name" json:"token_name"`
	TokenValue string `yaml:"token_value" json:"token_value"`
	VerifySSL bool   `yaml:"verify_ssl" json:"verify_ssl"`
	Timeout   int    `yaml:"timeout" json:"timeout"`
	Node      string `yaml:"node" json:"node"`
}

// LoadProxmoxConfig loads Proxmox configuration from file and environment variables
func LoadProxmoxConfig() (*types.AuthConfig, error) {
	config := &ProxmoxConfig{
		Timeout:   30,
		VerifySSL: false,
	}

	// Try to load from YAML file first
	configFile := "proxmox-auth.yaml"
	if _, err := os.Stat(configFile); err == nil {
		if err := loadFromYAML(configFile, config); err != nil {
			return nil, fmt.Errorf("failed to load config from %s: %w", configFile, err)
		}
	}

	// Override with environment variables if they exist
	loadFromEnv(config)

	// Validate required fields
	if config.Host == "" {
		return nil, fmt.Errorf("proxmox host is required (set PROXMOX_HOST env var or configure in %s)", configFile)
	}

	// Convert to AuthConfig
	authConfig := &types.AuthConfig{
		Proxmox: struct {
			Host       string `yaml:"host" json:"host"`
			User       string `yaml:"user" json:"user"`
			Password   string `yaml:"password" json:"password"`
			TokenName  string `yaml:"token_name" json:"token_name"`
			TokenValue string `yaml:"token_value" json:"token_value"`
			VerifySSL  bool   `yaml:"verify_ssl" json:"verify_ssl"`
			Timeout    int    `yaml:"timeout" json:"timeout"`
			Node       string `yaml:"node" json:"node"`
		}{
			Host:       config.Host,
			User:       config.User,
			Password:   config.Password,
			TokenName:  config.TokenName,
			TokenValue: config.TokenValue,
			VerifySSL:  config.VerifySSL,
			Timeout:    config.Timeout,
			Node:       config.Node,
		},
	}

	return authConfig, nil
}

// loadFromYAML loads configuration from YAML file
func loadFromYAML(filename string, config *ProxmoxConfig) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// Parse YAML into a map first to handle nested structure
	var yamlData map[string]interface{}
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		return err
	}

	// Extract proxmox section
	if proxmoxSection, ok := yamlData["proxmox"].(map[string]interface{}); ok {
		if host, ok := proxmoxSection["host"].(string); ok {
			config.Host = host
		}
		if user, ok := proxmoxSection["user"].(string); ok {
			config.User = user
		}
		if password, ok := proxmoxSection["password"].(string); ok {
			config.Password = password
		}
		if tokenName, ok := proxmoxSection["token_name"].(string); ok {
			config.TokenName = tokenName
		}
		if tokenValue, ok := proxmoxSection["token_value"].(string); ok {
			config.TokenValue = tokenValue
		}
		if verifySSL, ok := proxmoxSection["verify_ssl"].(bool); ok {
			config.VerifySSL = verifySSL
		}
		if timeout, ok := proxmoxSection["timeout"].(int); ok {
			config.Timeout = timeout
		}
		if node, ok := proxmoxSection["node"].(string); ok {
			config.Node = node
		}
	}

	return nil
}

// loadFromEnv loads configuration from environment variables
func loadFromEnv(config *ProxmoxConfig) {
	if host := os.Getenv("PROXMOX_HOST"); host != "" {
		config.Host = host
	}
	if user := os.Getenv("PROXMOX_USER"); user != "" {
		config.User = user
	}
	if password := os.Getenv("PROXMOX_PASSWORD"); password != "" {
		config.Password = password
	}
	if tokenName := os.Getenv("PROXMOX_TOKEN_NAME"); tokenName != "" {
		config.TokenName = tokenName
	}
	if tokenValue := os.Getenv("PROXMOX_TOKEN_VALUE"); tokenValue != "" {
		config.TokenValue = tokenValue
	}
	if verifySSL := os.Getenv("PROXMOX_VERIFY_SSL"); verifySSL != "" {
		config.VerifySSL = verifySSL == "true"
	}
	if timeout := os.Getenv("PROXMOX_TIMEOUT"); timeout != "" {
		if _, err := fmt.Sscanf(timeout, "%d", &config.Timeout); err != nil {
			// Timeout parse error handled by default value
		}
	}
	if node := os.Getenv("PROXMOX_NODE"); node != "" {
		config.Node = node
	}
}

// GetConfigPath returns the path to the configuration file
func GetConfigPath() string {
	// Check for config in current directory
	if _, err := os.Stat("proxmox-auth.yaml"); err == nil {
		return "proxmox-auth.yaml"
	}
	
	// Check for config in home directory
	homeDir, err := os.UserHomeDir()
	if err == nil {
		configPath := filepath.Join(homeDir, ".config", "mini-mcp", "proxmox-auth.yaml")
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}
	}
	
	// Default to current directory
	return "proxmox-auth.yaml"
}
