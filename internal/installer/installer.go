package installer

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"mini-mcp/internal/shared/logging"
)

// Installer handles building and installing the mini-mcp tool
type Installer struct {
	projectRoot string
	binaryName  string
	installPath string
}

// NewInstaller creates a new installer instance
func NewInstaller(projectRoot string) *Installer {
	binaryName := "mini-mcp"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	// Determine install path based on OS
	var installPath string
	switch runtime.GOOS {
	case "windows":
		installPath = filepath.Join(os.Getenv("USERPROFILE"), "bin", binaryName)
	case "darwin", "linux":
		installPath = filepath.Join("/usr", "local", "bin", binaryName)
	default:
		installPath = filepath.Join(os.Getenv("HOME"), "bin", binaryName)
	}

	return &Installer{
		projectRoot: projectRoot,
		binaryName:  binaryName,
		installPath: installPath,
	}
}

// Build builds the mini-mcp binary
func (i *Installer) Build() error {
	logger := logging.GetGlobalLogger()
	logger.Info("Building mini-mcp binary", map[string]any{
		"project_root": i.projectRoot,
		"binary_name":  i.binaryName,
	})

	// Change to project root
	if err := os.Chdir(i.projectRoot); err != nil {
		return fmt.Errorf("failed to change to project root: %w", err)
	}

	// Build the binary
	cmd := exec.Command("go", "build", "-o", i.binaryName, "./cmd/mini-mcp")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build binary: %w", err)
	}

	logger.Info("Binary built successfully", map[string]any{
		"binary_path": filepath.Join(i.projectRoot, i.binaryName),
	})

	return nil
}

// Install installs the binary to the system PATH
func (i *Installer) Install() error {
	logger := logging.GetGlobalLogger()

	// Check if binary exists
	binaryPath := filepath.Join(i.projectRoot, i.binaryName)
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return fmt.Errorf("binary not found at %s, please build first", binaryPath)
	}

	// Create install directory if it doesn't exist
	installDir := filepath.Dir(i.installPath)
	if err := os.MkdirAll(installDir, 0755); err != nil {
		return fmt.Errorf("failed to create install directory: %w", err)
	}

	// Copy binary to install location
	cmd := exec.Command("cp", binaryPath, i.installPath)
	if runtime.GOOS == "windows" {
		cmd = exec.Command("copy", binaryPath, i.installPath)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to copy binary: %w", err)
	}

	// Make binary executable on Unix systems
	if runtime.GOOS != "windows" {
		cmd = exec.Command("chmod", "+x", i.installPath)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to make binary executable: %w", err)
		}
	}

	logger.Info("Binary installed successfully", map[string]any{
		"install_path": i.installPath,
	})

	return nil
}

// Uninstall removes the binary from the system
func (i *Installer) Uninstall() error {
	logger := logging.GetGlobalLogger()

	if _, err := os.Stat(i.installPath); os.IsNotExist(err) {
		logger.Info("Binary not found, nothing to uninstall", map[string]any{
			"install_path": i.installPath,
		})
		return nil
	}

	if err := os.Remove(i.installPath); err != nil {
		return fmt.Errorf("failed to remove binary: %w", err)
	}

	logger.Info("Binary uninstalled successfully", map[string]any{
		"install_path": i.installPath,
	})

	return nil
}

// IsInstalled checks if the binary is installed
func (i *Installer) IsInstalled() bool {
	_, err := os.Stat(i.installPath)
	return err == nil
}

// GetInstallPath returns the install path
func (i *Installer) GetInstallPath() string {
	return i.installPath
}

// GetBinaryName returns the binary name
func (i *Installer) GetBinaryName() string {
	return i.binaryName
}

// UpdateVSCodeSettings updates VS Code settings.json with MCP configuration
func (i *Installer) UpdateVSCodeSettings() error {
	logger := logging.GetGlobalLogger()

	// Find VS Code settings directory
	var settingsPath string
	switch runtime.GOOS {
	case "windows":
		settingsPath = filepath.Join(os.Getenv("APPDATA"), "Code", "User", "settings.json")
	case "darwin":
		settingsPath = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "Code", "User", "settings.json")
	case "linux":
		settingsPath = filepath.Join(os.Getenv("HOME"), ".config", "Code", "User", "settings.json")
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// Check if settings.json exists
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		logger.Info("VS Code settings.json not found, creating new one", map[string]any{
			"settings_path": settingsPath,
		})

		// Create directory if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(settingsPath), 0755); err != nil {
			return fmt.Errorf("failed to create VS Code settings directory: %w", err)
		}
	}

	// Read existing settings
	settings := make(map[string]interface{})
	if data, err := os.ReadFile(settingsPath); err == nil {
		if err := json.Unmarshal(data, &settings); err != nil {
			logger.Warning("Failed to parse existing settings.json, creating new one", map[string]any{
				"settings_path": settingsPath,
				"error":         err.Error(),
			})
		}
	}

	// Update MCP servers configuration
	if settings["mcpServers"] == nil {
		settings["mcpServers"] = make(map[string]interface{})
	}

	mcpServers := settings["mcpServers"].(map[string]interface{})
	mcpServers["mini-mcp"] = map[string]interface{}{
		"command": i.installPath,
		"args":    []string{},
		"env": map[string]string{
			"ENVIRONMENT": "production",
			"LOG_LEVEL":   "INFO",
			"PORT":        ":8080",
		},
		"cwd":         i.projectRoot,
		"description": "Enhanced Mini MCP Infrastructure Management Tool",
	}

	// Write updated settings
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write settings: %w", err)
	}

	logger.Info("VS Code settings updated successfully", map[string]any{
		"settings_path": settingsPath,
		"mcp_server":    "mini-mcp",
	})

	return nil
}

// UpdateCursorSettings updates Cursor settings.json with MCP configuration
func (i *Installer) UpdateCursorSettings() error {
	logger := logging.GetGlobalLogger()

	// Find Cursor settings directory
	var settingsPath string
	switch runtime.GOOS {
	case "windows":
		settingsPath = filepath.Join(os.Getenv("APPDATA"), "Cursor", "User", "settings.json")
	case "darwin":
		settingsPath = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "Cursor", "User", "settings.json")
	case "linux":
		settingsPath = filepath.Join(os.Getenv("HOME"), ".config", "Cursor", "User", "settings.json")
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// Check if settings.json exists
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		logger.Info("Cursor settings.json not found, creating new one", map[string]any{
			"settings_path": settingsPath,
		})

		// Create directory if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(settingsPath), 0755); err != nil {
			return fmt.Errorf("failed to create Cursor settings directory: %w", err)
		}
	}

	// Read existing settings
	settings := make(map[string]interface{})
	if data, err := os.ReadFile(settingsPath); err == nil {
		if err := json.Unmarshal(data, &settings); err != nil {
			logger.Warning("Failed to parse existing settings.json, creating new one", map[string]any{
				"settings_path": settingsPath,
				"error":         err.Error(),
			})
		}
	}

	// Update MCP servers configuration
	if settings["mcpServers"] == nil {
		settings["mcpServers"] = make(map[string]interface{})
	}

	mcpServers := settings["mcpServers"].(map[string]interface{})
	mcpServers["mini-mcp"] = map[string]interface{}{
		"command": i.installPath,
		"args":    []string{},
		"env": map[string]string{
			"ENVIRONMENT": "production",
			"LOG_LEVEL":   "INFO",
			"PORT":        ":8080",
		},
		"cwd":         i.projectRoot,
		"description": "Enhanced Mini MCP Infrastructure Management Tool",
	}

	// Write updated settings
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write settings: %w", err)
	}

	logger.Info("Cursor settings updated successfully", map[string]any{
		"settings_path": settingsPath,
		"mcp_server":    "mini-mcp",
	})

	return nil
}

// InstallAndConfigure performs a complete installation with configuration
func (i *Installer) InstallAndConfigure() error {
	logger := logging.GetGlobalLogger()

	// Build the binary
	if err := i.Build(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	// Install the binary
	if err := i.Install(); err != nil {
		return fmt.Errorf("install failed: %w", err)
	}

	// Update VS Code settings
	if err := i.UpdateVSCodeSettings(); err != nil {
		logger.Warning("Failed to update VS Code settings", map[string]any{
			"error": err.Error(),
		})
	}

	// Update Cursor settings
	if err := i.UpdateCursorSettings(); err != nil {
		logger.Warning("Failed to update Cursor settings", map[string]any{
			"error": err.Error(),
		})
	}

	logger.Info("Installation and configuration completed successfully", map[string]any{
		"install_path": i.installPath,
		"binary_name":  i.binaryName,
	})

	return nil
}

// GetStatus returns the current installation status
func (i *Installer) GetStatus() map[string]interface{} {
	status := map[string]interface{}{
		"binary_name":  i.binaryName,
		"install_path": i.installPath,
		"is_installed": i.IsInstalled(),
		"project_root": i.projectRoot,
	}

	// Check if binary exists in project
	binaryPath := filepath.Join(i.projectRoot, i.binaryName)
	if _, err := os.Stat(binaryPath); err == nil {
		status["binary_exists"] = true
		status["binary_path"] = binaryPath
	} else {
		status["binary_exists"] = false
	}

	return status
}
