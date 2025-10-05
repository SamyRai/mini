package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"mini-mcp/internal/installer"
	"mini-mcp/internal/shared/logging"
)

var (
	editor string
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure VS Code and Cursor settings",
	Long: `Configure VS Code and Cursor settings for mini-mcp.

This command will update the MCP server configuration in your editor settings.
You can specify which editor to configure or configure both.

Examples:
  mini-mcp-cli configure                  # Configure both VS Code and Cursor
  mini-mcp-cli configure --editor vscode  # Configure VS Code only
  mini-mcp-cli configure --editor cursor  # Configure Cursor only`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logging.GetGlobalLogger()
		logger.Info("Starting configuration process", map[string]any{
			"project_root": projectRoot,
			"editor":       editor,
		})

		// Create installer instance
		inst := installer.NewInstaller(projectRoot)

		// Configure based on editor selection
		if err := configureEditors(inst); err != nil {
			logger.Error("Configuration failed", err, map[string]any{})
			printError(fmt.Sprintf("Configuration failed: %v", err))
			os.Exit(1)
		}

		logger.Info("Configuration completed successfully", map[string]any{})
		printSuccess("Configuration completed successfully")
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)

	// Configure-specific flags
	configureCmd.Flags().StringVar(&editor, "editor", "", "Editor to configure (vscode, cursor, or both)")
}
