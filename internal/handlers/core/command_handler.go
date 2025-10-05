package core

import (
	"context"
	"fmt"
	"time"

	appcommand "mini-mcp/internal/application/command"
	"mini-mcp/internal/shared/logging"
)

// CommandHandler handles command execution requests
type CommandHandler interface {
	ExecuteCommand(ctx context.Context, args map[string]any) (string, error)
}

// CommandHandlerImpl implements the CommandHandler interface
type CommandHandlerImpl struct {
	commandService appcommand.Service
	logger         logging.Logger
}

// NewCommandHandler creates a new command handler
func NewCommandHandler(commandService appcommand.Service, logger logging.Logger) CommandHandler {
	return &CommandHandlerImpl{
		commandService: commandService,
		logger:         logger,
	}
}

// ExecuteCommand executes a shell command securely
func (h *CommandHandlerImpl) ExecuteCommand(ctx context.Context, args map[string]any) (string, error) {
	start := time.Now()

	// Extract command from args
	command, ok := args["command"].(string)
	if !ok {
		h.logger.Error("Invalid command argument", fmt.Errorf("invalid command argument"), map[string]any{"args": args})
		return "", fmt.Errorf("invalid command argument")
	}

	// Extract timeout from args
	timeout := 30 // default timeout
	if timeoutVal, ok := args["timeout"].(float64); ok {
		timeout = int(timeoutVal)
	}

	// Record metrics
	defer func() {
		duration := time.Since(start)
		success := true

		// Get logger metrics
		if logger, ok := h.logger.(*logging.LoggerImpl); ok {
			logger.GetMetrics().RecordRequest("run", duration, success)
		}
	}()

	// Execute command
	result, err := h.commandService.ExecuteCommand(ctx, command, timeout)
	if err != nil {
		// Record error in metrics
		if logger, ok := h.logger.(*logging.LoggerImpl); ok {
			logger.GetMetrics().RecordError("run")
		}

		h.logger.Error("Command execution failed", err, map[string]any{
			"command":  command,
			"timeout":  timeout,
			"duration": time.Since(start).String(),
		})
		return "", err
	}

	h.logger.Info("Command executed successfully", map[string]any{
		"command":       command,
		"output_length": len(result),
		"duration":      time.Since(start).String(),
	})

	return result, nil
}
