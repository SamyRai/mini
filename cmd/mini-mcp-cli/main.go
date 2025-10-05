package main

import (
	"os"

	"mini-mcp/cmd/mini-mcp-cli/cmd"
	"mini-mcp/internal/shared/logging"
)

func main() {
	// Set up logging
	logging.InitGlobalLogger(logging.LogLevelInfo)

	// Execute the root command
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
