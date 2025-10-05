package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"mini-mcp/internal/installer"
	"mini-mcp/internal/shared/logging"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove the mini-mcp binary from system PATH",
	Long: `Remove the mini-mcp binary from your system PATH.

This command will:
1. Remove the binary from the system PATH
2. Leave project files intact

Examples:
  mini-mcp-cli uninstall                 # Remove binary from system`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logging.GetGlobalLogger()
		logger.Info("Starting uninstallation process", map[string]any{
			"project_root": projectRoot,
		})

		// Create installer instance
		inst := installer.NewInstaller(projectRoot)

		// Check if installed
		if !inst.IsInstalled() {
			logger.Info("Binary not installed, nothing to uninstall", map[string]any{})
			printInfo("Binary not installed, nothing to uninstall")
			return
		}

		// Uninstall the binary
		if err := inst.Uninstall(); err != nil {
			logger.Error("Uninstallation failed", err, map[string]any{})
			printError(fmt.Sprintf("Uninstallation failed: %v", err))
			os.Exit(1)
		}

		logger.Info("Uninstallation completed successfully", map[string]any{})
		printSuccess("Uninstallation completed successfully")
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
