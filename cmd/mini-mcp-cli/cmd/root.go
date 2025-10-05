package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"mini-mcp/internal/shared/logging"
)

var (
	projectRoot string
	logLevel    string
	verbose     bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mini-mcp-cli",
	Short: "Mini MCP CLI - Build, install, and manage the mini-mcp tool",
	Long: `Mini MCP CLI is a comprehensive tool for building, installing, and managing
the mini-mcp infrastructure management tool. It provides commands for building
the binary, installing it to your system PATH, configuring VS Code and Cursor
settings, and managing the installation.

Examples:
  mini-mcp-cli build                    # Build the mini-mcp binary
  mini-mcp-cli install                  # Install to system PATH
  mini-mcp-cli install --configure     # Install and configure editors
  mini-mcp-cli status                   # Show installation status
  mini-mcp-cli uninstall                # Remove from system`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Set up logging based on flags
		var lvl logging.LogLevel
		switch logLevel {
		case "DEBUG":
			lvl = logging.LogLevelDebug
		case "INFO":
			lvl = logging.LogLevelInfo
		case "WARN":
			lvl = logging.LogLevelWarning
		case "ERROR":
			lvl = logging.LogLevelError
		default:
			lvl = logging.LogLevelInfo
		}

		// Enable verbose logging if requested
		if verbose {
			lvl = logging.LogLevelDebug
		}

		logging.InitGlobalLogger(lvl)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&projectRoot, "project-root", "", "Project root directory (default: current directory)")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "INFO", "Log level (DEBUG, INFO, WARN, ERROR)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")

	// Set default project root to current directory
	cobra.OnInitialize(func() {
		if projectRoot == "" {
			wd, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
				os.Exit(1)
			}
			projectRoot = wd
		}

		// Resolve absolute path
		absPath, err := filepath.Abs(projectRoot)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error resolving project root: %v\n", err)
			os.Exit(1)
		}
		projectRoot = absPath

		// Validate project root (only for commands that need it)
		// Skip validation for help and version commands
		if len(os.Args) > 1 && os.Args[1] != "help" && os.Args[1] != "--help" && os.Args[1] != "-h" {
			if err := validateProjectRoot(); err != nil {
				fmt.Fprintf(os.Stderr, "Invalid project root: %v\n", err)
				os.Exit(1)
			}
		}
	})
}
