package command

import (
	"context"

	"mini-mcp/internal/domain/command"
)

// Service defines the interface for command application services
type Service interface {
	ExecuteCommand(ctx context.Context, command string, timeout int) (string, error)
}

// ServiceImpl implements the command service
type ServiceImpl struct {
	commandDomainService command.Service
}

// NewService creates a new command application service
func NewService(commandDomainService command.Service) Service {
	return &ServiceImpl{
		commandDomainService: commandDomainService,
	}
}

// ExecuteCommand executes a command through the domain service
func (s *ServiceImpl) ExecuteCommand(ctx context.Context, commandStr string, timeout int) (string, error) {
	// Create domain command
	cmd := &command.Command{
		Command: commandStr,
		Timeout: timeout,
	}

	// Execute through domain service
	result, err := s.commandDomainService.ExecuteCommand(ctx, cmd)
	if err != nil {
		return "", err
	}

	return result.Output, nil
}
