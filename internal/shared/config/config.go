package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"mini-mcp/internal/shared/auth"
	"mini-mcp/internal/shared/security"
	"mini-mcp/internal/shared/validation"
)

// Environment represents the deployment environment
type Environment string

const (
	EnvironmentDevelopment Environment = "development"
	EnvironmentStaging    Environment = "staging"
	EnvironmentProduction Environment = "production"
)

// Config holds the application configuration
type Config struct {
	Environment string `json:"environment"`
	LogLevel    string `json:"log_level"`
	Port        string `json:"port"`
	
	Security   SecurityConfig   `json:"security"`
	Auth       AuthConfig       `json:"auth"`
	Performance PerformanceConfig `json:"performance"`
	
	// Feature flags
	Features FeatureFlags `json:"features"`
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	// Command execution settings
	AllowedCommands []string `json:"allowed_commands"`
	WorkingDirectory string `json:"working_directory"`
	CommandTimeout time.Duration `json:"command_timeout"`
	MaxOutputSize int64 `json:"max_output_size"`
	
	// Path restrictions
	AllowedPaths []string `json:"allowed_paths"`
	BlockedPaths []string `json:"blocked_paths"`
	
	// Environment variables
	AllowedEnvVars []string `json:"allowed_env_vars"`
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	APIKeys     map[string]string `json:"api_keys"`
	RateLimiting time.Duration    `json:"rate_limiting"`
	IPWhitelist  []string         `json:"ip_whitelist"`
	MaxRequests  int              `json:"max_requests"`
	WindowSize   time.Duration    `json:"window_size"`
}

// PerformanceConfig holds performance-related configuration
type PerformanceConfig struct {
	MaxConcurrentRequests int           `json:"max_concurrent_requests"`
	RequestTimeout        time.Duration `json:"request_timeout"`
	IdleTimeout           time.Duration `json:"idle_timeout"`
	ReadTimeout           time.Duration `json:"read_timeout"`
	WriteTimeout          time.Duration `json:"write_timeout"`
	
	// Caching settings
	CacheEnabled bool          `json:"cache_enabled"`
	CacheTTL     time.Duration `json:"cache_ttl"`
	CacheSize    int           `json:"cache_size"`
}

// FeatureFlags holds feature flag configuration
type FeatureFlags struct {
	BatchOperations bool `json:"batch_operations"`
	AsyncOperations bool `json:"async_operations"`
	WorkflowSupport bool `json:"workflow_support"`
	StateManagement bool `json:"state_management"`
}

// LoadConfig loads configuration from environment variables, files, and defaults
func LoadConfig() (*Config, error) {
	// Load environment-specific defaults first
	config := loadEnvironmentDefaults()

	// Load configuration from file if specified
	if configFile := getEnv("CONFIG_FILE", ""); configFile != "" {
		if err := loadConfigFromFile(configFile, config); err != nil {
			return nil, fmt.Errorf("failed to load config file: %w", err)
		}
	}

	// Override with environment variables
	loadConfigFromEnv(config)

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// loadEnvironmentDefaults loads environment-specific default configurations
func loadEnvironmentDefaults() *Config {
	environment := Environment(getEnv("ENVIRONMENT", string(EnvironmentDevelopment)))

	config := &Config{
		Environment: string(environment),
		LogLevel:    getEnv("LOG_LEVEL", "INFO"),
		Port:        getEnv("PORT", ":8080"),
	}

	// Environment-specific defaults
	switch environment {
	case EnvironmentDevelopment:
		config.LogLevel = "DEBUG"
		config.Security = loadDevelopmentSecurityConfig()
		config.Auth = loadDevelopmentAuthConfig()
		config.Performance = loadDevelopmentPerformanceConfig()
	case EnvironmentStaging:
		config.LogLevel = "INFO"
		config.Security = loadStagingSecurityConfig()
		config.Auth = loadStagingAuthConfig()
		config.Performance = loadStagingPerformanceConfig()
	case EnvironmentProduction:
		config.LogLevel = "WARNING"
		config.Security = loadProductionSecurityConfig()
		config.Auth = loadProductionAuthConfig()
		config.Performance = loadProductionPerformanceConfig()
	default:
		// Fall back to basic configuration
		config.Security = loadSecurityConfig()
		config.Auth = loadAuthConfig()
		config.Performance = loadPerformanceConfig()
	}

	config.Features = loadFeatureFlags()

	return config
}

// loadSecurityConfig loads security configuration
func loadSecurityConfig() SecurityConfig {
	config := SecurityConfig{
		AllowedCommands: []string{
			"ls", "cat", "head", "tail", "grep", "find", "wc", "sort", "uniq",
			"ps", "top", "df", "du", "free", "uptime", "who", "w",
			"git", "docker", "nomad", "consul", "terraform",
		},
		WorkingDirectory: getEnv("SECURITY_WORKING_DIR", "/tmp"),
		CommandTimeout:   getDurationEnv("SECURITY_COMMAND_TIMEOUT", 30*time.Second),
		MaxOutputSize:    getInt64Env("SECURITY_MAX_OUTPUT_SIZE", 1024*1024), // 1MB
		AllowedPaths:     []string{"/tmp", "/var/log", "/proc"},
		BlockedPaths:     []string{"/etc/passwd", "/etc/shadow", "/root", "/home"},
		AllowedEnvVars:   []string{"PATH", "HOME", "USER", "PWD"},
	}
	
	// Override with environment variables if provided
	if allowedCommands := getEnv("SECURITY_ALLOWED_COMMANDS", ""); allowedCommands != "" {
		config.AllowedCommands = strings.Split(allowedCommands, ",")
	}
	
	if allowedPaths := getEnv("SECURITY_ALLOWED_PATHS", ""); allowedPaths != "" {
		config.AllowedPaths = strings.Split(allowedPaths, ",")
	}
	
	if blockedPaths := getEnv("SECURITY_BLOCKED_PATHS", ""); blockedPaths != "" {
		config.BlockedPaths = strings.Split(blockedPaths, ",")
	}
	
	if allowedEnvVars := getEnv("SECURITY_ALLOWED_ENV_VARS", ""); allowedEnvVars != "" {
		config.AllowedEnvVars = strings.Split(allowedEnvVars, ",")
	}
	
	return config
}

// loadAuthConfig loads authentication configuration
func loadAuthConfig() AuthConfig {
	config := AuthConfig{
		APIKeys:      make(map[string]string),
		RateLimiting: getDurationEnv("AUTH_RATE_LIMITING", 100*time.Millisecond),
		IPWhitelist:  []string{"127.0.0.1", "::1"},
		MaxRequests:  getIntEnv("AUTH_MAX_REQUESTS", 1000),
		WindowSize:   getDurationEnv("AUTH_WINDOW_SIZE", 1*time.Hour),
	}
	
	// Load API keys from environment
	if apiKeys := getEnv("AUTH_API_KEYS", ""); apiKeys != "" {
		pairs := strings.Split(apiKeys, ",")
		for _, pair := range pairs {
			parts := strings.SplitN(pair, ":", 2)
			if len(parts) == 2 {
				config.APIKeys[parts[0]] = parts[1]
			}
		}
	}
	
	// Load IP whitelist
	if ipWhitelist := getEnv("AUTH_IP_WHITELIST", ""); ipWhitelist != "" {
		config.IPWhitelist = strings.Split(ipWhitelist, ",")
	}
	
	return config
}

// loadPerformanceConfig loads performance configuration
func loadPerformanceConfig() PerformanceConfig {
	return PerformanceConfig{
		MaxConcurrentRequests: getIntEnv("PERF_MAX_CONCURRENT_REQUESTS", 100),
		RequestTimeout:        getDurationEnv("PERF_REQUEST_TIMEOUT", 30*time.Second),
		IdleTimeout:           getDurationEnv("PERF_IDLE_TIMEOUT", 120*time.Second),
		ReadTimeout:           getDurationEnv("PERF_READ_TIMEOUT", 30*time.Second),
		WriteTimeout:          getDurationEnv("PERF_WRITE_TIMEOUT", 30*time.Second),
		CacheEnabled:          getBoolEnv("PERF_CACHE_ENABLED", false),
		CacheTTL:              getDurationEnv("PERF_CACHE_TTL", 5*time.Minute),
		CacheSize:             getIntEnv("PERF_CACHE_SIZE", 1000),
	}
}

// loadFeatureFlags loads feature flag configuration
func loadFeatureFlags() FeatureFlags {
	return FeatureFlags{
		BatchOperations: getBoolEnv("FEATURE_BATCH_OPERATIONS", false),
		AsyncOperations: getBoolEnv("FEATURE_ASYNC_OPERATIONS", false),
		WorkflowSupport: getBoolEnv("FEATURE_WORKFLOW_SUPPORT", false),
		StateManagement: getBoolEnv("FEATURE_STATE_MANAGEMENT", false),
	}
}

// ToSecurityConfig converts the security configuration to the security package format
func (c *Config) ToSecurityConfig() *security.SecurityConfig {
	return &security.SecurityConfig{
		AllowedCommands: c.Security.AllowedCommands,
		WorkingDirectory: c.Security.WorkingDirectory,
		CommandTimeout: c.Security.CommandTimeout,
		MaxOutputSize: c.Security.MaxOutputSize,
		AllowedEnvVars: c.Security.AllowedEnvVars,
		AllowedPaths: c.Security.AllowedPaths,
		BlockedPaths: c.Security.BlockedPaths,
	}
}

// ToAuthConfig converts the auth configuration to the auth package format
func (c *Config) ToAuthConfig() *auth.AuthConfig {
	return &auth.AuthConfig{
		APIKeys: c.Auth.APIKeys,
		RateLimiting: c.Auth.RateLimiting,
		IPWhitelist: c.Auth.IPWhitelist,
		MaxRequests: c.Auth.MaxRequests,
		WindowSize: c.Auth.WindowSize,
	}
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.Environment == string(EnvironmentDevelopment)
}

// IsStaging returns true if the environment is staging
func (c *Config) IsStaging() bool {
	return c.Environment == string(EnvironmentStaging)
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	return c.Environment == string(EnvironmentProduction)
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate port
	if err := validation.StringRequired("port", c.Port); err != nil {
		return err
	}
	
	// Validate log level
	logLevelValidator := validation.EnumString("DEBUG", "INFO", "WARNING", "ERROR")
	if err := logLevelValidator("log_level", c.LogLevel); err != nil {
		return err
	}

	// Validate environment
	envValidator := validation.EnumString("development", "staging", "production")
	if err := envValidator("environment", c.Environment); err != nil {
		return err
	}
	
	// Validate command timeout
	vf := validation.NewValidationFactory()
	if err := vf.DurationPositive("command_timeout", c.Security.CommandTimeout); err != nil {
		return err
	}

	// Validate max output size
	sizeValidator := vf.RangeInt64(1024, 104857600)
	if err := sizeValidator("max_output_size", c.Security.MaxOutputSize); err != nil {
		return err
	}

	// Validate max requests
	if err := vf.Positive("max_requests", c.Auth.MaxRequests); err != nil {
		return err
	}

	// Validate window size
	if err := vf.DurationPositive("window_size", c.Auth.WindowSize); err != nil {
		return err
	}
	
	return nil
}

// ToJSON converts the configuration to JSON
func (c *Config) ToJSON() ([]byte, error) {
	return json.MarshalIndent(c, "", "  ")
}

// Helper functions for environment variable parsing

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getInt64Env(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// loadConfigFromFile loads configuration from a JSON file
func loadConfigFromFile(filename string, config *Config) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	
	if err := json.Unmarshal(data, config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}
	
	return nil
}

// loadConfigFromEnv loads configuration from environment variables
func loadConfigFromEnv(config *Config) {
	// Override basic settings
	if env := getEnv("ENVIRONMENT", ""); env != "" {
		config.Environment = env
	}
	if logLevel := getEnv("LOG_LEVEL", ""); logLevel != "" {
		config.LogLevel = logLevel
	}
	if port := getEnv("PORT", ""); port != "" {
		config.Port = port
	}
	
	// Override security settings
	if allowedCommands := getEnv("SECURITY_ALLOWED_COMMANDS", ""); allowedCommands != "" {
		config.Security.AllowedCommands = strings.Split(allowedCommands, ",")
	}
	if workingDir := getEnv("SECURITY_WORKING_DIR", ""); workingDir != "" {
		config.Security.WorkingDirectory = workingDir
	}
	if timeout := getEnv("SECURITY_COMMAND_TIMEOUT", ""); timeout != "" {
		if duration, err := time.ParseDuration(timeout); err == nil {
			config.Security.CommandTimeout = duration
		}
	}
	if maxOutput := getEnv("SECURITY_MAX_OUTPUT_SIZE", ""); maxOutput != "" {
		if size, err := strconv.ParseInt(maxOutput, 10, 64); err == nil {
			config.Security.MaxOutputSize = size
		}
	}
	
	// Override auth settings
	if rateLimit := getEnv("AUTH_RATE_LIMITING", ""); rateLimit != "" {
		if duration, err := time.ParseDuration(rateLimit); err == nil {
			config.Auth.RateLimiting = duration
		}
	}
	if maxRequests := getEnv("AUTH_MAX_REQUESTS", ""); maxRequests != "" {
		if requests, err := strconv.Atoi(maxRequests); err == nil {
			config.Auth.MaxRequests = requests
		}
	}
	if windowSize := getEnv("AUTH_WINDOW_SIZE", ""); windowSize != "" {
		if duration, err := time.ParseDuration(windowSize); err == nil {
			config.Auth.WindowSize = duration
		}
	}
	
	// Override performance settings
	if maxConcurrent := getEnv("PERF_MAX_CONCURRENT_REQUESTS", ""); maxConcurrent != "" {
		if concurrent, err := strconv.Atoi(maxConcurrent); err == nil {
			config.Performance.MaxConcurrentRequests = concurrent
		}
	}
	if requestTimeout := getEnv("PERF_REQUEST_TIMEOUT", ""); requestTimeout != "" {
		if duration, err := time.ParseDuration(requestTimeout); err == nil {
			config.Performance.RequestTimeout = duration
		}
	}
	if cacheEnabled := getEnv("PERF_CACHE_ENABLED", ""); cacheEnabled != "" {
		if enabled, err := strconv.ParseBool(cacheEnabled); err == nil {
			config.Performance.CacheEnabled = enabled
		}
	}
}

// Environment-specific configuration loaders

func loadDevelopmentSecurityConfig() SecurityConfig {
	config := loadSecurityConfig()
	config.AllowedCommands = append(config.AllowedCommands, "debug", "trace", "profile")
	config.CommandTimeout = 60 * time.Second
	config.MaxOutputSize = 10 * 1024 * 1024 // 10MB for development
	return config
}

func loadDevelopmentAuthConfig() AuthConfig {
	config := loadAuthConfig()
	config.RateLimiting = 10 * time.Millisecond // More permissive for development
	config.MaxRequests = 10000
	return config
}

func loadDevelopmentPerformanceConfig() PerformanceConfig {
	config := loadPerformanceConfig()
	config.MaxConcurrentRequests = 50
	config.RequestTimeout = 60 * time.Second
	config.CacheEnabled = false // Disable cache in development for easier debugging
	return config
}

func loadStagingSecurityConfig() SecurityConfig {
	config := loadSecurityConfig()
	config.CommandTimeout = 30 * time.Second
	config.MaxOutputSize = 5 * 1024 * 1024 // 5MB for staging
	return config
}

func loadStagingAuthConfig() AuthConfig {
	config := loadAuthConfig()
	config.RateLimiting = 50 * time.Millisecond
	config.MaxRequests = 5000
	return config
}

func loadStagingPerformanceConfig() PerformanceConfig {
	config := loadPerformanceConfig()
	config.MaxConcurrentRequests = 75
	config.RequestTimeout = 45 * time.Second
	config.CacheEnabled = true
	config.CacheTTL = 2 * time.Minute
	return config
}

func loadProductionSecurityConfig() SecurityConfig {
	config := loadSecurityConfig()
	config.CommandTimeout = 15 * time.Second
	config.MaxOutputSize = 1024 * 1024 // 1MB for production
	// Remove potentially dangerous commands in production
	var safeCommands []string
	for _, cmd := range config.AllowedCommands {
		if cmd != "debug" && cmd != "trace" && cmd != "profile" {
			safeCommands = append(safeCommands, cmd)
		}
	}
	config.AllowedCommands = safeCommands
	return config
}

func loadProductionAuthConfig() AuthConfig {
	config := loadAuthConfig()
	config.RateLimiting = 100 * time.Millisecond
	config.MaxRequests = 1000
	return config
}

func loadProductionPerformanceConfig() PerformanceConfig {
	config := loadPerformanceConfig()
	config.MaxConcurrentRequests = 100
	config.RequestTimeout = 30 * time.Second
	config.CacheEnabled = true
	config.CacheTTL = 5 * time.Minute
	config.CacheSize = 5000
	return config
}
