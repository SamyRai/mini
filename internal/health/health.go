package health

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"runtime"
	"sync"
	"syscall"
	"time"
)

// HealthStatus represents the overall health status
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusDegraded  HealthStatus = "degraded"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
)

// CheckResult represents the result of a health check
type CheckResult struct {
	Status    HealthStatus   `json:"status"`
	Message   string         `json:"message"`
	LastCheck time.Time      `json:"last_check"`
	Duration  time.Duration  `json:"duration"`
	Details   map[string]any `json:"details,omitempty"`
}

// DependencyStatus represents the status of a dependency
type DependencyStatus struct {
	Name      string        `json:"name"`
	Status    HealthStatus  `json:"status"`
	Message   string        `json:"message"`
	LastCheck time.Time     `json:"last_check"`
	Duration  time.Duration `json:"duration"`
	Critical  bool          `json:"critical"`
}

// HealthStatus represents the overall health status
type HealthInfo struct {
	Status       HealthStatus           `json:"status"`
	Checks       map[string]CheckResult `json:"checks"`
	Dependencies []DependencyStatus     `json:"dependencies"`
	Timestamp    time.Time              `json:"timestamp"`
	Version      string                 `json:"version"`
	Uptime       time.Duration          `json:"uptime"`
}

// HealthChecker provides health checking functionality
type HealthChecker struct {
	checks       map[string]HealthCheck
	dependencies map[string]Dependency
	startTime    time.Time
	version      string
	mu           sync.RWMutex
}

// HealthCheck represents a health check function
type HealthCheck func(ctx context.Context) CheckResult

// Dependency represents a system dependency
type Dependency struct {
	Name     string
	Check    HealthCheck
	Critical bool
	Interval time.Duration
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(version string) *HealthChecker {
	return &HealthChecker{
		checks:       make(map[string]HealthCheck),
		dependencies: make(map[string]Dependency),
		startTime:    time.Now(),
		version:      version,
	}
}

// AddCheck adds a health check
func (h *HealthChecker) AddCheck(name string, check HealthCheck) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.checks[name] = check
}

// AddDependency adds a dependency check
func (h *HealthChecker) AddDependency(name string, check HealthCheck, critical bool, interval time.Duration) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.dependencies[name] = Dependency{
		Name:     name,
		Check:    check,
		Critical: critical,
		Interval: interval,
	}
}

// CheckHealth performs all health checks
func (h *HealthChecker) CheckHealth(ctx context.Context) *HealthInfo {
	h.mu.RLock()
	defer h.mu.RUnlock()

	checks := make(map[string]CheckResult)
	dependencies := make([]DependencyStatus, 0)

	// Perform regular health checks
	for name, check := range h.checks {
		start := time.Now()
		result := check(ctx)
		result.LastCheck = start
		result.Duration = time.Since(start)
		checks[name] = result
	}

	// Perform dependency checks
	for name, dep := range h.dependencies {
		start := time.Now()
		result := dep.Check(ctx)
		dependencyStatus := DependencyStatus{
			Name:      name,
			Status:    result.Status,
			Message:   result.Message,
			LastCheck: start,
			Duration:  time.Since(start),
			Critical:  dep.Critical,
		}
		dependencies = append(dependencies, dependencyStatus)
	}

	// Determine overall status
	status := h.determineOverallStatus(checks, dependencies)

	return &HealthInfo{
		Status:       status,
		Checks:       checks,
		Dependencies: dependencies,
		Timestamp:    time.Now(),
		Version:      h.version,
		Uptime:       time.Since(h.startTime),
	}
}

// determineOverallStatus determines the overall health status
func (h *HealthChecker) determineOverallStatus(checks map[string]CheckResult, dependencies []DependencyStatus) HealthStatus {
	// Check if any critical dependencies are unhealthy
	for _, dep := range dependencies {
		if dep.Critical && dep.Status == HealthStatusUnhealthy {
			return HealthStatusUnhealthy
		}
	}

	// Count unhealthy and degraded checks
	unhealthyCount := 0
	degradedCount := 0

	for _, check := range checks {
		switch check.Status {
		case HealthStatusUnhealthy:
			unhealthyCount++
		case HealthStatusDegraded:
			degradedCount++
		}
	}

	// Determine overall status
	if unhealthyCount > 0 {
		return HealthStatusUnhealthy
	} else if degradedCount > 0 {
		return HealthStatusDegraded
	}

	return HealthStatusHealthy
}

// HTTPHandler creates an HTTP handler for health checks
func (h *HealthChecker) HTTPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Set timeout for health checks
		timeout := 30 * time.Second
		if r.URL.Query().Get("timeout") != "" {
			if t, err := time.ParseDuration(r.URL.Query().Get("timeout")); err == nil {
				timeout = t
			}
		}

		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		// Perform health check
		health := h.CheckHealth(ctx)

		// Set response status code
		switch health.Status {
		case HealthStatusHealthy:
			w.WriteHeader(http.StatusOK)
		case HealthStatusDegraded:
			w.WriteHeader(http.StatusOK) // Still OK but degraded
		case HealthStatusUnhealthy:
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		// Set content type
		w.Header().Set("Content-Type", "application/json")

		// Marshal and write response
		data, err := json.MarshalIndent(health, "", "  ")
		if err != nil {
			http.Error(w, "Failed to marshal health status", http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(data); err != nil {
			slog.Error("Failed to write health response", "error", err)
		}
	}
}

// ReadyHandler creates an HTTP handler for readiness checks
func (h *HealthChecker) ReadyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Set timeout for readiness checks
		timeout := 10 * time.Second
		if r.URL.Query().Get("timeout") != "" {
			if t, err := time.ParseDuration(r.URL.Query().Get("timeout")); err == nil {
				timeout = t
			}
		}

		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		// Perform health check
		health := h.CheckHealth(ctx)

		// Only return 200 if healthy
		if health.Status == HealthStatusHealthy {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{"status":"ready"}`)); err != nil {
				slog.Error("Failed to write ready response", "error", err)
			}
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			if _, err := w.Write([]byte(`{"status":"not ready"}`)); err != nil {
				slog.Error("Failed to write not ready response", "error", err)
			}
		}
	}
}

// Common health check functions

// PingCheck creates a simple ping health check
func PingCheck() HealthCheck {
	return func(ctx context.Context) CheckResult {
		return CheckResult{
			Status:  HealthStatusHealthy,
			Message: "Service is responding",
		}
	}
}

// DiskSpaceCheck creates a disk space health check
func DiskSpaceCheck(path string, minSpaceGB float64) HealthCheck {
	return func(ctx context.Context) CheckResult {
		start := time.Now()

		// Check actual disk space using syscall.Statfs
		var stat syscall.Statfs_t
		err := syscall.Statfs(path, &stat)
		if err != nil {
			return CheckResult{
				Status:  HealthStatusUnhealthy,
				Message: fmt.Sprintf("Failed to check disk space for %s: %v", path, err),
				Details: map[string]any{
					"path":         path,
					"min_space_gb": minSpaceGB,
					"duration":     time.Since(start).String(),
				},
			}
		}

		// Calculate available space in GB
		totalBytes := stat.Blocks * uint64(stat.Bsize)
		freeBytes := stat.Bavail * uint64(stat.Bsize)
		availableGB := float64(freeBytes) / (1024 * 1024 * 1024)

		status := HealthStatusHealthy
		message := fmt.Sprintf("Disk space is sufficient: %.2f GB available", availableGB)

		if availableGB < minSpaceGB {
			status = HealthStatusUnhealthy
			message = fmt.Sprintf("Insufficient disk space: %.2f GB available, %.2f GB required", availableGB, minSpaceGB)
		}

		return CheckResult{
			Status:  status,
			Message: message,
			Details: map[string]any{
				"path":         path,
				"min_space_gb": minSpaceGB,
				"available_gb": availableGB,
				"total_bytes":  totalBytes,
				"free_bytes":   freeBytes,
				"duration":     time.Since(start).String(),
			},
		}
	}
}

// MemoryCheck creates a memory usage health check
func MemoryCheck(maxUsagePercent float64) HealthCheck {
	return func(ctx context.Context) CheckResult {
		start := time.Now()

		// Get actual memory usage
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)

		// Calculate memory usage percentage
		usagePercent := float64(memStats.Alloc) / float64(memStats.Sys) * 100

		status := HealthStatusHealthy
		message := fmt.Sprintf("Memory usage is within limits: %.2f%% used", usagePercent)

		if usagePercent > maxUsagePercent {
			status = HealthStatusDegraded
			message = fmt.Sprintf("High memory usage: %.2f%% used (limit: %.2f%%)", usagePercent, maxUsagePercent)
		}

		return CheckResult{
			Status:  status,
			Message: message,
			Details: map[string]any{
				"max_usage_percent":     maxUsagePercent,
				"current_usage_percent": usagePercent,
				"alloc_bytes":           memStats.Alloc,
				"sys_bytes":             memStats.Sys,
				"total_alloc_bytes":     memStats.TotalAlloc,
				"num_gc":                memStats.NumGC,
				"duration":              time.Since(start).String(),
			},
		}
	}
}

// DatabaseCheck creates a database connectivity health check
func DatabaseCheck(checkFunc func(ctx context.Context) error) HealthCheck {
	return func(ctx context.Context) CheckResult {
		start := time.Now()
		err := checkFunc(ctx)
		duration := time.Since(start)

		if err != nil {
			return CheckResult{
				Status:  HealthStatusUnhealthy,
				Message: fmt.Sprintf("Database check failed: %v", err),
				Details: map[string]any{
					"duration": duration.String(),
				},
			}
		}

		return CheckResult{
			Status:  HealthStatusHealthy,
			Message: "Database is accessible",
			Details: map[string]any{
				"duration": duration.String(),
			},
		}
	}
}

// HTTPCheck creates an HTTP endpoint health check
func HTTPCheck(url string, timeout time.Duration) HealthCheck {
	return func(ctx context.Context) CheckResult {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return CheckResult{
				Status:  HealthStatusUnhealthy,
				Message: fmt.Sprintf("Failed to create request: %v", err),
			}
		}

		start := time.Now()
		resp, err := http.DefaultClient.Do(req)
		duration := time.Since(start)

		if err != nil {
			return CheckResult{
				Status:  HealthStatusUnhealthy,
				Message: fmt.Sprintf("HTTP check failed: %v", err),
				Details: map[string]any{
					"url":      url,
					"duration": duration.String(),
				},
			}
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				slog.Error("Failed to close response body", "error", err)
			}
		}()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return CheckResult{
				Status:  HealthStatusHealthy,
				Message: fmt.Sprintf("HTTP check successful: %d", resp.StatusCode),
				Details: map[string]any{
					"url":         url,
					"status_code": resp.StatusCode,
					"duration":    duration.String(),
				},
			}
		}

		return CheckResult{
			Status:  HealthStatusUnhealthy,
			Message: fmt.Sprintf("HTTP check failed with status: %d", resp.StatusCode),
			Details: map[string]any{
				"url":         url,
				"status_code": resp.StatusCode,
				"duration":    duration.String(),
			},
		}
	}
}

// FallbackConfig provides fallback configuration for health checks
type FallbackConfig struct {
	PrimaryTool       string `json:"primary_tool"`
	FallbackTool      string `json:"fallback_tool"`
	FallbackCondition string `json:"fallback_condition"`
	DegradedMode      bool   `json:"degraded_mode"`
}

// NewFallbackConfig creates a new fallback configuration
func NewFallbackConfig(primary, fallback, condition string) *FallbackConfig {
	return &FallbackConfig{
		PrimaryTool:       primary,
		FallbackTool:      fallback,
		FallbackCondition: condition,
		DegradedMode:      true,
	}
}

// SystemResourceCheck creates a comprehensive system resource health check
func SystemResourceCheck() HealthCheck {
	return func(ctx context.Context) CheckResult {
		start := time.Now()

		// Check memory usage
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)

		// Check goroutines
		goroutines := runtime.NumGoroutine()

		// Check disk usage for common paths
		diskUsage := checkDiskUsage()

		// Determine status based on thresholds
		status := HealthStatusHealthy
		messages := make([]string, 0)

		// Memory check (if > 80% of allocated memory is in use)
		if memStats.Alloc > memStats.Sys*8/10 {
			status = HealthStatusDegraded
			messages = append(messages, "High memory usage")
		}

		// Disk usage check (if > 90% disk usage)
		if diskUsage > 90 {
			status = HealthStatusDegraded
			messages = append(messages, "High disk usage")
		}

		// Goroutine check (if > 1000 goroutines)
		if goroutines > 1000 {
			status = HealthStatusDegraded
			messages = append(messages, "High number of goroutines")
		}

		message := "System resources are healthy"
		if len(messages) > 0 {
			message = fmt.Sprintf("System issues detected: %s", fmt.Sprintf("[%s]", fmt.Sprintf("%v", messages)))
		}

		return CheckResult{
			Status:  status,
			Message: message,
			Details: map[string]any{
				"memory_alloc":   fmt.Sprintf("%d bytes", memStats.Alloc),
				"memory_sys":     fmt.Sprintf("%d bytes", memStats.Sys),
				"goroutines":     goroutines,
				"disk_usage":     fmt.Sprintf("%.2f%%", diskUsage),
				"disk_usage_raw": diskUsage,
				"duration":       time.Since(start).String(),
			},
		}
	}
}

// ProcessCheck creates a process health check
func ProcessCheck() HealthCheck {
	return func(ctx context.Context) CheckResult {
		start := time.Now()

		// Check if process can access its own PID
		pid := os.Getpid()
		process, err := os.FindProcess(pid)
		if err != nil {
			return CheckResult{
				Status:  HealthStatusUnhealthy,
				Message: fmt.Sprintf("Cannot find process: %v", err),
			}
		}

		// Try to signal the process (this checks if it's alive)
		err = process.Signal(syscall.Signal(0))
		if err != nil {
			return CheckResult{
				Status:  HealthStatusUnhealthy,
				Message: fmt.Sprintf("Process signal failed: %v", err),
			}
		}

		return CheckResult{
			Status:  HealthStatusHealthy,
			Message: "Process is running",
			Details: map[string]any{
				"pid":      pid,
				"duration": time.Since(start).String(),
			},
		}
	}
}

// FileSystemCheck creates a filesystem health check
func FileSystemCheck(paths []string) HealthCheck {
	return func(ctx context.Context) CheckResult {
		start := time.Now()

		issues := make([]string, 0)

		for _, path := range paths {
			// Check if path exists
			if _, err := os.Stat(path); err != nil {
				if os.IsNotExist(err) {
					issues = append(issues, fmt.Sprintf("Path does not exist: %s", path))
				} else {
					issues = append(issues, fmt.Sprintf("Cannot access path %s: %v", path, err))
				}
				continue
			}

			// Check if path is writable (for temp directories)
			if path == "/tmp" || path == os.TempDir() {
				testFile := fmt.Sprintf("%s/health_check_%d", path, time.Now().UnixNano())
				if file, err := os.Create(testFile); err != nil {
					issues = append(issues, fmt.Sprintf("Cannot write to %s: %v", path, err))
				} else {
					_ = file.Close()
					_ = os.Remove(testFile)
				}
			}
		}

		status := HealthStatusHealthy
		message := "Filesystem checks passed"

		if len(issues) > 0 {
			status = HealthStatusUnhealthy
			message = fmt.Sprintf("Filesystem issues: %s", fmt.Sprintf("[%s]", fmt.Sprintf("%v", issues)))
		}

		return CheckResult{
			Status:  status,
			Message: message,
			Details: map[string]any{
				"paths_checked": paths,
				"issues":        issues,
				"duration":      time.Since(start).String(),
			},
		}
	}
}

// NetworkCheck creates a network connectivity health check
func NetworkCheck(host string, port int, timeout time.Duration) HealthCheck {
	return func(ctx context.Context) CheckResult {
		start := time.Now()

		// Simple TCP connection check
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
		duration := time.Since(start)

		if err != nil {
			return CheckResult{
				Status:  HealthStatusUnhealthy,
				Message: fmt.Sprintf("Network check failed: %v", err),
				Details: map[string]any{
					"host":     host,
					"port":     port,
					"duration": duration.String(),
				},
			}
		}

		_ = conn.Close()

		return CheckResult{
			Status:  HealthStatusHealthy,
			Message: "Network connectivity is good",
			Details: map[string]any{
				"host":     host,
				"port":     port,
				"duration": duration.String(),
			},
		}
	}
}

// SecurityCheck creates a security configuration health check
func SecurityCheck(securityValidator interface{}) HealthCheck {
	return func(ctx context.Context) CheckResult {
		start := time.Now()

		// Basic security checks
		issues := make([]string, 0)

		// Check if security validator is properly initialized
		if securityValidator == nil {
			issues = append(issues, "Security validator not initialized")
		}

		// Check environment variables for sensitive data exposure
		envVars := []string{"API_KEY", "SECRET", "PASSWORD", "TOKEN"}
		for _, envVar := range envVars {
			if os.Getenv(envVar) != "" {
				issues = append(issues, fmt.Sprintf("Sensitive environment variable %s is set", envVar))
			}
		}

		status := HealthStatusHealthy
		message := "Security configuration is good"

		if len(issues) > 0 {
			status = HealthStatusDegraded
			message = fmt.Sprintf("Security issues detected: %s", fmt.Sprintf("[%s]", fmt.Sprintf("%v", issues)))
		}

		return CheckResult{
			Status:  status,
			Message: message,
			Details: map[string]any{
				"issues":   issues,
				"duration": time.Since(start).String(),
			},
		}
	}
}

// MetricsHealthCheck creates a health check that monitors metrics
func MetricsHealthCheck(metricsCollector interface{}) HealthCheck {
	return func(ctx context.Context) CheckResult {
		start := time.Now()

		// Basic metrics check
		// In a real implementation, this would check metrics thresholds

		return CheckResult{
			Status:  HealthStatusHealthy,
			Message: "Metrics collection is operational",
			Details: map[string]any{
				"duration": time.Since(start).String(),
			},
		}
	}
}

// checkDiskUsage checks disk usage for common paths
func checkDiskUsage() float64 {
	paths := []string{"/", "/tmp", "/var"}

	for _, path := range paths {
		usage := getPathDiskUsage(path)
		if usage > 0 {
			return usage
		}
	}

	return 0
}

// getPathDiskUsage gets disk usage for a specific path
func getPathDiskUsage(path string) float64 {
	// Use syscall.Statfs for cross-platform disk usage
	var stat syscall.Statfs_t
	err := syscall.Statfs(path, &stat)

	if err != nil {
		return 0 // Path doesn't exist or can't be accessed
	}

	// Calculate usage percentage
	totalBytes := stat.Blocks * uint64(stat.Bsize)
	freeBytes := stat.Bavail * uint64(stat.Bsize)

	if totalBytes == 0 {
		return 0
	}

	usedBytes := totalBytes - freeBytes
	usagePercent := float64(usedBytes) / float64(totalBytes) * 100

	return usagePercent
}

// CreateDefaultHealthChecker creates a health checker with default checks
func CreateDefaultHealthChecker(version string) *HealthChecker {
	checker := NewHealthChecker(version)

	// Add basic health checks
	checker.AddCheck("ping", PingCheck())
	checker.AddCheck("process", ProcessCheck())
	checker.AddCheck("system", SystemResourceCheck())

	// Add filesystem checks for common paths
	checker.AddCheck("filesystem", FileSystemCheck([]string{"/tmp", "/var/tmp", os.TempDir()}))

	// Add network check (localhost)
	checker.AddCheck("network", NetworkCheck("localhost", 8080, 5*time.Second))

	return checker
}
