package core

import (
	"context"
	"fmt"

	appfile "mini-mcp/internal/application/file"
	"mini-mcp/internal/shared/logging"
)

// FileHandler handles file operation requests
type FileHandler interface {
	ReadFile(ctx context.Context, args map[string]any) (string, error)
	WriteFile(ctx context.Context, args map[string]any) (string, error)
	ListDirectory(ctx context.Context, args map[string]any) (string, error)
	DeleteFile(ctx context.Context, args map[string]any) (string, error)
}

// FileHandlerImpl implements the FileHandler interface
type FileHandlerImpl struct {
	fileService appfile.Service
	logger      logging.Logger
}

// NewFileHandler creates a new file handler
func NewFileHandler(fileService appfile.Service, logger logging.Logger) FileHandler {
	return &FileHandlerImpl{
		fileService: fileService,
		logger:      logger,
	}
}

// ReadFile reads a file
func (h *FileHandlerImpl) ReadFile(ctx context.Context, args map[string]any) (string, error) {
	path, ok := args["path"].(string)
	if !ok {
		h.logger.Error("Invalid path argument", fmt.Errorf("invalid path argument"), map[string]any{"args": args})
		return "", fmt.Errorf("invalid path argument")
	}

	result, err := h.fileService.ReadFile(ctx, path)
	if err != nil {
		h.logger.Error("File read failed", err, map[string]any{"path": path})
		return "", err
	}

	h.logger.Info("File read successfully", map[string]any{"path": path})
	return result, nil
}

// WriteFile writes content to a file
func (h *FileHandlerImpl) WriteFile(ctx context.Context, args map[string]any) (string, error) {
	path, ok := args["path"].(string)
	if !ok {
		h.logger.Error("Invalid path argument", fmt.Errorf("invalid path argument"), map[string]any{"args": args})
		return "", fmt.Errorf("invalid path argument")
	}

	content, ok := args["content"].(string)
	if !ok {
		h.logger.Error("Invalid content argument", fmt.Errorf("invalid content argument"), map[string]any{"args": args})
		return "", fmt.Errorf("invalid content argument")
	}

	err := h.fileService.WriteFile(ctx, path, content)
	if err != nil {
		h.logger.Error("File write failed", err, map[string]any{"path": path})
		return "", err
	}

	h.logger.Info("File written successfully", map[string]any{"path": path})
	return "File written successfully", nil
}

// ListDirectory lists directory contents
func (h *FileHandlerImpl) ListDirectory(ctx context.Context, args map[string]any) (string, error) {
	path, ok := args["path"].(string)
	if !ok {
		h.logger.Error("Invalid path argument", fmt.Errorf("invalid path argument"), map[string]any{"args": args})
		return "", fmt.Errorf("invalid path argument")
	}

	result, err := h.fileService.ListDirectory(ctx, path)
	if err != nil {
		h.logger.Error("Directory listing failed", err, map[string]any{"path": path})
		return "", err
	}

	h.logger.Info("Directory listed successfully", map[string]any{"path": path})
	return result, nil
}

// DeleteFile deletes a file or directory
func (h *FileHandlerImpl) DeleteFile(ctx context.Context, args map[string]any) (string, error) {
	path, ok := args["path"].(string)
	if !ok {
		h.logger.Error("Invalid path argument", fmt.Errorf("invalid path argument"), map[string]any{"args": args})
		return "", fmt.Errorf("invalid path argument")
	}

	err := h.fileService.DeleteFile(ctx, path)
	if err != nil {
		h.logger.Error("File deletion failed", err, map[string]any{"path": path})
		return "", err
	}

	h.logger.Info("File deleted successfully", map[string]any{"path": path})
	return "File deleted successfully", nil
}
