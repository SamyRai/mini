package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mini-mcp/internal/health"
	"mini-mcp/internal/server"
	"mini-mcp/internal/shared/logging"
	"mini-mcp/internal/shared/security"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
    version := flag.String("version", "dev", "Version of the mini-mcp server")
    logLevel := flag.String("log-level", "INFO", "Log level: DEBUG, INFO, WARNING, ERROR, FATAL")
    flag.Parse()

    // Initialize global logger
    lvl := logging.LogLevel(*logLevel)
    logging.InitGlobalLogger(lvl)
    logger := logging.GetGlobalLogger()

    // Logger initialization handled by structured logging

    // Log startup
    logger.Info("Starting mini-mcp server", map[string]any{
        "version":   *version,
        "log_level": string(lvl),
    })

    // Security executor initialization handled by structured logging
    // Initialize security executor with defaults
    sec := security.NewSecureCommandExecutor(nil)
    // Security executor initialization handled by structured logging
    logger.Info("Security executor initialized", nil)

    // Create health checker
    healthChecker := health.CreateDefaultHealthChecker(*version)

    deps := server.Deps{
        Logger:        logger,
        Security:      sec,
        HealthChecker: healthChecker,
    }

    // Server build handled by structured logging
    s := server.BuildServer(deps, *version)
    // Server build success handled by structured logging
    logger.Info("MCP server built successfully", map[string]any{
        "capabilities": "core_tools,infrastructure_tools,port_process_tools,resources",
    })

    // Set up graceful shutdown
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Handle shutdown signals
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    // Start server in a goroutine
    serverErr := make(chan error, 1)
    go func() {
        logger.Info("Starting MCP server on stdio transport", nil)
        if err := s.Run(ctx, &mcp.StdioTransport{}); err != nil {
            serverErr <- err
        }
    }()

    // Wait for shutdown signal or server error
    select {
    case sig := <-sigChan:
        logger.Info("Received shutdown signal, initiating graceful shutdown", map[string]any{
            "signal": sig.String(),
        })

        // Cancel context to signal shutdown
        cancel()

        // Give the server time to shutdown gracefully
        shutdownTimeout := 10 * time.Second
        shutdownTimer := time.NewTimer(shutdownTimeout)
       	defer shutdownTimer.Stop()

       	select {
       	case <-shutdownTimer.C:
           	logger.Warning("Shutdown timeout reached, forcing exit", nil)
           	os.Exit(1)
       	default:
       	}

    case err := <-serverErr:
        logger.Error("MCP server failed", err, nil)
        fmt.Fprintf(os.Stderr, "failed to run mcp server: %v\n", err)
        os.Exit(1)
    }

    // Perform cleanup
    performCleanup(logger, sec)

    logger.Info("MCP server shutdown gracefully", nil)
}

// performCleanup handles resource cleanup during shutdown
func performCleanup(logger logging.Logger, security *security.SecureCommandExecutor) {
	logger.Info("Performing cleanup before shutdown", nil)

	// Cleanup security executor (kill any running processes)
	if security != nil {
		logger.Debug("Cleaning up security executor", nil)
		security.Cleanup()
	}

	// Log final metrics
	if metrics := logger.GetMetrics(); metrics != nil {
		logger.Info("Final metrics summary", metrics.GetMetricsSummary())
	}

	// Flush any pending logs
	time.Sleep(100 * time.Millisecond) // Brief pause to ensure logs are written

	logger.Info("Cleanup completed", nil)
}
