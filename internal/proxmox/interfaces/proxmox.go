package interfaces

import (
	"context"
	"io"
	"time"

	"mini-mcp/internal/proxmox/types"
)

// ProxmoxClient defines the interface for Proxmox API operations
type ProxmoxClient interface {
	// Authentication
	Authenticate(ctx context.Context) error
	IsAuthenticated() bool

	// Node operations
	GetNodes(ctx context.Context) ([]types.Node, error)
	GetNodeStatus(ctx context.Context, nodeName string) (*types.NodeStatus, error)

	// VM operations
	GetVMs(ctx context.Context, nodeName string) ([]types.VM, error)
	GetVM(ctx context.Context, nodeName string, vmid int) (*types.VM, error)
	GetVMStatus(ctx context.Context, nodeName string, vmid int) (*types.VMStatus, error)
	GetVMConfig(ctx context.Context, nodeName string, vmid int) (*types.VMConfig, error)
	CreateVM(ctx context.Context, nodeName string, config types.VMCreateRequest) error
	UpdateVMConfig(ctx context.Context, nodeName string, vmid int, config map[string]string) error
	StartVM(ctx context.Context, nodeName string, vmid int) error
	StopVM(ctx context.Context, nodeName string, vmid int) error
	ShutdownVM(ctx context.Context, nodeName string, vmid int) error
	RebootVM(ctx context.Context, nodeName string, vmid int) error
	DeleteVM(ctx context.Context, nodeName string, vmid int) error

	// VM Clone and Migration
	CloneVM(ctx context.Context, nodeName string, vmid int, config types.VMCloneRequest) error
	MigrateVM(ctx context.Context, nodeName string, vmid int, config types.VMMigrateRequest) error

	// VM Snapshots
	CreateVMSnapshot(ctx context.Context, nodeName string, vmid int, config types.VMSnapshotCreateRequest) error
	DeleteVMSnapshot(ctx context.Context, nodeName string, vmid int, snapName string) error
	RollbackVMSnapshot(ctx context.Context, nodeName string, vmid int, config types.VMSnapshotRollbackRequest) error
	GetVMSnapshots(ctx context.Context, nodeName string, vmid int) ([]types.VMSnapshot, error)

	// VM Disk Management
	AddVMDisk(ctx context.Context, nodeName string, vmid int, config types.VMDiskAddRequest) error
	RemoveVMDisk(ctx context.Context, nodeName string, vmid int, diskID string) error
	ResizeVMDisk(ctx context.Context, nodeName string, vmid int, config types.VMDiskResizeRequest) error
	MoveVMDisk(ctx context.Context, nodeName string, vmid int, config types.VMDiskMoveRequest) error

	// VM Console Access
	GetVMConsole(ctx context.Context, nodeName string, vmid int, config types.VMConsoleRequest) (*types.VMConsole, error)

	// VM Backup and Restore
	CreateVMBackup(ctx context.Context, nodeName string, vmid int, config types.VMBackupCreateRequest) error
	GetVMBackups(ctx context.Context, nodeName string, storageName string) ([]types.VMBackup, error)
	RestoreVMBackup(ctx context.Context, nodeName string, vmid int, config types.VMBackupRestoreRequest) error

	// VM Templates
	CreateVMTemplate(ctx context.Context, nodeName string, vmid int, config types.VMTemplateCreateRequest) error
	DeployFromTemplate(ctx context.Context, nodeName string, config types.VMTemplateDeployRequest) error
	GetVMTemplates(ctx context.Context, nodeName string) ([]types.VMTemplate, error)

	// VM Monitoring and Logs
	GetVMPerformanceData(ctx context.Context, nodeName string, vmid int, config types.VMStatistics) (*types.VMPerformanceData, error)
	GetVMLogs(ctx context.Context, nodeName string, vmid int, limit int) ([]types.VMLogEntry, error)
	GetVMEvents(ctx context.Context, nodeName string, vmid int, limit int) ([]types.VMEvent, error)

	// Storage operations
	GetStorages(ctx context.Context, nodeName string) ([]types.Storage, error)
	GetStorageContent(ctx context.Context, nodeName, storageName string) ([]types.StorageContent, error)
	UploadFile(ctx context.Context, nodeName, storageName, filename string, data io.Reader) error
	DownloadFile(ctx context.Context, nodeName, storageName, filename string) (io.ReadCloser, error)

	// Image operations
	GetImages(ctx context.Context, nodeName string) ([]types.Image, error)
	ImportImage(ctx context.Context, nodeName, storageName, imagePath string) error
	DeleteImage(ctx context.Context, nodeName, storageName, imageName string) error

	// Network operations
	GetNetworkInterfaces(ctx context.Context, nodeName string) ([]types.NetworkInterface, error)
	GetNetworkInterface(ctx context.Context, nodeName, iface string) (*types.NetworkInterface, error)
	CreateNetworkInterface(ctx context.Context, nodeName, iface string, config types.NetworkInterfaceConfig) error
	UpdateNetworkInterface(ctx context.Context, nodeName, iface string, config types.NetworkInterfaceConfig) error
	DeleteNetworkInterface(ctx context.Context, nodeName, iface string) error

	// Firewall operations
	GetFirewallRules(ctx context.Context, nodeName string) ([]types.FirewallRule, error)
	GetFirewallRule(ctx context.Context, nodeName string, pos int) (*types.FirewallRule, error)
	CreateFirewallRule(ctx context.Context, nodeName string, rule types.FirewallRule) error
	UpdateFirewallRule(ctx context.Context, nodeName string, pos int, rule types.FirewallRule) error
	DeleteFirewallRule(ctx context.Context, nodeName string, pos int) error
	GetFirewallOptions(ctx context.Context, nodeName string) (*types.FirewallOptions, error)
	UpdateFirewallOptions(ctx context.Context, nodeName string, options types.FirewallOptions) error
	GetFirewallAliases(ctx context.Context, nodeName string) ([]types.FirewallAlias, error)
	CreateFirewallAlias(ctx context.Context, nodeName string, alias types.FirewallAlias) error
	UpdateFirewallAlias(ctx context.Context, nodeName, name string, alias types.FirewallAlias) error
	DeleteFirewallAlias(ctx context.Context, nodeName, name string) error
	GetFirewallGroups(ctx context.Context, nodeName string) ([]types.FirewallGroup, error)
	CreateFirewallGroup(ctx context.Context, nodeName string, group types.FirewallGroup) error
	DeleteFirewallGroup(ctx context.Context, nodeName, group string) error
}

// VMService defines the interface for VM business logic
type VMService interface {
	ListVMs(ctx context.Context) ([]types.VM, error)
	GetVM(ctx context.Context, vmid int) (*types.VM, error)
	GetVMStatus(ctx context.Context, vmid int) (*types.VMStatus, error)
	GetVMConfig(ctx context.Context, vmid int) (*types.VMConfig, error)
	CreateVM(ctx context.Context, config types.VMCreateRequest) error
	UpdateVMConfig(ctx context.Context, vmid int, config map[string]string) error
	StartVM(ctx context.Context, vmid int) error
	StopVM(ctx context.Context, vmid int) error
	ShutdownVM(ctx context.Context, vmid int) error
	RebootVM(ctx context.Context, vmid int) error
	DeleteVM(ctx context.Context, vmid int) error
	PrintVMStatus(ctx context.Context) error

	// VM Clone and Migration
	CloneVM(ctx context.Context, vmid int, config types.VMCloneRequest) error
	MigrateVM(ctx context.Context, vmid int, config types.VMMigrateRequest) error

	// VM Snapshots
	CreateVMSnapshot(ctx context.Context, vmid int, config types.VMSnapshotCreateRequest) error
	DeleteVMSnapshot(ctx context.Context, vmid int, snapName string) error
	RollbackVMSnapshot(ctx context.Context, vmid int, config types.VMSnapshotRollbackRequest) error
	GetVMSnapshots(ctx context.Context, vmid int) ([]types.VMSnapshot, error)

	// VM Disk Management
	AddVMDisk(ctx context.Context, vmid int, config types.VMDiskAddRequest) error
	RemoveVMDisk(ctx context.Context, vmid int, diskID string) error
	ResizeVMDisk(ctx context.Context, vmid int, config types.VMDiskResizeRequest) error
	MoveVMDisk(ctx context.Context, vmid int, config types.VMDiskMoveRequest) error

	// VM Console Access
	GetVMConsole(ctx context.Context, vmid int, config types.VMConsoleRequest) (*types.VMConsole, error)

	// VM Backup and Restore
	CreateVMBackup(ctx context.Context, vmid int, config types.VMBackupCreateRequest) error
	GetVMBackups(ctx context.Context, storageName string) ([]types.VMBackup, error)
	RestoreVMBackup(ctx context.Context, vmid int, config types.VMBackupRestoreRequest) error

	// VM Templates
	CreateVMTemplate(ctx context.Context, vmid int, config types.VMTemplateCreateRequest) error
	DeployFromTemplate(ctx context.Context, config types.VMTemplateDeployRequest) error
	GetVMTemplates(ctx context.Context) ([]types.VMTemplate, error)

	// VM Monitoring and Logs
	GetVMPerformanceData(ctx context.Context, vmid int, config types.VMStatistics) (*types.VMPerformanceData, error)
	GetVMLogs(ctx context.Context, vmid int, limit int) ([]types.VMLogEntry, error)
	GetVMEvents(ctx context.Context, vmid int, limit int) ([]types.VMEvent, error)
}

// ImageService defines the interface for image management
type ImageService interface {
	ListImages(ctx context.Context) ([]types.Image, error)
	DownloadUbuntu(ctx context.Context, version string) error
	ImportImage(ctx context.Context, imagePath, storageName string) error
	DeleteImage(ctx context.Context, imageName, storageName string) error
	GetImageInfo(ctx context.Context, imageName, storageName string) (*types.Image, error)
	GetAvailableImageSources() []types.ImageSource
	GetImageSourceByTypeAndVersion(imageType, version string) (*types.ImageSource, error)
}

// NodeService defines the interface for node management
type NodeService interface {
	ListNodes(ctx context.Context) ([]types.Node, error)
	GetNodeStatus(ctx context.Context, nodeName string) (*types.NodeStatus, error)
	GetNodeInfo(ctx context.Context, nodeName string) (*types.Node, error)
}

// NetworkService defines the interface for network management
type NetworkService interface {
	ListNetworkInterfaces(ctx context.Context, nodeName string) ([]types.NetworkInterface, error)
	GetNetworkInterface(ctx context.Context, nodeName, iface string) (*types.NetworkInterface, error)
	CreateNetworkInterface(ctx context.Context, nodeName string, config types.NetworkInterfaceConfig) error
	UpdateNetworkInterface(ctx context.Context, nodeName, iface string, config types.NetworkInterfaceConfig) error
	DeleteNetworkInterface(ctx context.Context, nodeName, iface string) error
}

// FirewallService defines the interface for firewall management
type FirewallService interface {
	// Firewall rules
	ListFirewallRules(ctx context.Context, nodeName string) ([]types.FirewallRule, error)
	GetFirewallRule(ctx context.Context, nodeName string, pos int) (*types.FirewallRule, error)
	CreateFirewallRule(ctx context.Context, nodeName string, rule types.FirewallRule) error
	UpdateFirewallRule(ctx context.Context, nodeName string, pos int, rule types.FirewallRule) error
	DeleteFirewallRule(ctx context.Context, nodeName string, pos int) error

	// Firewall options
	GetFirewallOptions(ctx context.Context, nodeName string) (*types.FirewallOptions, error)
	UpdateFirewallOptions(ctx context.Context, nodeName string, options types.FirewallOptions) error

	// Firewall aliases
	ListFirewallAliases(ctx context.Context, nodeName string) ([]types.FirewallAlias, error)
	CreateFirewallAlias(ctx context.Context, nodeName string, alias types.FirewallAlias) error
	UpdateFirewallAlias(ctx context.Context, nodeName, name string, alias types.FirewallAlias) error
	DeleteFirewallAlias(ctx context.Context, nodeName, name string) error

	// Firewall groups
	ListFirewallGroups(ctx context.Context, nodeName string) ([]types.FirewallGroup, error)
	CreateFirewallGroup(ctx context.Context, nodeName string, group types.FirewallGroup) error
	DeleteFirewallGroup(ctx context.Context, nodeName, group string) error
}

// StorageService defines the interface for storage management
type StorageService interface {
	ListStorages(ctx context.Context) ([]types.Storage, error)
	GetStorageInfo(ctx context.Context, storageName string) (*types.Storage, error)
	GetStorageContent(ctx context.Context, storageName string) ([]types.StorageContent, error)
	UploadFile(ctx context.Context, storageName, filename string, data io.Reader) error
	DownloadFile(ctx context.Context, storageName, filename string) (io.ReadCloser, error)
}

// CloudInitService provides helper operations for cloud-init snippets
type CloudInitService interface {
	UploadAndAttachSnippet(ctx context.Context, vmid int, storage, filename string, data io.Reader) error
}

// DeployService defines the interface for VM deployment
type DeployService interface {
	DeployUbuntu(ctx context.Context, vmid int, config types.UbuntuDeployConfig) (string, error)
	DeployFromTemplate(ctx context.Context, vmid int, templateID int, config types.TemplateDeployConfig) error
	DeployCustom(ctx context.Context, vmid int, config types.CustomDeployConfig) error
}

// ConfigService defines the interface for configuration management
type ConfigService interface {
	LoadAuthConfig(filename string) (*types.AuthConfig, error)
	LoadMainConfig(filename string) (*types.MainConfig, error)
	SaveConfig(config interface{}, filename string) error
	ValidateConfig(config interface{}) error
}

// SSHService defines the interface for SSH operations
type SSHService interface {
	ExecuteCommand(ctx context.Context, host, command string) (string, error)
	ExecuteCommandWithTimeout(ctx context.Context, host, command string, timeout time.Duration) (string, error)
}

// Logger defines the interface for structured logging with type safety
type Logger interface {
	Debug(msg string, fields ...LogField)
	Info(msg string, fields ...LogField)
	Warn(msg string, fields ...LogField)
	Error(msg string, fields ...LogField)
	Fatal(msg string, fields ...LogField)
	WithField(key string, value any) Logger
	WithFields(fields map[string]any) Logger
}

// LogField represents a structured log field with type safety
type LogField struct {
	Key   string
	Value any
}

// NewLogField creates a new log field
func NewLogField(key string, value any) LogField {
	return LogField{Key: key, Value: value}
}
