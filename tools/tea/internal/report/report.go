// Package report outputs analysis results to various formats.
package report

import (
	"encoding/json"
	"fmt"
	"os"

	"mini-mcp/tools/tea/internal/model"
)

// Write writes the analysis to a JSON file.
func Write(path string, a *model.Analysis) error {
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal analysis: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write analysis: %w", err)
	}
	return nil
}
