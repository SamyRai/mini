package core

import (
	"context"

	appsystem "mini-mcp/internal/application/system"
	"mini-mcp/internal/shared/logging"
)

// SystemHandler handles system information requests
type SystemHandler interface {
	GetSystemInfo(ctx context.Context, args map[string]any) (string, error)
	GetHealth(ctx context.Context, args map[string]any) (string, error)
	GetMetrics(ctx context.Context, args map[string]any) (string, error)
}

// SystemHandlerImpl implements the SystemHandler interface
type SystemHandlerImpl struct {
	systemService appsystem.Service
	logger        logging.Logger
}

// NewSystemHandler creates a new system handler
func NewSystemHandler(systemService appsystem.Service, logger logging.Logger) SystemHandler {
	return &SystemHandlerImpl{
		systemService: systemService,
		logger:        logger,
	}
}

// GetSystemInfo returns system information
func (h *SystemHandlerImpl) GetSystemInfo(ctx context.Context, args map[string]any) (string, error) {
	result, err := h.systemService.GetSystemInfo(ctx)
	if err != nil {
		h.logger.Error("System info retrieval failed", err, nil)
		return "", err
	}

	h.logger.Info("System info retrieved successfully", nil)
	return result, nil
}

// GetHealth returns system health status
func (h *SystemHandlerImpl) GetHealth(ctx context.Context, args map[string]any) (string, error) {
	result, err := h.systemService.GetHealth(ctx)
	if err != nil {
		h.logger.Error("Health check failed", err, nil)
		return "", err
	}

	h.logger.Info("Health check completed successfully", nil)
	return result, nil
}

// GetMetrics returns system metrics
func (h *SystemHandlerImpl) GetMetrics(ctx context.Context, args map[string]any) (string, error) {
	result, err := h.systemService.GetMetrics(ctx)
	if err != nil {
		h.logger.Error("Metrics retrieval failed", err, nil)
		return "", err
	}

	h.logger.Info("Metrics retrieved successfully", nil)
	return result, nil
}
