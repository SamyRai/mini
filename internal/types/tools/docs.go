package tools

import (
	"mini-mcp/internal/shared/validation"
)

// DocArgs represents arguments for get_documentation.
// Example:
//
//	{"name": "docker", "type": "command"}
//	{"name": "https://github.com/user/repo", "type": "git"}
type DocArgs struct {
	// Name is the name of the command or Git repository to get documentation for
	Name string `json:"name"`
	// Type is the type of documentation to get ("command" or "git")
	// If not specified, defaults to "command"
	Type string `json:"type,omitempty"`
}

// Validate checks if the documentation arguments are valid.
func (args *DocArgs) Validate() error {
	// Validate name
	if err := validation.StringRequired("name", args.Name); err != nil {
		return err
	}
	
	// Validate type if provided
	if args.Type != "" {
		if args.Type != "command" && args.Type != "git" {
			return validation.NewInvalidFormatError("type", "must be 'command' or 'git'")
		}
	}
	
	return nil
}

// NewDocArgs creates a new DocArgs with the given name and optional type.
func NewDocArgs(name, docType string) *DocArgs {
	return &DocArgs{
		Name: name,
		Type: docType,
	}
}
