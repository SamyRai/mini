package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"mini-mcp/internal/installer"
	"mini-mcp/internal/shared/logging"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show installation status",
	Long: `Show the current installation status of mini-mcp.

This command displays:
- Binary name and paths
- Installation status
- Project information

Examples:
  mini-mcp-cli status                    # Show current status`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logging.GetGlobalLogger()
		logger.Info("Checking installation status", map[string]any{
			"project_root": projectRoot,
		})

		// Create installer instance
		inst := installer.NewInstaller(projectRoot)

		// Get status
		status := inst.GetStatus()

		// Display status
		fmt.Println("Mini MCP Installation Status:")
		fmt.Println("=============================")
		fmt.Printf("Binary Name:     %s\n", status["binary_name"])
		fmt.Printf("Install Path:    %s\n", status["install_path"])
		fmt.Printf("Is Installed:    %v\n", status["is_installed"])
		fmt.Printf("Project Root:    %s\n", status["project_root"])
		fmt.Printf("Binary Exists:   %v\n", status["binary_exists"])
		
		if status["binary_path"] != nil {
			fmt.Printf("Binary Path:     %s\n", status["binary_path"])
		}

		// Additional status information
		if status["is_installed"].(bool) {
			fmt.Println("\n✅ mini-mcp is installed and available in PATH")
		} else {
			fmt.Println("\n❌ mini-mcp is not installed")
		}

		if status["binary_exists"].(bool) {
			fmt.Println("✅ Binary exists in project directory")
		} else {
			fmt.Println("❌ Binary not found in project directory")
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
