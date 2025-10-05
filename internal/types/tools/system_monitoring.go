// Package tools contains types for MCP tool arguments and responses.
package tools

import (
	"mini-mcp/internal/shared/validation"
)

// SystemMetric represents the type of system metric to retrieve
type SystemMetric string

const (
	MetricProcesses   SystemMetric = "processes"
	MetricDiskUsage   SystemMetric = "disk_usage"
	MetricMemoryUsage SystemMetric = "memory_usage"
	MetricNetwork     SystemMetric = "network"
	MetricUptime      SystemMetric = "uptime"
)

// SystemMonitoringArgs represents arguments for the system_monitoring tool.
// This tool provides comprehensive system monitoring and performance analysis capabilities.
//
// Monitoring Capabilities:
// - processes: Detailed process information with CPU/memory usage, status, and command details
// - disk_usage: Filesystem usage statistics with capacity, used space, and availability
// - memory_usage: RAM consumption with total, used, free, and percentage calculations
// - network: Network interface information including IP addresses, status, and configuration
// - uptime: System uptime, load averages, boot time, and user session information
//
// Data Format:
// - All metrics return structured JSON with detailed information
// - Numeric values include both raw and formatted representations
// - Timestamps are provided in ISO 8601 format
// - Percentages and ratios are calculated automatically
//
// Performance Considerations:
// - Process listing is limited to top 20 processes for performance
// - Disk usage includes all mounted filesystems
// - Memory calculations include swap and cache information
// - Network data includes all active interfaces
//
// Examples:
//
//	{"metric": "processes"}     // Get running processes
//	{"metric": "disk_usage"}    // Check filesystem usage
//	{"metric": "memory_usage"}  // Monitor RAM consumption
//	{"metric": "network"}       // List network interfaces
//	{"metric": "uptime"}        // Get system uptime and load
type SystemMonitoringArgs struct {
	// Metric is the system metric to retrieve
	// Must be one of: processes, disk_usage, memory_usage, network, uptime
	Metric SystemMetric `json:"metric"`
}

// Validate checks if the system monitoring arguments are valid.
func (args *SystemMonitoringArgs) Validate() error {
	metric := string(args.Metric)
	if metric == "" {
		return validation.NewMissingRequiredError("metric")
	}
	if metric != "processes" && metric != "disk_usage" && metric != "memory_usage" && 
	   metric != "network" && metric != "uptime" {
		return validation.NewInvalidFormatError("metric", "must be one of: processes, disk_usage, memory_usage, network, uptime")
	}
	
	return nil
}

// NewSystemMonitoringArgs creates a new SystemMonitoringArgs with the given metric.
func NewSystemMonitoringArgs(metric SystemMetric) *SystemMonitoringArgs {
	return &SystemMonitoringArgs{
		Metric: metric,
	}
}


