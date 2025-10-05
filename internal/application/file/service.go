package file

import (
	"context"

	"mini-mcp/internal/domain/file"
	"mini-mcp/internal/shared/logging"
	"mini-mcp/internal/shared/security"
)

// Service defines the interface for file application services
type Service interface {
	ReadFile(ctx context.Context, path string) (string, error)
	WriteFile(ctx context.Context, path, content string) error
	ListDirectory(ctx context.Context, path string) (string, error)
	DeleteFile(ctx context.Context, path string) error
}

// ServiceImpl implements the file service
type ServiceImpl struct {
	fileDomainService file.Service
}

// NewService creates a new file application service
func NewService(fileDomainService file.Service) Service {
	return &ServiceImpl{
		fileDomainService: fileDomainService,
	}
}

// NewServiceWithDeps creates a new file application service with dependencies
func NewServiceWithDeps(securityValidator security.PathValidator, logger logging.Logger) Service {
	// Create domain service with dependencies
	domainService := file.NewService(securityValidator, logger)

	// Create application service with domain service
	return &ServiceImpl{
		fileDomainService: domainService,
	}
}

// ReadFile reads a file through the domain service
func (s *ServiceImpl) ReadFile(ctx context.Context, path string) (string, error) {
	return s.fileDomainService.ReadFile(ctx, path)
}

// WriteFile writes content to a file through the domain service
func (s *ServiceImpl) WriteFile(ctx context.Context, path, content string) error {
	return s.fileDomainService.WriteFile(ctx, path, content)
}

// ListDirectory lists directory contents through the domain service
func (s *ServiceImpl) ListDirectory(ctx context.Context, path string) (string, error) {
	return s.fileDomainService.ListDirectory(ctx, path)
}

// DeleteFile deletes a file through the domain service
func (s *ServiceImpl) DeleteFile(ctx context.Context, path string) error {
	return s.fileDomainService.DeleteFile(ctx, path)
}
