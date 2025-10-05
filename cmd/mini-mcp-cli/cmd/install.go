package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"mini-mcp/internal/installer"
	"mini-mcp/internal/shared/logging"
)

var (
	configure bool
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the mini-mcp binary to system PATH",
	Long: `Install the mini-mcp binary to your system PATH and optionally configure
VS Code and Cursor settings.

This command will:
1. Build the binary (if not already built)
2. Install it to the system PATH
3. Optionally configure VS Code and Cursor settings

Examples:
  mini-mcp-cli install                    # Install binary only
  mini-mcp-cli install --configure       # Install and configure editors
  mini-mcp-cli install -c                # Short form of --configure`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logging.GetGlobalLogger()
		logger.Info("Starting installation process", map[string]any{
			"project_root": projectRoot,
			"configure":    configure,
		})

		// Create installer instance
		inst := installer.NewInstaller(projectRoot)

		// Check if binary exists, build if not
		status := inst.GetStatus()
		if !status["binary_exists"].(bool) {
			logger.Info("Binary not found, building first...", map[string]any{})
			printInfo("Binary not found, building first...")
			if err := inst.Build(); err != nil {
				logger.Error("Build failed", err, map[string]any{})
				printError(fmt.Sprintf("Build failed: %v", err))
				os.Exit(1)
			}
		}

		// Install the binary
		if err := inst.Install(); err != nil {
			logger.Error("Installation failed", err, map[string]any{})
			printError(fmt.Sprintf("Installation failed: %v", err))
			os.Exit(1)
		}

		// Configure editors if requested
		if configure {
			if err := configureEditors(inst); err != nil {
				logger.Warning("Configuration failed", map[string]any{
					"error": err.Error(),
				})
				printWarning(fmt.Sprintf("Configuration failed: %v", err))
			}
		}

		logger.Info("Installation completed successfully", map[string]any{
			"install_path": inst.GetInstallPath(),
		})
		printSuccess("Installation completed successfully")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Install-specific flags
	installCmd.Flags().BoolVarP(&configure, "configure", "c", false, "Configure VS Code and Cursor settings")
}
