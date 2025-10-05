package command

import (
	"context"
)

// Service defines the interface for command domain services
type Service interface {
	ExecuteCommand(ctx context.Context, cmd *Command) (*Result, error)
}

// ServiceImpl implements the command domain service
type ServiceImpl struct {
	repo Repository
}

// NewService creates a new command domain service
func NewService(repo Repository) Service {
	return &ServiceImpl{
		repo: repo,
	}
}

// ExecuteCommand executes a command
func (s *ServiceImpl) ExecuteCommand(ctx context.Context, cmd *Command) (*Result, error) {
	// Validate command
	if err := s.repo.Validate(ctx, cmd); err != nil {
		return nil, err
	}

	// Sanitize command
	if err := s.repo.Sanitize(ctx, cmd); err != nil {
		return nil, err
	}

	// Execute command
	result, err := s.repo.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return result, nil
}
