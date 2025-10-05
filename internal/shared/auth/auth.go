package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// contextKey is a type for context keys to avoid collisions
type contextKey string

// AuthConfig holds authentication configuration
type AuthConfig struct {
	APIKeys     map[string]string `json:"api_keys"`
	RateLimiting time.Duration    `json:"rate_limiting"`
	IPWhitelist  []string         `json:"ip_whitelist"`
	MaxRequests  int              `json:"max_requests"`
	WindowSize   time.Duration    `json:"window_size"`
}

// DefaultAuthConfig returns a secure default authentication configuration
func DefaultAuthConfig() *AuthConfig {
	return &AuthConfig{
		APIKeys:      make(map[string]string),
		RateLimiting: 100 * time.Millisecond, // Minimum time between requests
		IPWhitelist:  []string{"127.0.0.1", "::1"}, // Localhost only by default
		MaxRequests:  1000, // Max requests per window
		WindowSize:   1 * time.Hour, // 1 hour window
	}
}

// Authenticator provides authentication and authorization services
type Authenticator struct {
	config      *AuthConfig
	rateLimiter *RateLimiter
	mu          sync.RWMutex
	sessions    map[string]time.Time
}

// NewAuthenticator creates a new authenticator
func NewAuthenticator(config *AuthConfig) *Authenticator {
	if config == nil {
		config = DefaultAuthConfig()
	}
	return &Authenticator{
		config:      config,
		rateLimiter: NewRateLimiter(config.MaxRequests, config.WindowSize),
		sessions:    make(map[string]time.Time),
	}
}

// Cleanup performs cleanup of expired sessions and resources
func (a *Authenticator) Cleanup() {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Clean up expired sessions
	now := time.Now()
	for id, lastAccess := range a.sessions {
		if now.Sub(lastAccess) > a.config.WindowSize {
			delete(a.sessions, id)
		}
	}

	// Reset rate limiter if exists
	if a.rateLimiter != nil {
		a.rateLimiter.Reset()
	}
}

// AuthenticateRequest authenticates an HTTP request
func (a *Authenticator) AuthenticateRequest(r *http.Request) (*AuthResult, error) {
	// Check IP whitelist
	if err := a.checkIPWhitelist(r); err != nil {
		return nil, fmt.Errorf("IP not allowed: %w", err)
	}
	
	// Extract API key
	apiKey := a.extractAPIKey(r)
	if apiKey == "" {
		return nil, fmt.Errorf("missing API key")
	}
	
	// Validate API key
	userID, err := a.validateAPIKey(apiKey)
	if err != nil {
		return nil, fmt.Errorf("invalid API key: %w", err)
	}
	
	// Check rate limiting
	if err := a.rateLimiter.CheckLimit(userID); err != nil {
		return nil, fmt.Errorf("rate limit exceeded: %w", err)
	}
	
	return &AuthResult{
		UserID:    userID,
		APIKey:    apiKey,
		IPAddress: a.getClientIP(r),
		Timestamp: time.Now(),
	}, nil
}

// AuthResult contains authentication result information
type AuthResult struct {
	UserID    string    `json:"user_id"`
	APIKey    string    `json:"api_key"`
	IPAddress string    `json:"ip_address"`
	Timestamp time.Time `json:"timestamp"`
}

// checkIPWhitelist checks if the client IP is in the whitelist
func (a *Authenticator) checkIPWhitelist(r *http.Request) error {
	clientIP := a.getClientIP(r)
	
	// If no whitelist is configured, allow all
	if len(a.config.IPWhitelist) == 0 {
		return nil
	}
	
	for _, allowedIP := range a.config.IPWhitelist {
		if clientIP == allowedIP {
			return nil
		}
	}
	
	return fmt.Errorf("IP %s not in whitelist", clientIP)
}

// getClientIP extracts the real client IP from the request
func (a *Authenticator) getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		ips := strings.Split(forwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	
	// Check X-Real-IP header
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}
	
	// Fall back to remote address
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	
	return host
}

// extractAPIKey extracts the API key from the request
func (a *Authenticator) extractAPIKey(r *http.Request) string {
	// Check Authorization header
	if auth := r.Header.Get("Authorization"); auth != "" {
		if strings.HasPrefix(auth, "Bearer ") {
			return strings.TrimPrefix(auth, "Bearer ")
		}
		if strings.HasPrefix(auth, "ApiKey ") {
			return strings.TrimPrefix(auth, "ApiKey ")
		}
	}
	
	// Check X-API-Key header
	if apiKey := r.Header.Get("X-API-Key"); apiKey != "" {
		return apiKey
	}
	
	// Check query parameter
	if apiKey := r.URL.Query().Get("api_key"); apiKey != "" {
		return apiKey
	}
	
	return ""
}

// validateAPIKey validates the provided API key
func (a *Authenticator) validateAPIKey(apiKey string) (string, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	for userID, key := range a.config.APIKeys {
		if key == apiKey {
			return userID, nil
		}
	}
	
	return "", fmt.Errorf("invalid API key")
}

// AddAPIKey adds a new API key for a user
func (a *Authenticator) AddAPIKey(userID string) (string, error) {
	// Generate a secure random API key
	apiKey, err := a.generateAPIKey()
	if err != nil {
		return "", fmt.Errorf("failed to generate API key: %w", err)
	}
	
	a.mu.Lock()
	defer a.mu.Unlock()
	
	a.config.APIKeys[userID] = apiKey
	return apiKey, nil
}

// RemoveAPIKey removes an API key for a user
func (a *Authenticator) RemoveAPIKey(userID string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	delete(a.config.APIKeys, userID)
	return nil
}

// generateAPIKey generates a secure random API key
func (a *Authenticator) generateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	
	// Create a hash for better formatting
	hash := sha256.Sum256(bytes)
	return hex.EncodeToString(hash[:]), nil
}

// RateLimiter provides rate limiting functionality
type RateLimiter struct {
	maxRequests int
	windowSize  time.Duration
	requests    map[string][]time.Time
	mu          sync.RWMutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(maxRequests int, windowSize time.Duration) *RateLimiter {
	return &RateLimiter{
		maxRequests: maxRequests,
		windowSize:  windowSize,
		requests:    make(map[string][]time.Time),
	}
}

// CheckLimit checks if a user has exceeded their rate limit
func (r *RateLimiter) CheckLimit(userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	now := time.Now()
	windowStart := now.Add(-r.windowSize)
	
	// Get user's request history
	userRequests, exists := r.requests[userID]
	if !exists {
		userRequests = []time.Time{}
	}
	
	// Remove old requests outside the window
	var validRequests []time.Time
	for _, reqTime := range userRequests {
		if reqTime.After(windowStart) {
			validRequests = append(validRequests, reqTime)
		}
	}
	
	// Check if user has exceeded the limit
	if len(validRequests) >= r.maxRequests {
		return fmt.Errorf("rate limit exceeded: %d requests in %v", r.maxRequests, r.windowSize)
	}
	
	// Add current request
	validRequests = append(validRequests, now)
	r.requests[userID] = validRequests
	
	return nil
}

// GetRateLimitInfo returns rate limit information for a user
func (r *RateLimiter) GetRateLimitInfo(userID string) *RateLimitInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	now := time.Now()
	windowStart := now.Add(-r.windowSize)
	
	userRequests, exists := r.requests[userID]
	if !exists {
		return &RateLimitInfo{
			UserID:      userID,
			Requests:    0,
			MaxRequests: r.maxRequests,
			WindowSize:  r.windowSize,
			ResetTime:   now.Add(r.windowSize),
		}
	}
	
	// Count valid requests
	var validRequests int
	for _, reqTime := range userRequests {
		if reqTime.After(windowStart) {
			validRequests++
		}
	}
	
	return &RateLimitInfo{
		UserID:      userID,
		Requests:    validRequests,
		MaxRequests: r.maxRequests,
		WindowSize:  r.windowSize,
		ResetTime:   now.Add(r.windowSize),
	}
}

// Reset clears the internal request tracking state.
func (r *RateLimiter) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.requests = make(map[string][]time.Time)
}

// RateLimitInfo contains rate limit information
type RateLimitInfo struct {
	UserID      string        `json:"user_id"`
	Requests    int           `json:"requests"`
	MaxRequests int           `json:"max_requests"`
	WindowSize  time.Duration `json:"window_size"`
	ResetTime   time.Time     `json:"reset_time"`
}

// Middleware creates an HTTP middleware for authentication
func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip authentication for health check endpoints
		if r.URL.Path == "/health" || r.URL.Path == "/ready" {
			next.ServeHTTP(w, r)
			return
		}
		
		// Authenticate request
		authResult, err := a.AuthenticateRequest(r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Authentication failed: %v", err), http.StatusUnauthorized)
			return
		}
		
		// Add authentication result to request context
		ctx := context.WithValue(r.Context(), contextKey("auth"), authResult)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetAuthFromContext extracts authentication result from context
func GetAuthFromContext(ctx context.Context) (*AuthResult, bool) {
	auth, ok := ctx.Value("auth").(*AuthResult)
	return auth, ok
}
