// Package tools contains types for MCP tool arguments and responses.
package tools

import (
	"mini-mcp/internal/shared/validation"
)

// FileOperation represents the type of file operation to perform
type FileOperation string

const (
	FileOpRead   FileOperation = "read"
	FileOpWrite  FileOperation = "write"
	FileOpList   FileOperation = "list"
	FileOpDelete FileOperation = "delete"
)

// FileOperationsArgs represents arguments for the file_operations tool.
// This tool provides comprehensive file system operations with advanced security features.
//
// Security Features:
// - Path validation and sanitization
// - Directory traversal attack prevention
// - Restricted access to system directories (/etc, /var, /usr)
// - Automatic directory creation for write operations
// - Permission handling and metadata preservation
//
// Operations:
// - read: Read file content with detailed metadata (size, permissions, timestamps)
// - write: Create or overwrite files with automatic directory creation
// - list: List directory contents with file information and total size
// - delete: Remove files or directories with safety checks
//
// Path Restrictions:
// - Cannot access system directories: /etc, /var, /usr, /bin, /sbin
// - Cannot use directory traversal: ../, ../
// - Supports relative and absolute paths
// - Automatic path normalization
//
// Examples:
//
//	{"operation": "read", "path": "/tmp/test.txt"}
//	{"operation": "write", "path": "/tmp/test.txt", "content": "Hello World"}
//	{"operation": "list", "path": "/tmp"}
//	{"operation": "delete", "path": "/tmp/test.txt"}
type FileOperationsArgs struct {
	// Operation is the file operation to perform
	// Must be one of: read, write, list, delete
	Operation FileOperation `json:"operation"`

	// Path is the file or directory path
	// Supports relative and absolute paths
	// Restricted from accessing system directories
	Path string `json:"path"`

	// Content is the content to write (for write operations)
	// Required when operation is "write"
	// Can be any string content (text, JSON, configuration, etc.)
	Content string `json:"content,omitempty"`
}

// Validate checks if the file operations arguments are valid.
func (args *FileOperationsArgs) Validate() error {
	// Validate operation
	if string(args.Operation) == "" {
		return validation.NewMissingRequiredError("operation")
	}
	if string(args.Operation) != "read" && string(args.Operation) != "write" &&
		string(args.Operation) != "list" && string(args.Operation) != "delete" {
		return validation.NewInvalidFormatError("operation", "must be 'read', 'write', 'list', or 'delete'")
	}

	// Validate path
	if err := validation.StringPath("path", args.Path); err != nil {
		return err
	}

	// Validate content for write operations
	if args.Operation == FileOpWrite {
		if args.Content == "" {
			return validation.NewMissingRequiredError("content")
		}
		if len(args.Content) > 1000000 { // 1MB max content
			return validation.NewInvalidFormatError("content", "content too large (max 1MB)")
		}
	}

	return nil
}

// NewFileOperationsArgs creates a new FileOperationsArgs with the given parameters.
func NewFileOperationsArgs(operation FileOperation, path string, content ...string) *FileOperationsArgs {
	args := &FileOperationsArgs{
		Operation: operation,
		Path:      path,
	}

	if len(content) > 0 {
		args.Content = content[0]
	}

	return args
}
