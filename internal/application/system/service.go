package system

import (
	"context"

	"mini-mcp/internal/domain/system"
)

// Service defines the interface for system application services
type Service interface {
	GetSystemInfo(ctx context.Context) (string, error)
	GetHealth(ctx context.Context) (string, error)
	GetMetrics(ctx context.Context) (string, error)
}

// ServiceImpl implements the system service
type ServiceImpl struct {
	systemDomainService system.Service
}

// NewService creates a new system application service
func NewService(systemDomainService system.Service) Service {
	return &ServiceImpl{
		systemDomainService: systemDomainService,
	}
}

// GetSystemInfo gets system information through the domain service
func (s *ServiceImpl) GetSystemInfo(ctx context.Context) (string, error) {
	return s.systemDomainService.GetSystemInfo(ctx)
}

// GetHealth gets system health through the domain service
func (s *ServiceImpl) GetHealth(ctx context.Context) (string, error) {
	return s.systemDomainService.GetHealth(ctx)
}

// GetMetrics gets system metrics through the domain service
func (s *ServiceImpl) GetMetrics(ctx context.Context) (string, error) {
	return s.systemDomainService.GetMetrics(ctx)
}
