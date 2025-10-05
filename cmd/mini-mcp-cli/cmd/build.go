package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"mini-mcp/internal/installer"
	"mini-mcp/internal/shared/logging"
)

var (
	outputPath string
	buildTags  string
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the mini-mcp binary",
	Long: `Build the mini-mcp binary from source code.

This command compiles the mini-mcp tool and creates an executable binary.
The binary will be created in the project root directory by default, or
in the specified output path.

Examples:
  mini-mcp-cli build                    # Build with default settings
  mini-mcp-cli build -o /usr/local/bin  # Build to specific directory
  mini-mcp-cli build --tags "dev"       # Build with specific tags`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logging.GetGlobalLogger()
		logger.Info("Starting build process", map[string]any{
			"project_root": projectRoot,
			"output_path":  outputPath,
			"build_tags":   buildTags,
		})

		// Create installer instance
		inst := installer.NewInstaller(projectRoot)

		// Override output path if specified
		if outputPath != "" {
			// This would require modifying the installer to accept custom output path
			logger.Warning("Custom output path not yet implemented, using default", map[string]any{
				"requested_path": outputPath,
			})
		}

		// Build the binary
		if err := inst.Build(); err != nil {
			logger.Error("Build failed", err, map[string]any{
				"project_root": projectRoot,
			})
			printError(fmt.Sprintf("Build failed: %v", err))
			os.Exit(1)
		}

		logger.Info("Build completed successfully", map[string]any{
			"binary_name": inst.GetBinaryName(),
		})
		printSuccess("Build completed successfully")
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Build-specific flags
	buildCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output path for the binary")
	buildCmd.Flags().StringVar(&buildTags, "tags", "", "Build tags to pass to go build")
}
