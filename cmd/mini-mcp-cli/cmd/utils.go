package cmd

import (
	"fmt"
	"os"

	"mini-mcp/internal/installer"
	"mini-mcp/internal/shared/logging"
)

// configureEditors configures VS Code and Cursor settings based on the editor flag
func configureEditors(inst *installer.Installer) error {
	logger := logging.GetGlobalLogger()
	var errors []error

	// Configure VS Code if requested or if no specific editor specified
	if editor == "" || editor == "vscode" || editor == "both" {
		logger.Info("Configuring VS Code settings", map[string]any{})
		if err := inst.UpdateVSCodeSettings(); err != nil {
			logger.Warning("VS Code configuration failed", map[string]any{
				"error": err.Error(),
			})
			errors = append(errors, fmt.Errorf("vscode: %w", err))
		} else {
			fmt.Println("✅ VS Code settings updated successfully")
		}
	}

	// Configure Cursor if requested or if no specific editor specified
	if editor == "" || editor == "cursor" || editor == "both" {
		logger.Info("Configuring Cursor settings", map[string]any{})
		if err := inst.UpdateCursorSettings(); err != nil {
			logger.Warning("Cursor configuration failed", map[string]any{
				"error": err.Error(),
			})
			errors = append(errors, fmt.Errorf("cursor: %w", err))
		} else {
			fmt.Println("✅ Cursor settings updated successfully")
		}
	}

	// Return error if any configuration failed
	if len(errors) > 0 {
		return fmt.Errorf("configuration errors: %v", errors)
	}

	return nil
}

// validateProjectRoot validates that the project root exists and contains necessary files
func validateProjectRoot() error {
	// Check if project root exists
	if _, err := os.Stat(projectRoot); os.IsNotExist(err) {
		return fmt.Errorf("project root does not exist: %s", projectRoot)
	}

	// Check for go.mod file
	goModPath := fmt.Sprintf("%s/go.mod", projectRoot)
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return fmt.Errorf("go.mod not found in project root: %s", projectRoot)
	}

	// Check for main.go file
	mainGoPath := fmt.Sprintf("%s/cmd/mini-mcp/main.go", projectRoot)
	if _, err := os.Stat(mainGoPath); os.IsNotExist(err) {
		return fmt.Errorf("main.go not found in expected location: %s", mainGoPath)
	}

	return nil
}

// printSuccess prints a success message with consistent formatting
func printSuccess(message string) {
	fmt.Printf("✅ %s\n", message)
}

// printError prints an error message with consistent formatting
func printError(message string) {
	fmt.Fprintf(os.Stderr, "❌ %s\n", message)
}

// printWarning prints a warning message with consistent formatting
func printWarning(message string) {
	fmt.Printf("⚠️  %s\n", message)
}

// printInfo prints an info message with consistent formatting
func printInfo(message string) {
	fmt.Printf("ℹ️  %s\n", message)
}
