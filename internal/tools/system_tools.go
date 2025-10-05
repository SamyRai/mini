package tools

import (
	"context"
	"encoding/json"

	"mini-mcp/internal/handlers/core"
	"mini-mcp/internal/health"
	"mini-mcp/internal/registry"
	"mini-mcp/internal/shared/logging"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterSystemTools registers system-related tools
func RegisterSystemTools(server *mcp.Server, toolRegistry *registry.TypeSafeToolRegistry, systemHandler *core.SystemHandlerImpl, healthChecker *health.HealthChecker, logger logging.Logger) {
	// system - Get system information
	systemBuilder := registry.NewToolBuilder[SystemInfoArgs](toolRegistry, "system", "Get system information and status")

	systemBuilder.
		WithHandler(func(ctx context.Context, req *mcp.CallToolRequest, args SystemInfoArgs) (*mcp.CallToolResult, any, error) {
			info, err := systemHandler.GetSystemInfo(ctx, map[string]any{
				"include_metrics": args.IncludeMetrics,
				"include_health":  args.IncludeHealth,
			})
			if err != nil {
				errorResult, _, _ := toolRegistry.CreateErrorResult(err.Error(), map[string]any{
					"include_metrics": args.IncludeMetrics,
					"include_health":  args.IncludeHealth,
				})
				return errorResult, nil, nil
			}

			successResult, _, _ := toolRegistry.CreateTextResult(info)
			return successResult, nil, nil
		}).
		WithValidator(func(args SystemInfoArgs) error {
			return nil
		})

	if err := systemBuilder.Register(); err != nil {
		// Log error but continue - tool registration failure should not crash the server
		return
	}

	// metrics - Get application metrics and performance data
	metricsBuilder := registry.NewToolBuilder[MetricsArgs](toolRegistry, "metrics", "Get application metrics, performance data, and observability information")

	metricsBuilder.
		WithHandler(func(ctx context.Context, req *mcp.CallToolRequest, args MetricsArgs) (*mcp.CallToolResult, any, error) {
			// Get metrics from logger
			metrics := logger.GetMetrics()
			metricsData := metrics.GetMetricsSummary()

			// Add health check status if requested
			if args.IncludeHealth && healthChecker != nil {
				healthInfo := healthChecker.CheckHealth(ctx)
				metricsData["health_status"] = healthInfo.Status
				metricsData["health_checks"] = healthInfo.Checks
			}

			// Convert to JSON for better readability
			jsonData, err := json.MarshalIndent(metricsData, "", "  ")
			if err != nil {
				errorResult, _, _ := toolRegistry.CreateErrorResult("Failed to format metrics", map[string]any{
					"error": err.Error(),
				})
				return errorResult, nil, nil
			}

			successResult, _, _ := toolRegistry.CreateTextResult(string(jsonData))
			return successResult, nil, nil
		}).
		WithValidator(func(args MetricsArgs) error {
			return nil
		})

	if err := metricsBuilder.Register(); err != nil {
		// Log error but continue - tool registration failure should not crash the server
		return
	}
}

// SystemInfoArgs represents arguments for system information
type SystemInfoArgs struct {
	IncludeMetrics bool `json:"include_metrics,omitempty" jsonschema:"Include performance metrics"`
	IncludeHealth  bool `json:"include_health,omitempty" jsonschema:"Include health check information"`
}

// MetricsArgs represents arguments for metrics
type MetricsArgs struct {
	IncludeHealth bool `json:"include_health,omitempty" jsonschema:"Include health check information"`
}
