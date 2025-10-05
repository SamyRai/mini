package file

import (
	"context"
	stderrors "errors"
	"os"
	"path/filepath"

	"mini-mcp/internal/shared/errors"
	"mini-mcp/internal/shared/logging"
	"mini-mcp/internal/shared/security"
)

var (
	ErrPathNotAllowed = stderrors.New("path not allowed")
	ErrInvalidPath    = stderrors.New("invalid path")
)

// Service defines the interface for file domain services
type Service interface {
	ReadFile(ctx context.Context, path string) (string, error)
	WriteFile(ctx context.Context, path, content string) error
	ListDirectory(ctx context.Context, path string) (string, error)
	DeleteFile(ctx context.Context, path string) error
}

// ServiceImpl implements the file domain service
type ServiceImpl struct {
	securityValidator security.PathValidator
	logger           logging.Logger
}

// NewService creates a new file domain service
func NewService(securityValidator security.PathValidator, logger logging.Logger) Service {
	return &ServiceImpl{
		securityValidator: securityValidator,
		logger:           logger,
	}
}

// ReadFile reads a file
func (s *ServiceImpl) ReadFile(ctx context.Context, path string) (string, error) {
	// Validate path using security validator
	if err := s.securityValidator.ValidatePath(path); err != nil {
		s.logger.Error("Path validation failed for read operation", err, map[string]any{
			"path": path,
		})
		return "", errors.WrapError(err, errors.ErrorCodePathBlocked, "Path validation failed for read operation")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			s.logger.Error("File not found", err, map[string]any{
				"path": path,
			})
			return "", errors.NewFileNotFoundError(path)
		}
		if os.IsPermission(err) {
			s.logger.Error("Permission denied", err, map[string]any{
				"path": path,
			})
			return "", errors.NewPermissionDeniedError(path)
		}

		s.logger.Error("Failed to read file", err, map[string]any{
			"path": path,
		})
		return "", errors.WrapError(err, errors.ErrorCodeInternalError, "Failed to read file")
	}

	s.logger.Debug("File read successfully", map[string]any{
		"path":         path,
		"content_size": len(content),
	})

	return string(content), nil
}

// WriteFile writes content to a file
func (s *ServiceImpl) WriteFile(ctx context.Context, path, content string) error {
	// Validate path using security validator
	if err := s.securityValidator.ValidatePath(path); err != nil {
		s.logger.Error("Path validation failed for write operation", err, map[string]any{
			"path": path,
		})
		return errors.WrapError(err, errors.ErrorCodePathBlocked, "Path validation failed for write operation")
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		s.logger.Error("Failed to create directory", err, map[string]any{
			"directory": dir,
		})
		return errors.WrapError(err, errors.ErrorCodeInternalError, "Failed to create directory")
	}

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		if os.IsPermission(err) {
			s.logger.Error("Permission denied", err, map[string]any{
				"path": path,
			})
			return errors.NewPermissionDeniedError(path)
		}

		s.logger.Error("Failed to write file", err, map[string]any{
			"path": path,
		})
		return errors.WrapError(err, errors.ErrorCodeInternalError, "Failed to write file")
	}

	s.logger.Debug("File written successfully", map[string]any{
		"path":         path,
		"content_size": len(content),
	})

	return nil
}

// ListDirectory lists directory contents
func (s *ServiceImpl) ListDirectory(ctx context.Context, path string) (string, error) {
	// Validate path using security validator
	if err := s.securityValidator.ValidatePath(path); err != nil {
		s.logger.Error("Path validation failed for list operation", err, map[string]any{
			"path": path,
		})
		return "", errors.WrapError(err, errors.ErrorCodePathBlocked, "Path validation failed for list operation")
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			s.logger.Error("Directory not found", err, map[string]any{
				"path": path,
			})
			return "", errors.NewFileNotFoundError(path)
		}
		if os.IsPermission(err) {
			s.logger.Error("Permission denied", err, map[string]any{
				"path": path,
			})
			return "", errors.NewPermissionDeniedError(path)
		}

		s.logger.Error("Failed to read directory", err, map[string]any{
			"path": path,
		})
		return "", errors.WrapError(err, errors.ErrorCodeInternalError, "Failed to read directory")
	}

	var result string
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			s.logger.Warning("Failed to get file info", map[string]any{
				"entry": entry.Name(),
			})
			continue
		}
		result += info.Name() + "\n"
	}

	s.logger.Debug("Directory listed successfully", map[string]any{
		"path":         path,
		"entry_count":  len(entries),
		"result_size":  len(result),
	})

	return result, nil
}

// DeleteFile deletes a file or directory
func (s *ServiceImpl) DeleteFile(ctx context.Context, path string) error {
	// Validate path using security validator
	if err := s.securityValidator.ValidatePath(path); err != nil {
		s.logger.Error("Path validation failed for delete operation", err, map[string]any{
			"path": path,
		})
		return errors.WrapError(err, errors.ErrorCodePathBlocked, "Path validation failed for delete operation")
	}

	err := os.RemoveAll(path)
	if err != nil {
		if os.IsNotExist(err) {
			s.logger.Error("File/directory not found", err, map[string]any{
				"path": path,
			})
			return errors.NewFileNotFoundError(path)
		}
		if os.IsPermission(err) {
			s.logger.Error("Permission denied", err, map[string]any{
				"path": path,
			})
			return errors.NewPermissionDeniedError(path)
		}

		s.logger.Error("Failed to delete file/directory", err, map[string]any{
			"path": path,
		})
		return errors.WrapError(err, errors.ErrorCodeInternalError, "Failed to delete file/directory")
	}

	s.logger.Debug("File/directory deleted successfully", map[string]any{
		"path": path,
	})

	return nil
}
