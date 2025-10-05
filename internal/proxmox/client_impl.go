package proxmox

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"mini-mcp/internal/proxmox/config"
	"mini-mcp/internal/proxmox/interfaces"
	"mini-mcp/internal/proxmox/types"
)

// client implements the ProxmoxClient interface
type client struct {
	*BaseClient
}

// NewClient creates a new Proxmox client
func NewClient(authConfig *types.AuthConfig) interfaces.ProxmoxClient {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !config.GetVerifySSL(authConfig),
		},
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(config.GetTimeout(authConfig)) * time.Second,
	}

	baseURL := fmt.Sprintf("https://%s:8006/api2/json", config.GetHost(authConfig))
	baseClient := NewBaseClient(baseURL, httpClient, authConfig)

	return &client{
		BaseClient: baseClient,
	}
}

// Authenticate authenticates with the Proxmox API
func (c *client) Authenticate(ctx context.Context) error {
	if c.useTokenAuth {
		// Token-based authentication - no need to call /access/ticket
		c.authenticated = true
		return nil
	}

	// Username/password authentication
	data := url.Values{
		"username": {c.username},
		"password": {c.password},
	}

	body, err := c.Post(ctx, "/access/ticket", data)
	if err != nil {
		return fmt.Errorf("authentication request failed: %v", err)
	}

	var authResp struct {
		Data struct {
			Ticket    string `json:"ticket"`
			CSRFToken string `json:"CSRFPreventionToken"`
			Username  string `json:"username"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &authResp); err != nil {
		return fmt.Errorf("failed to unmarshal auth response: %v", err)
	}

	c.SetAuthCredentials(authResp.Data.Ticket, authResp.Data.CSRFToken)
	c.username = authResp.Data.Username

	return nil
}

// IsAuthenticated returns whether the client is authenticated
func (c *client) IsAuthenticated() bool {
	return c.BaseClient.IsAuthenticated()
}

// GetNodes returns a list of all nodes
func (c *client) GetNodes(ctx context.Context) ([]types.Node, error) {
	var nodes []types.Node
	err := c.GetListAndUnmarshal(ctx, "/nodes", nil, &nodes)
	return nodes, err
}

// GetNodeStatus returns the status of a specific node
func (c *client) GetNodeStatus(ctx context.Context, nodeName string) (*types.NodeStatus, error) {
	var status types.NodeStatus
	err := c.GetAndUnmarshal(ctx, c.NodeStatusEndpoint(nodeName), nil, &status)
	return &status, err
}

// GetVMs returns a list of all VMs on a specific node
func (c *client) GetVMs(ctx context.Context, nodeName string) ([]types.VM, error) {
	var vms []types.VM
	err := c.GetListAndUnmarshal(ctx, c.NodeQemuEndpoint(nodeName), nil, &vms)
	return vms, err
}

// GetVM returns a specific VM
func (c *client) GetVM(ctx context.Context, nodeName string, vmid int) (*types.VM, error) {
	vms, err := c.GetVMs(ctx, nodeName)
	if err != nil {
		return nil, err
	}

	for _, vm := range vms {
		if vm.VMID == vmid {
			return &vm, nil
		}
	}

	return nil, fmt.Errorf("VM %d not found", vmid)
}

// GetVMStatus returns the status of a specific VM
func (c *client) GetVMStatus(ctx context.Context, nodeName string, vmid int) (*types.VMStatus, error) {
	var status types.VMStatus
	err := c.GetAndUnmarshal(ctx, c.VMStatusEndpoint(nodeName, vmid), nil, &status)
	return &status, err
}

// GetVMConfig returns the configuration of a specific VM
func (c *client) GetVMConfig(ctx context.Context, nodeName string, vmid int) (*types.VMConfig, error) {
	var config types.VMConfig
	err := c.GetAndUnmarshal(ctx, c.VMConfigEndpoint(nodeName, vmid), nil, &config)
	return &config, err
}

// CreateVM creates a new VM
func (c *client) CreateVM(ctx context.Context, nodeName string, config types.VMCreateRequest) error {
	form := url.Values{}
	form.Set("vmid", fmt.Sprintf("%d", config.VMID))
	form.Set("name", config.Name)
	form.Set("memory", fmt.Sprintf("%d", config.Memory))
	form.Set("cores", fmt.Sprintf("%d", config.Cores))
	form.Set("ostype", config.OSType)
	form.Set("scsihw", config.SCSIHw)
	form.Set("net0", config.Net0)
	form.Set("agent", config.Agent)
	form.Set("vga", config.VGA)
	form.Set("boot", config.Boot)
	form.Set("bootdisk", config.BootDisk)

	if config.OnBoot {
		form.Set("onboot", "1")
	}
	if config.Start {
		form.Set("start", "1")
	}
	if config.IPConfig0 != "" {
		form.Set("ipconfig0", config.IPConfig0)
	}
	if config.Nameserver != "" {
		form.Set("nameserver", config.Nameserver)
	}
	if config.SearchDomain != "" {
		form.Set("searchdomain", config.SearchDomain)
	}

	// Cloud-init configuration
	if config.CIUser != "" {
		form.Set("ciuser", config.CIUser)
	}
	if config.CIPassword != "" {
		form.Set("cipassword", config.CIPassword)
	}
	if config.SSHKeys != "" {
		form.Set("sshkeys", config.SSHKeys)
	}
	// Note: SSH keys are set after VM creation via UpdateVMConfig to avoid encoding issues
	if config.Nameserver != "" {
		form.Set("nameserver", config.Nameserver)
	}
	if config.SearchDomain != "" {
		form.Set("searchdomain", config.SearchDomain)
	}
	if config.CIDNS != "" {
		form.Set("cidns", config.CIDNS)
	}
	if config.CISearch != "" {
		form.Set("cisearch", config.CISearch)
	}
	// Remove ciuserdata as it's not supported; use cicustom for custom files
	if config.CINetwork != "" {
		form.Set("cinetwork", config.CINetwork)
	}

	// Add cloud-init drive if cloud-init is configured
	if config.CIUser != "" || config.CIPassword != "" {
		form.Set("ide2", config.Storage+":cloudinit")
	}

	// Storage and disk configuration
	if config.Storage != "" {
		form.Set("storage", config.Storage)
	}
	if config.IDE0 != "" {
		form.Set("ide0", config.IDE0)
	}
	if config.SATA0 != "" {
		form.Set("sata0", config.SATA0)
	}
	if config.SCSI0 != "" {
		form.Set("scsi0", config.SCSI0)
	}

	_, err := c.Post(ctx, c.NodeQemuEndpoint(nodeName), form)
	return err
}

// UpdateVMConfig updates VM configuration
func (c *client) UpdateVMConfig(ctx context.Context, nodeName string, vmid int, config map[string]string) error {
	// If sshkeys is present, send it as a raw, unencoded value. Proxmox's
	// parameter parser is sensitive to certain encodings for multiline
	// ssh key blocks; sending the raw newline-separated keys in the body
	// (sshkeys=<raw keys>) has proven to work against some Proxmox versions.
	if val, ok := config["sshkeys"]; ok {
		bodyStr := "sshkeys=" + val
		_, err := c.Put(ctx, c.VMConfigEndpoint(nodeName, vmid), bodyStr)
		return err
	}

	// Default: use url.Values (encoded) for other parameters
	form := url.Values{}
	for key, value := range config {
		form.Set(key, value)
	}

	_, err := c.Put(ctx, c.VMConfigEndpoint(nodeName, vmid), form)
	return err
}

// StartVM starts a VM
func (c *client) StartVM(ctx context.Context, nodeName string, vmid int) error {
	_, err := c.Post(ctx, c.VMStartEndpoint(nodeName, vmid), nil)
	return err
}

// StopVM stops a VM
func (c *client) StopVM(ctx context.Context, nodeName string, vmid int) error {
	_, err := c.Post(ctx, c.VMStopEndpoint(nodeName, vmid), nil)
	return err
}

// ShutdownVM shuts down a VM gracefully
func (c *client) ShutdownVM(ctx context.Context, nodeName string, vmid int) error {
	_, err := c.Post(ctx, c.VMShutdownEndpoint(nodeName, vmid), nil)
	return err
}

// RebootVM reboots a VM
func (c *client) RebootVM(ctx context.Context, nodeName string, vmid int) error {
	_, err := c.Post(ctx, c.VMRebootEndpoint(nodeName, vmid), nil)
	return err
}

// DeleteVM deletes a VM
func (c *client) DeleteVM(ctx context.Context, nodeName string, vmid int) error {
	_, err := c.Delete(ctx, c.VMDeleteEndpoint(nodeName, vmid))
	return err
}

// GetStorages returns a list of all storages on a specific node
func (c *client) GetStorages(ctx context.Context, nodeName string) ([]types.Storage, error) {
	var storages []types.Storage
	err := c.GetListAndUnmarshal(ctx, c.NodeStorageEndpoint(nodeName), nil, &storages)
	return storages, err
}

// GetStorageContent returns the content of a specific storage
func (c *client) GetStorageContent(ctx context.Context, nodeName, storageName string) ([]types.StorageContent, error) {
	var content []types.StorageContent
	err := c.GetListAndUnmarshal(ctx, c.StorageContentEndpoint(nodeName, storageName), nil, &content)
	return content, err
}

// UploadFile uploads a file to storage
func (c *client) UploadFile(ctx context.Context, nodeName, storageName, filename string, data io.Reader) error {
	// Read the file data into memory so we can write it into multiple parts
	// if necessary (some Proxmox servers accept 'file', others expect
	// 'filename'). Keep content in memory since user-data files are small.
	fileBuf := &bytes.Buffer{}
	if _, err := io.Copy(fileBuf, data); err != nil {
		return fmt.Errorf("failed to read file data: %v", err)
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Decide content type based on filename. Cloud-init user-data files (yaml)
	// should be uploaded with content 'snippets'. Default to 'iso' for other
	// cases to preserve existing behaviour.
	contentType := "iso"
	if strings.HasSuffix(strings.ToLower(filename), ".yaml") || strings.HasSuffix(strings.ToLower(filename), ".yml") {
		contentType = "snippets"
	}

	// Write fields in the same order as the curl example used during manual
	// testing: content, filename, then the file part named 'file'. Some
	// Proxmox API implementations validate field ordering and names strictly.
	if err := writer.WriteField("content", contentType); err != nil {
		_ = writer.Close()
		return fmt.Errorf("failed to write content field: %v", err)
	}
	if err := writer.WriteField("filename", filename); err != nil {
		_ = writer.Close()
		return fmt.Errorf("failed to write filename field: %v", err)
	}

	filePart, err := writer.CreateFormFile("file", filename)
	if err != nil {
		_ = writer.Close()
		return fmt.Errorf("failed to create file part: %v", err)
	}
	if _, err := io.Copy(filePart, bytes.NewReader(fileBuf.Bytes())); err != nil {
		_ = writer.Close()
		return fmt.Errorf("failed to write file part: %v", err)
	}

	// (fields were already written above)

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close multipart writer: %v", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.baseURL+c.StorageUploadEndpoint(nodeName, storageName), &buf)
	if err != nil {
		return fmt.Errorf("failed to create upload request: %v", err)
	}

	// Debug: log request metadata (content-type and size)
	contentType = writer.FormDataContentType()
	req.Header.Set("Content-Type", contentType)
	c.setAuthHeaders(req)
	// Debug logging removed for production

	// Perform request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		// Read response body for debugging
		body, _ := io.ReadAll(resp.Body)

		// Log response headers and body for diagnosis
		// Error logging handled by structured logging

		return fmt.Errorf("failed to upload file with status %d: %s", resp.StatusCode, string(body))
	}

	// Upload successful - logged via structured logging
	return nil
}

// DownloadFile downloads a file from storage
func (c *client) DownloadFile(ctx context.Context, nodeName, storageName, filename string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		c.baseURL+c.StorageContentFileEndpoint(nodeName, storageName, filename), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create download request: %v", err)
	}

	c.setAuthHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("failed to download file with status %d", resp.StatusCode)
	}

	return resp.Body, nil
}

// GetNetworkInterfaces returns a list of network interfaces on a node
func (c *client) GetNetworkInterfaces(ctx context.Context, nodeName string) ([]types.NetworkInterface, error) {
	var result []types.NetworkInterface
	err := c.GetAndUnmarshal(ctx, c.NodeNetworkEndpoint(nodeName), nil, &result)
	return result, err
}

// GetNetworkInterface returns a specific network interface configuration
func (c *client) GetNetworkInterface(ctx context.Context, nodeName, iface string) (*types.NetworkInterface, error) {
	var result types.NetworkInterface
	err := c.GetAndUnmarshal(ctx, c.NodeNetworkInterfaceEndpoint(nodeName, iface), nil, &result)
	return &result, err
}

// CreateNetworkInterface creates a new network interface
func (c *client) CreateNetworkInterface(ctx context.Context, nodeName, iface string, config types.NetworkInterfaceConfig) error {
	form := url.Values{}
	form.Set("type", config.Type)

	if config.Bridge != "" {
		form.Set("bridge", config.Bridge)
	}
	if config.VLAN > 0 {
		form.Set("vlan", fmt.Sprintf("%d", config.VLAN))
	}
	if config.Address != "" {
		form.Set("address", config.Address)
	}
	if config.Netmask != "" {
		form.Set("netmask", config.Netmask)
	}
	if config.Gateway != "" {
		form.Set("gateway", config.Gateway)
	}
	if config.BondMode != "" {
		form.Set("bond_mode", config.BondMode)
	}
	if config.BondSlaves != "" {
		form.Set("bond_slaves", config.BondSlaves)
	}
	if config.Comments != "" {
		form.Set("comments", config.Comments)
	}

	_, err := c.Post(ctx, c.NodeNetworkEndpoint(nodeName), form)
	return err
}

// UpdateNetworkInterface updates an existing network interface
func (c *client) UpdateNetworkInterface(ctx context.Context, nodeName, iface string, config types.NetworkInterfaceConfig) error {
	form := url.Values{}
	form.Set("type", config.Type)

	if config.Bridge != "" {
		form.Set("bridge", config.Bridge)
	}
	if config.VLAN > 0 {
		form.Set("vlan", fmt.Sprintf("%d", config.VLAN))
	}
	if config.Address != "" {
		form.Set("address", config.Address)
	}
	if config.Netmask != "" {
		form.Set("netmask", config.Netmask)
	}
	if config.Gateway != "" {
		form.Set("gateway", config.Gateway)
	}
	if config.BondMode != "" {
		form.Set("bond_mode", config.BondMode)
	}
	if config.BondSlaves != "" {
		form.Set("bond_slaves", config.BondSlaves)
	}
	if config.Comments != "" {
		form.Set("comments", config.Comments)
	}

	_, err := c.Put(ctx, c.NodeNetworkInterfaceEndpoint(nodeName, iface), form)
	return err
}

// DeleteNetworkInterface deletes a network interface
func (c *client) DeleteNetworkInterface(ctx context.Context, nodeName, iface string) error {
	_, err := c.Delete(ctx, c.NodeNetworkInterfaceEndpoint(nodeName, iface))
	return err
}

// GetImages returns a list of all images
func (c *client) GetImages(ctx context.Context, nodeName string) ([]types.Image, error) {
	// Get all storages first
	storages, err := c.GetStorages(ctx, nodeName)
	if err != nil {
		return nil, fmt.Errorf("failed to get storages: %v", err)
	}

	var allImages []types.Image
	for _, storage := range storages {
		// Get content for each storage
		content, err := c.GetStorageContent(ctx, nodeName, storage.Storage)
		if err != nil {
			continue // Skip storages that can't be accessed
		}

		// Filter for images (content types: iso, vztmpl, etc.)
		for _, item := range content {
			if item.Content == "iso" || item.Content == "vztmpl" || item.Content == "images" {
				image := types.Image{
					VolID:   item.VolID,
					Format:  item.Format,
					Size:    item.Size,
					CTime:   time.Unix(item.CTime, 0),
					Path:    item.Path,
					Content: item.Content,
					Name:    item.VolID, // Use VolID as name for now
					Storage: storage.Storage,
				}
				allImages = append(allImages, image)
			}
		}
	}

	return allImages, nil
}

// ImportImage imports an image
func (c *client) ImportImage(ctx context.Context, nodeName, storageName, imagePath string) error {
	// Create form data for image import
	form := url.Values{}
	form.Set("filename", imagePath)
	form.Set("content", "iso")

	_, err := c.Post(ctx, c.StorageUploadEndpoint(nodeName, storageName), form)
	return err
}

// DeleteImage deletes an image
func (c *client) DeleteImage(ctx context.Context, nodeName, storageName, imageName string) error {
	_, err := c.Delete(ctx, c.StorageContentFileEndpoint(nodeName, storageName, imageName))
	return err
}

// GetFirewallRules returns a list of firewall rules
func (c *client) GetFirewallRules(ctx context.Context, nodeName string) ([]types.FirewallRule, error) {
	var result []types.FirewallRule
	err := c.GetAndUnmarshal(ctx, c.NodeFirewallRulesEndpoint(nodeName), nil, &result)
	return result, err
}

// GetFirewallRule returns a specific firewall rule
func (c *client) GetFirewallRule(ctx context.Context, nodeName string, pos int) (*types.FirewallRule, error) {
	var result types.FirewallRule
	err := c.GetAndUnmarshal(ctx, c.NodeFirewallRuleEndpoint(nodeName, pos), nil, &result)
	return &result, err
}

// CreateFirewallRule creates a new firewall rule
func (c *client) CreateFirewallRule(ctx context.Context, nodeName string, rule types.FirewallRule) error {
	form := url.Values{}
	form.Set("action", rule.Action)
	form.Set("type", rule.Type)

	if rule.Source != "" {
		form.Set("source", rule.Source)
	}
	if rule.Dest != "" {
		form.Set("dest", rule.Dest)
	}
	if rule.Proto != "" {
		form.Set("proto", rule.Proto)
	}
	if rule.Sport != "" {
		form.Set("sport", rule.Sport)
	}
	if rule.Dport != "" {
		form.Set("dport", rule.Dport)
	}
	if rule.Macro != "" {
		form.Set("macro", rule.Macro)
	}
	if rule.Target != "" {
		form.Set("target", rule.Target)
	}
	if rule.Log != "" {
	}
	if rule.Comment != "" {
		form.Set("comment", rule.Comment)
	}
	form.Set("enable", fmt.Sprintf("%d", rule.Enable))

	_, err := c.Post(ctx, c.NodeFirewallRulesEndpoint(nodeName), form)
	return err
}

// UpdateFirewallRule updates an existing firewall rule
func (c *client) UpdateFirewallRule(ctx context.Context, nodeName string, pos int, rule types.FirewallRule) error {
	form := url.Values{}
	form.Set("action", rule.Action)
	form.Set("type", rule.Type)

	if rule.Source != "" {
		form.Set("source", rule.Source)
	}
	if rule.Dest != "" {
		form.Set("dest", rule.Dest)
	}
	if rule.Proto != "" {
		form.Set("proto", rule.Proto)
	}
	if rule.Sport != "" {
		form.Set("sport", rule.Sport)
	}
	if rule.Dport != "" {
		form.Set("dport", rule.Dport)
	}
	if rule.Macro != "" {
		form.Set("macro", rule.Macro)
	}
	if rule.Target != "" {
		form.Set("target", rule.Target)
	}
	if rule.Log != "" {
	}
	if rule.Comment != "" {
		form.Set("comment", rule.Comment)
	}
	form.Set("enable", fmt.Sprintf("%d", rule.Enable))

	_, err := c.Put(ctx, c.NodeFirewallRuleEndpoint(nodeName, pos), form)
	return err
}

// DeleteFirewallRule deletes a firewall rule
func (c *client) DeleteFirewallRule(ctx context.Context, nodeName string, pos int) error {
	_, err := c.Delete(ctx, c.NodeFirewallRuleEndpoint(nodeName, pos))
	return err
}

// GetFirewallOptions returns firewall options
func (c *client) GetFirewallOptions(ctx context.Context, nodeName string) (*types.FirewallOptions, error) {
	var result types.FirewallOptions
	err := c.GetAndUnmarshal(ctx, c.NodeFirewallOptionsEndpoint(nodeName), nil, &result)
	return &result, err
}

// UpdateFirewallOptions updates firewall options
func (c *client) UpdateFirewallOptions(ctx context.Context, nodeName string, options types.FirewallOptions) error {
	form := url.Values{}
	form.Set("enable", fmt.Sprintf("%d", options.Enable))

	// Only set basic options that are supported by the API
	// Other options may not be supported in all Proxmox versions

	_, err := c.Put(ctx, c.NodeFirewallOptionsEndpoint(nodeName), form)
	return err
}

// GetFirewallAliases returns firewall aliases
func (c *client) GetFirewallAliases(ctx context.Context, nodeName string) ([]types.FirewallAlias, error) {
	var result []types.FirewallAlias
	err := c.GetAndUnmarshal(ctx, c.NodeFirewallAliasesEndpoint(nodeName), nil, &result)
	return result, err
}

// CreateFirewallAlias creates a firewall alias
func (c *client) CreateFirewallAlias(ctx context.Context, nodeName string, alias types.FirewallAlias) error {
	form := url.Values{}
	form.Set("name", alias.Name)
	form.Set("cidr", alias.CIDR)
	if alias.Comment != "" {
		form.Set("comment", alias.Comment)
	}

	_, err := c.Post(ctx, c.NodeFirewallAliasesEndpoint(nodeName), form)
	return err
}

// UpdateFirewallAlias updates a firewall alias
func (c *client) UpdateFirewallAlias(ctx context.Context, nodeName, name string, alias types.FirewallAlias) error {
	form := url.Values{}
	form.Set("cidr", alias.CIDR)
	if alias.Comment != "" {
		form.Set("comment", alias.Comment)
	}

	_, err := c.Put(ctx, c.NodeFirewallAliasEndpoint(nodeName, name), form)
	return err
}

// DeleteFirewallAlias deletes a firewall alias
func (c *client) DeleteFirewallAlias(ctx context.Context, nodeName, name string) error {
	_, err := c.Delete(ctx, c.NodeFirewallAliasEndpoint(nodeName, name))
	return err
}

// GetFirewallGroups returns firewall groups
func (c *client) GetFirewallGroups(ctx context.Context, nodeName string) ([]types.FirewallGroup, error) {
	var result []types.FirewallGroup
	err := c.GetAndUnmarshal(ctx, c.NodeFirewallGroupsEndpoint(nodeName), nil, &result)
	return result, err
}

// CreateFirewallGroup creates a firewall group
func (c *client) CreateFirewallGroup(ctx context.Context, nodeName string, group types.FirewallGroup) error {
	form := url.Values{}
	form.Set("group", group.Group)
	if group.Comment != "" {
		form.Set("comment", group.Comment)
	}

	_, err := c.Post(ctx, c.NodeFirewallGroupsEndpoint(nodeName), form)
	return err
}

// DeleteFirewallGroup deletes a firewall group
func (c *client) DeleteFirewallGroup(ctx context.Context, nodeName, group string) error {
	_, err := c.Delete(ctx, c.NodeFirewallGroupEndpoint(nodeName, group))
	return err
}

// SetTransport sets the HTTP transport for the client
func (c *client) SetTransport(transport http.RoundTripper) {
	c.BaseClient.SetTransport(transport)
}
