package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// LogLevel represents the logging level
type LogLevel string

// Logger interface defines the contract for logging operations
type Logger interface {
	Log(level LogLevel, message string, fields map[string]any)
	Debug(message string, fields map[string]any)
	Info(message string, fields map[string]any)
	Warning(message string, fields map[string]any)
	Error(message string, err error, fields map[string]any)
	Fatal(message string, fields map[string]any)
	GetMetrics() *Metrics
}

const (
	LogLevelDebug   LogLevel = "DEBUG"
	LogLevelInfo    LogLevel = "INFO"
	LogLevelWarning LogLevel = "WARNING"
	LogLevelError   LogLevel = "ERROR"
	LogLevelFatal   LogLevel = "FATAL"
)

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp   time.Time         `json:"timestamp"`
	Level       LogLevel          `json:"level"`
	Tool        string            `json:"tool,omitempty"`
	Duration    time.Duration     `json:"duration,omitempty"`
	Success     bool              `json:"success,omitempty"`
	UserAgent   string            `json:"user_agent,omitempty"`
	RequestID   string            `json:"request_id,omitempty"`
	Message     string            `json:"message"`
	Fields      map[string]any    `json:"fields,omitempty"`
	Error       string            `json:"error,omitempty"`
	Stack       string            `json:"stack,omitempty"`
}

// LoggerImpl provides structured logging functionality
type LoggerImpl struct {
	output   io.Writer
	level    LogLevel
	mu       sync.Mutex
	metrics  *Metrics
}

// NewLogger creates a new logger
func NewLogger(output io.Writer, level LogLevel) Logger {
	if output == nil {
		output = os.Stderr
	}
	
	return &LoggerImpl{
		output:  output,
		level:   level,
		metrics: NewMetrics(),
	}
}

// Log logs a message with the given level and fields
func (l *LoggerImpl) Log(level LogLevel, message string, fields map[string]any) {
	if !l.shouldLog(level) {
		return
	}
	
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Fields:    fields,
	}
	
	l.writeEntry(entry)
}

// Debug logs a debug message
func (l *LoggerImpl) Debug(message string, fields map[string]any) {
	l.Log(LogLevelDebug, message, fields)
}

// Info logs an info message
func (l *LoggerImpl) Info(message string, fields map[string]any) {
	l.Log(LogLevelInfo, message, fields)
}

// Warning logs a warning message
func (l *LoggerImpl) Warning(message string, fields map[string]any) {
	l.Log(LogLevelWarning, message, fields)
}

// Error logs an error message
func (l *LoggerImpl) Error(message string, err error, fields map[string]any) {
	if fields == nil {
		fields = make(map[string]any)
	}
	
	if err != nil {
		fields["error"] = err.Error()
	}
	
	l.Log(LogLevelError, message, fields)
}

// Fatal logs a fatal message and exits
func (l *LoggerImpl) Fatal(message string, fields map[string]any) {
	l.Log(LogLevelFatal, message, fields)
	os.Exit(1)
}

// shouldLog determines if a message should be logged based on the current level
func (l *LoggerImpl) shouldLog(level LogLevel) bool {
	levels := map[LogLevel]int{
		LogLevelDebug:   0,
		LogLevelInfo:    1,
		LogLevelWarning: 2,
		LogLevelError:   3,
		LogLevelFatal:   4,
	}
	
	return levels[level] >= levels[l.level]
}

// writeEntry writes a log entry to the output
func (l *LoggerImpl) writeEntry(entry LogEntry) {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	// Update metrics
	l.metrics.IncrementLogCount(entry.Level, entry.Tool)
	
	// Marshal to JSON
	data, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal log entry: %v\n", err)
		return
	}
	
	// Write to output
	if _, err := fmt.Fprintln(l.output, string(data)); err != nil {
		// Log error but don't fail the operation
		GetGlobalLogger().Error("Failed to write log output", err, map[string]any{
			"output_type": "log_entry",
		})
	}
}

// WithContext creates a logger with context information
func (l *LoggerImpl) WithContext(ctx context.Context) *ContextLogger {
	return &ContextLogger{
		logger: l,
		ctx:    ctx,
	}
}

// WithTool creates a logger with tool information
func (l *LoggerImpl) WithTool(tool string) *ToolLogger {
	return &ToolLogger{
		logger: l,
		tool:   tool,
	}
}

// GetMetrics returns the current metrics
func (l *LoggerImpl) GetMetrics() *Metrics {
	return l.metrics
}

// ContextLogger provides logging with context information
type ContextLogger struct {
	logger Logger
	ctx    context.Context
}

// Log logs a message with context information
func (l *ContextLogger) Log(level LogLevel, message string, fields map[string]any) {
	if fields == nil {
		fields = make(map[string]any)
	}
	
	// Extract context information
	if requestID, ok := l.ctx.Value("request_id").(string); ok {
		fields["request_id"] = requestID
	}
	
	if userAgent, ok := l.ctx.Value("user_agent").(string); ok {
		fields["user_agent"] = userAgent
	}
	
	l.logger.Log(level, message, fields)
}

// Debug logs a debug message with context
func (l *ContextLogger) Debug(message string, fields map[string]any) {
	l.Log(LogLevelDebug, message, fields)
}

// Info logs an info message with context
func (l *ContextLogger) Info(message string, fields map[string]any) {
	l.Log(LogLevelInfo, message, fields)
}

// Warning logs a warning message with context
func (l *ContextLogger) Warning(message string, fields map[string]any) {
	l.Log(LogLevelWarning, message, fields)
}

// Error logs an error message with context
func (l *ContextLogger) Error(message string, err error, fields map[string]any) {
	l.logger.Error(message, err, fields)
}

// NewToolLogger creates a new tool logger
func NewToolLogger(logger Logger) *ToolLogger {
	return &ToolLogger{
		logger: logger,
		tool:   "mcp",
	}
}

// ToolLogger provides logging with tool information
type ToolLogger struct {
	logger Logger
	tool   string
}

// Log logs a message with tool information
func (l *ToolLogger) Log(level LogLevel, message string, fields map[string]any) {
	if fields == nil {
		fields = make(map[string]any)
	}
	
	fields["tool"] = l.tool
	l.logger.Log(level, message, fields)
}

// Debug logs a debug message with tool context
func (l *ToolLogger) Debug(message string, fields map[string]any) {
	l.Log(LogLevelDebug, message, fields)
}

// Info logs an info message with tool context
func (l *ToolLogger) Info(message string, fields map[string]any) {
	l.Log(LogLevelInfo, message, fields)
}

// Warning logs a warning message with tool context
func (l *ToolLogger) Warning(message string, fields map[string]any) {
	l.Log(LogLevelWarning, message, fields)
}

// Error logs an error message with tool context
func (l *ToolLogger) Error(message string, err error, fields map[string]any) {
	if fields == nil {
		fields = make(map[string]any)
	}
	
	fields["tool"] = l.tool
	l.logger.Error(message, err, fields)
}

// Metrics provides metrics collection for logging
type Metrics struct {
	mu sync.RWMutex

	// Log counts by level and tool
	LogCounts map[string]int64 `json:"log_counts"`

	// Response times by tool
	ResponseTimes map[string][]time.Duration `json:"response_times"`

	// Error rates by tool
	ErrorRates map[string]float64 `json:"error_rates"`

	// Active connections
	ActiveConnections int `json:"active_connections"`

	// Request counts by tool
	RequestCounts map[string]int64 `json:"request_counts"`

	// Performance metrics
	PerformanceMetrics map[string]*PerformanceMetrics `json:"performance_metrics"`
}

// PerformanceMetrics tracks performance data for a specific tool
type PerformanceMetrics struct {
	TotalRequests     int64         `json:"total_requests"`
	TotalDuration     time.Duration `json:"total_duration"`
	AverageDuration   time.Duration `json:"average_duration"`
	MinDuration       time.Duration `json:"min_duration"`
	MaxDuration       time.Duration `json:"max_duration"`
	P95Duration       time.Duration `json:"p95_duration"`
	P99Duration       time.Duration `json:"p99_duration"`
	RequestsPerSecond float64       `json:"requests_per_second"`
	ErrorCount        int64         `json:"error_count"`
	ErrorRate         float64       `json:"error_rate"`
	LastUpdated       time.Time     `json:"last_updated"`
}

// NewMetrics creates new metrics
func NewMetrics() *Metrics {
	return &Metrics{
		LogCounts:          make(map[string]int64),
		ResponseTimes:      make(map[string][]time.Duration),
		ErrorRates:         make(map[string]float64),
		RequestCounts:      make(map[string]int64),
		PerformanceMetrics: make(map[string]*PerformanceMetrics),
	}
}

// IncrementLogCount increments the log count for a level and tool
func (m *Metrics) IncrementLogCount(level LogLevel, tool string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	key := fmt.Sprintf("%s:%s", level, tool)
	m.LogCounts[key]++
}

// RecordResponseTime records a response time for a tool
func (m *Metrics) RecordResponseTime(tool string, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if m.ResponseTimes[tool] == nil {
		m.ResponseTimes[tool] = make([]time.Duration, 0)
	}
	
	m.ResponseTimes[tool] = append(m.ResponseTimes[tool], duration)
	
	// Keep only the last 1000 measurements
	if len(m.ResponseTimes[tool]) > 1000 {
		m.ResponseTimes[tool] = m.ResponseTimes[tool][len(m.ResponseTimes[tool])-1000:]
	}
}

// UpdateErrorRate updates the error rate for a tool
func (m *Metrics) UpdateErrorRate(tool string, rate float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.ErrorRates[tool] = rate
}

// SetActiveConnections sets the number of active connections
func (m *Metrics) SetActiveConnections(count int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.ActiveConnections = count
}

// GetAverageResponseTime gets the average response time for a tool
func (m *Metrics) GetAverageResponseTime(tool string) time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	times, exists := m.ResponseTimes[tool]
	if !exists || len(times) == 0 {
		return 0
	}
	
	var total time.Duration
	for _, t := range times {
		total += t
	}
	
	return total / time.Duration(len(times))
}

// RecordRequest records a request for a tool
func (m *Metrics) RecordRequest(tool string, duration time.Duration, success bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Increment request count
	m.RequestCounts[tool]++

	// Update response times
	if m.ResponseTimes[tool] == nil {
		m.ResponseTimes[tool] = make([]time.Duration, 0)
	}
	m.ResponseTimes[tool] = append(m.ResponseTimes[tool], duration)

	// Keep only the last 1000 measurements
	if len(m.ResponseTimes[tool]) > 1000 {
		m.ResponseTimes[tool] = m.ResponseTimes[tool][len(m.ResponseTimes[tool])-1000:]
	}

	// Update performance metrics
	m.updatePerformanceMetrics(tool, duration, success)
}

// RecordError records an error for a tool
func (m *Metrics) RecordError(tool string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Update error rates
	totalRequests := m.RequestCounts[tool]
	errorCount := float64(0)
	if perfMetrics, exists := m.PerformanceMetrics[tool]; exists {
		perfMetrics.ErrorCount++
		errorCount = float64(perfMetrics.ErrorCount)
	}
	if totalRequests > 0 {
		m.ErrorRates[tool] = errorCount / float64(totalRequests)
	}
}

// updatePerformanceMetrics updates performance metrics for a tool
func (m *Metrics) updatePerformanceMetrics(tool string, duration time.Duration, success bool) {
	if m.PerformanceMetrics[tool] == nil {
		m.PerformanceMetrics[tool] = &PerformanceMetrics{
			MinDuration: duration,
			MaxDuration: duration,
		}
	}

	perf := m.PerformanceMetrics[tool]
	perf.TotalRequests++
	perf.TotalDuration += duration
	perf.AverageDuration = perf.TotalDuration / time.Duration(perf.TotalRequests)

	// Update min/max
	if duration < perf.MinDuration {
		perf.MinDuration = duration
	}
	if duration > perf.MaxDuration {
		perf.MaxDuration = duration
	}

	// Update percentiles (simplified calculation)
	responseTimes := m.ResponseTimes[tool]
	if len(responseTimes) > 0 {
		sortedTimes := make([]time.Duration, len(responseTimes))
		copy(sortedTimes, responseTimes)

		// Simple percentile calculation
		p95Index := int(float64(len(sortedTimes)) * 0.95)
		p99Index := int(float64(len(sortedTimes)) * 0.99)

		if p95Index >= len(sortedTimes) {
			p95Index = len(sortedTimes) - 1
		}
		if p99Index >= len(sortedTimes) {
			p99Index = len(sortedTimes) - 1
		}

		perf.P95Duration = sortedTimes[p95Index]
		perf.P99Duration = sortedTimes[p99Index]
	}

	// Update error count and rate
	if !success {
		perf.ErrorCount++
	}
	totalReqs := float64(perf.TotalRequests)
	perf.ErrorRate = float64(perf.ErrorCount) / totalReqs

	// Calculate requests per second (over last minute)
	elapsed := time.Since(perf.LastUpdated)
	if elapsed > time.Minute {
		perf.RequestsPerSecond = float64(perf.TotalRequests) / elapsed.Seconds()
		perf.LastUpdated = time.Now()
	}
}

// IncrementActiveConnections increments the active connections count
func (m *Metrics) IncrementActiveConnections() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ActiveConnections++
}

// DecrementActiveConnections decrements the active connections count
func (m *Metrics) DecrementActiveConnections() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.ActiveConnections > 0 {
		m.ActiveConnections--
	}
}

// GetMetricsSummary returns a summary of all metrics
func (m *Metrics) GetMetricsSummary() map[string]any {
	m.mu.RLock()
	defer m.mu.RUnlock()

	summary := make(map[string]any)

	// Log counts
	summary["log_counts"] = m.LogCounts

	// Request counts
	summary["request_counts"] = m.RequestCounts

	// Average response times
	avgResponseTimes := make(map[string]time.Duration)
	for tool := range m.ResponseTimes {
		avgResponseTimes[tool] = m.GetAverageResponseTime(tool)
	}
	summary["average_response_times"] = avgResponseTimes

	// Error rates
	summary["error_rates"] = m.ErrorRates

	// Active connections
	summary["active_connections"] = m.ActiveConnections

	// Performance metrics
	summary["performance_metrics"] = m.PerformanceMetrics

	return summary
}

// GetToolMetrics returns metrics for a specific tool
func (m *Metrics) GetToolMetrics(tool string) map[string]any {
	m.mu.RLock()
	defer m.mu.RUnlock()

	metrics := make(map[string]any)

	// Request count
	metrics["requests"] = m.RequestCounts[tool]

	// Average response time
	metrics["avg_response_time"] = m.GetAverageResponseTime(tool)

	// Error rate
	metrics["error_rate"] = m.ErrorRates[tool]

	// Performance metrics
	if perf, exists := m.PerformanceMetrics[tool]; exists {
		metrics["performance"] = perf
	}

	return metrics
}

// ResetMetrics resets all metrics (useful for testing)
func (m *Metrics) ResetMetrics() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.LogCounts = make(map[string]int64)
	m.ResponseTimes = make(map[string][]time.Duration)
	m.ErrorRates = make(map[string]float64)
	m.RequestCounts = make(map[string]int64)
	m.PerformanceMetrics = make(map[string]*PerformanceMetrics)
	m.ActiveConnections = 0
}

// Global logger instance
var globalLogger Logger

// InitGlobalLogger initializes the global logger
func InitGlobalLogger(level LogLevel) {
	globalLogger = NewLogger(os.Stderr, level)
	// Test log to ensure logger is working
	globalLogger.Info("Global logger initialized", map[string]any{
		"level": string(level),
		"output": "stderr",
	})
}

// GetGlobalLogger returns the global logger
func GetGlobalLogger() Logger {
	if globalLogger == nil {
		InitGlobalLogger(LogLevelInfo)
	}
	return globalLogger
}

// Debug logs a debug message using the global logger
func Debug(message string, fields map[string]any) {
	GetGlobalLogger().Debug(message, fields)
}

// Info logs an info message using the global logger
func Info(message string, fields map[string]any) {
	GetGlobalLogger().Info(message, fields)
}

// Warning logs a warning message using the global logger
func Warning(message string, fields map[string]any) {
	GetGlobalLogger().Warning(message, fields)
}

// Error logs an error message using the global logger
func Error(message string, err error, fields map[string]any) {
	GetGlobalLogger().Error(message, err, fields)
}

// Fatal logs a fatal message using the global logger
func Fatal(message string, fields map[string]any) {
	GetGlobalLogger().Fatal(message, fields)
}
