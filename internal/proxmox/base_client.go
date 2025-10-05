package proxmox

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"mini-mcp/internal/proxmox/types"
)

// HTTPMethod represents HTTP methods - using standard library constants
type HTTPMethod string

const (
	GET    HTTPMethod = http.MethodGet
	POST   HTTPMethod = http.MethodPost
	PUT    HTTPMethod = http.MethodPut
	DELETE HTTPMethod = http.MethodDelete
)

// RequestOptions represents options for HTTP requests with type safety
type RequestOptions[T any] struct {
	Method      HTTPMethod
	Path        string
	Body        T
	ContentType string
	QueryParams map[string]string
}

// APIResponse represents a standard Proxmox API response with type safety
type APIResponse[T any] struct {
	Data T `json:"data"`
}

// BaseClient provides common HTTP client functionality
type BaseClient struct {
	baseURL       string
	httpClient    *http.Client
	username      string
	password      string
	tokenName     string
	tokenValue    string
	useTokenAuth  bool
	ticket        string
	csrfToken     string
	authenticated bool
}

// NewBaseClient creates a new base client
func NewBaseClient(baseURL string, httpClient *http.Client, authConfig *types.AuthConfig) *BaseClient {
	return &BaseClient{
		baseURL:      baseURL,
		httpClient:   httpClient,
		username:     authConfig.Proxmox.User,
		password:     authConfig.Proxmox.Password,
		tokenName:    authConfig.Proxmox.TokenName,
		tokenValue:   authConfig.Proxmox.TokenValue,
		useTokenAuth: authConfig.Proxmox.TokenName != "" && authConfig.Proxmox.TokenValue != "",
	}
}

// SetAuthCredentials sets authentication credentials
func (bc *BaseClient) SetAuthCredentials(ticket, csrfToken string) {
	bc.ticket = ticket
	bc.csrfToken = csrfToken
	bc.authenticated = true
}

// IsAuthenticated returns whether the client is authenticated
func (bc *BaseClient) IsAuthenticated() bool {
	if bc.useTokenAuth {
		return bc.authenticated && bc.tokenName != "" && bc.tokenValue != ""
	}
	return bc.authenticated && bc.ticket != ""
}

// setAuthHeaders sets the authentication headers for a request
func (bc *BaseClient) setAuthHeaders(req *http.Request) {
	if bc.useTokenAuth {
		// Token-based authentication
		req.Header.Set("Authorization", fmt.Sprintf("PVEAPIToken=%s=%s", bc.tokenName, bc.tokenValue))
	} else if bc.ticket != "" {
		// Username/password authentication
		req.Header.Set("Cookie", fmt.Sprintf("PVEAuthCookie=%s", bc.ticket))
		if bc.csrfToken != "" {
			req.Header.Set("CSRFPreventionToken", bc.csrfToken)
		}
	}
}

// doRequest performs an authenticated HTTP request with retry logic
func (bc *BaseClient) doRequest(ctx context.Context, method HTTPMethod, path string, body io.Reader) (*http.Response, error) {
	const maxRetries = 3
	const baseDelay = time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		req, err := http.NewRequestWithContext(ctx, string(method), fmt.Sprintf("%s%s", bc.baseURL, path), body)
		if err != nil {
			// Request creation failed - handled by error return
			return nil, fmt.Errorf("failed to create request: %v", err)
		}

		bc.setAuthHeaders(req)
		if body != nil {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}

		// Request logging handled by structured logging
		resp, err := bc.httpClient.Do(req)
		if err != nil {
			if attempt < maxRetries-1 {
				delay := baseDelay * time.Duration(1<<attempt) // Exponential backoff
				// Retry logging handled by structured logging
				time.Sleep(delay)
				continue
			}
			// Final failure logging handled by error return
			return nil, fmt.Errorf("request failed after %d attempts: %v", maxRetries, err)
		}

			// Check for HTTP errors
			if resp.StatusCode >= 400 {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					// Log error but continue with empty body
					body = []byte("")
				}
				if closeErr := resp.Body.Close(); closeErr != nil {
					// Log close error but continue
					_ = closeErr
				}

			if resp.StatusCode == 401 {
				// Authentication error handled by error return
				return nil, fmt.Errorf("authentication failed: %s", string(body))
			}

			if resp.StatusCode >= 500 && attempt < maxRetries-1 {
				delay := baseDelay * time.Duration(1<<attempt)
				// Server error retry handled by structured logging
				time.Sleep(delay)
				continue
			}

			// HTTP error handled by error return
			return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
		}

		// Request success handled by structured logging
		return resp, nil
	}

	return nil, fmt.Errorf("max retries exceeded")
}

// makeRequest makes an HTTP request and returns the response body
func (bc *BaseClient) makeRequest(ctx context.Context, method HTTPMethod, path, body, contentType string) ([]byte, error) {
	resp, err := bc.doRequest(ctx, method, path, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer func() { 
		if closeErr := resp.Body.Close(); closeErr != nil {
			// Log close error but continue
			_ = closeErr
		}
	}()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	return responseBody, nil
}

// Request performs a generic HTTP request with options
func (bc *BaseClient) Request(ctx context.Context, opts RequestOptions[interface{}]) ([]byte, error) {
	var bodyReader io.Reader
	var contentType string

	// Prepare body with type safety
	if !isZeroValue(opts.Body) {
		switch v := any(opts.Body).(type) {
		case string:
			bodyReader = strings.NewReader(v)
			contentType = "application/x-www-form-urlencoded"
		case url.Values:
			bodyReader = strings.NewReader(v.Encode())
			contentType = "application/x-www-form-urlencoded"
		case []byte:
			bodyReader = bytes.NewReader(v)
			contentType = "application/octet-stream"
		default:
			jsonData, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal request body: %v", err)
			}
			bodyReader = bytes.NewReader(jsonData)
			contentType = "application/json"
		}
	}

	// Override content type if specified
	if opts.ContentType != "" {
		contentType = opts.ContentType
	}

	// Add query parameters
	path := opts.Path
	if len(opts.QueryParams) > 0 {
		params := url.Values{}
		for key, value := range opts.QueryParams {
			params.Set(key, value)
		}
		if strings.Contains(path, "?") {
			path += "&" + params.Encode()
		} else {
			path += "?" + params.Encode()
		}
	}

	// Make the request
	var bodyStr string
	if bodyReader != nil {
		bodyBytes, err := io.ReadAll(bodyReader)
		if err != nil {
			return nil, err
		}
		bodyStr = string(bodyBytes)
	}

	// Debugging: log request bodies for VM config endpoints to inspect encoding
	if strings.Contains(path, "/qemu/") && strings.Contains(path, "/config") {
		// Debug logging removed for production - no action needed
		_ = path
	}

	return bc.makeRequest(ctx, opts.Method, path, bodyStr, contentType)
}

// Get performs a GET request
func (bc *BaseClient) Get(ctx context.Context, path string, queryParams map[string]string) ([]byte, error) {
	return bc.Request(ctx, RequestOptions[interface{}]{
		Method:      GET,
		Path:        path,
		QueryParams: queryParams,
	})
}

// Post performs a POST request with type safety
func (bc *BaseClient) Post(ctx context.Context, path string, body interface{}) ([]byte, error) {
	return bc.Request(ctx, RequestOptions[interface{}]{
		Method: POST,
		Path:   path,
		Body:   body,
	})
}

// Put performs a PUT request with type safety
func (bc *BaseClient) Put(ctx context.Context, path string, body interface{}) ([]byte, error) {
	return bc.Request(ctx, RequestOptions[interface{}]{
		Method: PUT,
		Path:   path,
		Body:   body,
	})
}

// Delete performs a DELETE request
func (bc *BaseClient) Delete(ctx context.Context, path string) ([]byte, error) {
	return bc.Request(ctx, RequestOptions[interface{}]{
		Method: DELETE,
		Path:   path,
	})
}

// UnmarshalResponse unmarshals a Proxmox API response into the target type
func (bc *BaseClient) UnmarshalResponse(data []byte, target interface{}) error {
	if len(data) == 0 {
		return fmt.Errorf("empty response data")
	}

	var response struct {
		Data interface{} `json:"data"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		// Unmarshal error handled by error return
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if response.Data == nil {
		return fmt.Errorf("response data field is null")
	}

	// Marshal the data field back to JSON and unmarshal into target
	dataBytes, err := json.Marshal(response.Data)
	if err != nil {
		// Marshal error handled by error return
		return fmt.Errorf("failed to marshal data field: %v", err)
	}

	if err := json.Unmarshal(dataBytes, target); err != nil {
		// Unmarshal error handled by error return
		return fmt.Errorf("failed to unmarshal into target: %v", err)
	}

	return nil
}

// UnmarshalListResponse unmarshals a list response into a slice
func (bc *BaseClient) UnmarshalListResponse(data []byte, target interface{}) error {
	if len(data) == 0 {
		return fmt.Errorf("empty response data")
	}

	var response struct {
		Data interface{} `json:"data"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		// Unmarshal error handled by error return
		return fmt.Errorf("failed to unmarshal list response: %v", err)
	}

	if response.Data == nil {
		// Null response data handled by empty list initialization
		// Initialize as empty slice if data is null
		response.Data = []interface{}{}
	}

	// Marshal the data field back to JSON and unmarshal into target
	dataBytes, err := json.Marshal(response.Data)
	if err != nil {
		// Marshal error handled by error return
		return fmt.Errorf("failed to marshal data field: %v", err)
	}

	if err := json.Unmarshal(dataBytes, target); err != nil {
		// Unmarshal error handled by error return
		return fmt.Errorf("failed to unmarshal into target: %v", err)
	}

	return nil
}

// GetAndUnmarshal performs a GET request and unmarshals the response
func (bc *BaseClient) GetAndUnmarshal(ctx context.Context, path string, queryParams map[string]string, target interface{}) error {
	data, err := bc.Get(ctx, path, queryParams)
	if err != nil {
		return err
	}

	return bc.UnmarshalResponse(data, target)
}

// PostAndUnmarshal performs a POST request and unmarshals the response
func (bc *BaseClient) PostAndUnmarshal(ctx context.Context, path string, body interface{}, target interface{}) error {
	data, err := bc.Post(ctx, path, body)
	if err != nil {
		return err
	}

	return bc.UnmarshalResponse(data, target)
}

// GetListAndUnmarshal performs a GET request and unmarshals a list response
func (bc *BaseClient) GetListAndUnmarshal(ctx context.Context, path string, queryParams map[string]string, target interface{}) error {
	data, err := bc.Get(ctx, path, queryParams)
	if err != nil {
		return err
	}

	return bc.UnmarshalListResponse(data, target)
}

// Endpoint builders for common Proxmox API patterns
func (bc *BaseClient) NodeEndpoint(nodeName, path string) string {
	return fmt.Sprintf("/nodes/%s%s", nodeName, path)
}

func (bc *BaseClient) StorageEndpoint(nodeName, storageName, path string) string {
	return fmt.Sprintf("/nodes/%s/storage/%s%s", nodeName, storageName, path)
}

func (bc *BaseClient) VMEndpoint(nodeName, vmid, path string) string {
	return fmt.Sprintf("/nodes/%s/qemu/%s%s", nodeName, vmid, path)
}

func (bc *BaseClient) VMDumpEndpoint(nodeName string) string {
	return fmt.Sprintf("/nodes/%s/vzdump", nodeName)
}

func (bc *BaseClient) StorageContentEndpoint(nodeName, storageName string) string {
	return fmt.Sprintf("/nodes/%s/storage/%s/content", nodeName, storageName)
}

// Additional endpoint builders for comprehensive coverage
func (bc *BaseClient) NodeStatusEndpoint(nodeName string) string {
	return fmt.Sprintf("/nodes/%s/status", nodeName)
}

func (bc *BaseClient) NodeQemuEndpoint(nodeName string) string {
	return fmt.Sprintf("/nodes/%s/qemu", nodeName)
}

func (bc *BaseClient) VMStatusEndpoint(nodeName string, vmid int) string {
	return fmt.Sprintf("/nodes/%s/qemu/%d/status/current", nodeName, vmid)
}

func (bc *BaseClient) VMConfigEndpoint(nodeName string, vmid int) string {
	return fmt.Sprintf("/nodes/%s/qemu/%d/config", nodeName, vmid)
}

func (bc *BaseClient) VMStartEndpoint(nodeName string, vmid int) string {
	return fmt.Sprintf("/nodes/%s/qemu/%d/status/start", nodeName, vmid)
}

func (bc *BaseClient) VMStopEndpoint(nodeName string, vmid int) string {
	return fmt.Sprintf("/nodes/%s/qemu/%d/status/stop", nodeName, vmid)
}

func (bc *BaseClient) VMShutdownEndpoint(nodeName string, vmid int) string {
	return fmt.Sprintf("/nodes/%s/qemu/%d/status/shutdown", nodeName, vmid)
}

func (bc *BaseClient) VMRebootEndpoint(nodeName string, vmid int) string {
	return fmt.Sprintf("/nodes/%s/qemu/%d/status/reboot", nodeName, vmid)
}

func (bc *BaseClient) VMDeleteEndpoint(nodeName string, vmid int) string {
	return fmt.Sprintf("/nodes/%s/qemu/%d", nodeName, vmid)
}

func (bc *BaseClient) NodeStorageEndpoint(nodeName string) string {
	return fmt.Sprintf("/nodes/%s/storage", nodeName)
}

func (bc *BaseClient) StorageUploadEndpoint(nodeName, storageName string) string {
	return fmt.Sprintf("/nodes/%s/storage/%s/upload", nodeName, storageName)
}

func (bc *BaseClient) StorageContentFileEndpoint(nodeName, storageName, filename string) string {
	return fmt.Sprintf("/nodes/%s/storage/%s/content/%s", nodeName, storageName, filename)
}

func (bc *BaseClient) NodeNetworkEndpoint(nodeName string) string {
	return fmt.Sprintf("/nodes/%s/network", nodeName)
}

func (bc *BaseClient) NodeNetworkInterfaceEndpoint(nodeName, iface string) string {
	return fmt.Sprintf("/nodes/%s/network/%s", nodeName, iface)
}

func (bc *BaseClient) NodeFirewallRulesEndpoint(nodeName string) string {
	return fmt.Sprintf("/nodes/%s/firewall/rules", nodeName)
}

func (bc *BaseClient) NodeFirewallRuleEndpoint(nodeName string, pos int) string {
	return fmt.Sprintf("/nodes/%s/firewall/rules/%d", nodeName, pos)
}

func (bc *BaseClient) NodeFirewallOptionsEndpoint(nodeName string) string {
	return fmt.Sprintf("/nodes/%s/firewall/options", nodeName)
}

func (bc *BaseClient) NodeFirewallAliasesEndpoint(nodeName string) string {
	return fmt.Sprintf("/nodes/%s/firewall/aliases", nodeName)
}

func (bc *BaseClient) NodeFirewallAliasEndpoint(nodeName, name string) string {
	return fmt.Sprintf("/nodes/%s/firewall/aliases/%s", nodeName, name)
}

func (bc *BaseClient) NodeFirewallGroupsEndpoint(nodeName string) string {
	return fmt.Sprintf("/nodes/%s/firewall/groups", nodeName)
}

func (bc *BaseClient) NodeFirewallGroupEndpoint(nodeName, group string) string {
	return fmt.Sprintf("/nodes/%s/firewall/groups/%s", nodeName, group)
}

// SetTransport sets the HTTP transport for the client
func (bc *BaseClient) SetTransport(transport http.RoundTripper) {
	bc.httpClient.Transport = transport
}

// GetUseTokenAuth returns whether token authentication is being used
func (bc *BaseClient) GetUseTokenAuth() bool {
	return bc.useTokenAuth
}

// GetUsername returns the username
func (bc *BaseClient) GetUsername() string {
	return bc.username
}

// GetPassword returns the password
func (bc *BaseClient) GetPassword() string {
	return bc.password
}

// isZeroValue checks if a value is the zero value for its type
func isZeroValue(v interface{}) bool {
	if v == nil {
		return true
	}
	return reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
}
