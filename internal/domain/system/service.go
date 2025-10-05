package system

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"
)

// Service defines the interface for system domain services
type Service interface {
	GetSystemInfo(ctx context.Context) (string, error)
	GetHealth(ctx context.Context) (string, error)
	GetMetrics(ctx context.Context) (string, error)
}

// ServiceImpl implements the system domain service
type ServiceImpl struct {
	// System service implementation with proper error handling
}

// NewService creates a new system domain service
func NewService() Service {
	return &ServiceImpl{}
}

// GetSystemInfo returns basic system information
func (s *ServiceImpl) GetSystemInfo(ctx context.Context) (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	wd, err := os.Getwd()
	if err != nil {
		wd = "unknown"
	}

	info := "System Information:\n"
	info += fmt.Sprintf("Hostname: %s\n", hostname)
	info += fmt.Sprintf("Working Directory: %s\n", wd)
	info += fmt.Sprintf("OS: %s %s\n", runtime.GOOS, runtime.GOARCH)
	info += fmt.Sprintf("CPU Cores: %d\n", runtime.NumCPU())
	info += fmt.Sprintf("Go Version: %s\n", runtime.Version())

	return info, nil
}

// GetHealth returns system health status
func (s *ServiceImpl) GetHealth(ctx context.Context) (string, error) {
	// Basic health checks
	health := "System Health: OK\n"

	// Check if we can read /proc/version (Linux) or similar
	if runtime.GOOS == "linux" {
		if _, err := os.Stat("/proc/version"); err == nil {
			health += "Linux kernel: Available\n"
		}
	}

	// Check memory usage (basic)
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	health += fmt.Sprintf("Memory Usage: %d MB\n", memStats.Alloc/1024/1024)

	return health, nil
}

// GetMetrics returns basic system metrics
func (s *ServiceImpl) GetMetrics(ctx context.Context) (string, error) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	metrics := "System Metrics:\n"
	metrics += fmt.Sprintf("Memory Allocated: %d MB\n", memStats.Alloc/1024/1024)
	metrics += fmt.Sprintf("Memory Total: %d MB\n", memStats.TotalAlloc/1024/1024)
	metrics += fmt.Sprintf("Memory System: %d MB\n", memStats.Sys/1024/1024)
	metrics += fmt.Sprintf("GC Cycles: %d\n", memStats.NumGC)
	metrics += fmt.Sprintf("Goroutines: %d\n", runtime.NumGoroutine())
	metrics += fmt.Sprintf("Uptime: %v\n", time.Since(time.Now().Add(-time.Duration(memStats.PauseTotalNs))))

	return metrics, nil
}
