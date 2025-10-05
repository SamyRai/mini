package tools

import (
	"mini-mcp/internal/shared/validation"
)

// DuckDuckGoSearchArgs represents arguments for the duckduckgo_search tool.
// Example:
// {"query": "golang readability", "max_results": 3, "save_dir": "/tmp/mini-mcp", "timeout": 20}
type DuckDuckGoSearchArgs struct {
	// Query is the search query string
	Query string `json:"query"`
	// MaxResults limits number of result links to fetch (1-10, default 3)
	MaxResults int `json:"max_results,omitempty"`
	// SaveDir is an optional directory to save fetched page text files
	SaveDir string `json:"save_dir,omitempty"`
	// Timeout in seconds for each HTTP fetch (default 15, max 60)
	Timeout int `json:"timeout,omitempty"`
}

// Validate checks if the arguments are valid.
func (a *DuckDuckGoSearchArgs) Validate() error {
	if err := validation.StringRequired("query", a.Query); err != nil {
		return err
	}

	if a.MaxResults != 0 {
		if err := validation.Range[int](1, 10)("max_results", a.MaxResults); err != nil {
			return err
		}
	}

	if a.SaveDir != "" {
		if err := validation.Path("save_dir", a.SaveDir); err != nil {
			return err
		}
	}

	if a.Timeout != 0 {
		if err := validation.Range[int](1, 60)("timeout", a.Timeout); err != nil {
			return err
		}
	}
	return nil
}

// NewDuckDuckGoSearchArgs creates a new args object.
func NewDuckDuckGoSearchArgs(query string, maxResults int, saveDir string, timeout int) *DuckDuckGoSearchArgs {
	return &DuckDuckGoSearchArgs{
		Query:      query,
		MaxResults: maxResults,
		SaveDir:    saveDir,
		Timeout:    timeout,
	}
}
