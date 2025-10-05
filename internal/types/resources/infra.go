// Package resources contains types for MCP resource responses.
package resources

import "time"

// InfrastructureInfo represents comprehensive infrastructure information.
// This is the root structure returned by the infrastructure_info tool.
type InfrastructureInfo struct {
	// CollectionTime is the time when this information was collected
	CollectionTime time.Time `json:"collection_time"`
	// CollectionLevel indicates the detail level of the collected information
	CollectionLevel string `json:"collection_level"`
	// SystemInfo contains basic system information
	SystemInfo *SystemInfo `json:"system_info,omitempty"`
	// ProcessInfo contains process information
	ProcessInfo *ProcessInfo `json:"process_info,omitempty"`
	// NetworkInfo contains network information
	NetworkInfo *NetworkInfo `json:"network_info,omitempty"`
	// ServiceInfo contains service information
	ServiceInfo *ServiceInfo `json:"service_info,omitempty"`
	// DockerInfo contains Docker information
	DockerInfo *DockerInfo `json:"docker_info,omitempty"`
	// KubernetesInfo contains Kubernetes information
	KubernetesInfo *KubernetesInfo `json:"kubernetes_info,omitempty"`
	// CloudInfo contains cloud provider information
	CloudInfo *CloudInfo `json:"cloud_info,omitempty"`
	// ProxmoxInfo contains Proxmox information
	ProxmoxInfo *ProxmoxInfo `json:"proxmox_info,omitempty"`
	// CephInfo contains Ceph information
	CephInfo *CephInfo `json:"ceph_info,omitempty"`
	// MetricSummary contains aggregated metrics from various sources
	MetricSummary *MetricSummary `json:"metric_summary,omitempty"`
	// Alerts contains any alerts or warnings detected during collection
	Alerts []Alert `json:"alerts,omitempty"`
}

// NewInfrastructureInfo creates a new InfrastructureInfo with the given parameters
func NewInfrastructureInfo(collectionLevel string) *InfrastructureInfo {
	return &InfrastructureInfo{
		CollectionTime:  time.Now(),
		CollectionLevel: collectionLevel,
		Alerts:          make([]Alert, 0),
	}
}

// ProcessInfo represents information about processes
type ProcessInfo struct {
	// TotalProcesses is the total number of processes
	TotalProcesses int `json:"total_processes"`
	// Processes is a list of processes
	Processes []Process `json:"processes,omitempty"`
}

// Process represents a single process
type Process struct {
	// PID is the process ID
	PID int `json:"pid"`
	// Name is the process name
	Name string `json:"name"`
	// User is the user running the process
	User string `json:"user,omitempty"`
	// CommandLine is the command line of the process
	CommandLine string `json:"command_line,omitempty"`
	// StartTime is when the process started
	StartTime time.Time `json:"start_time,omitempty"`
	// CPUPercent is the CPU usage percentage
	CPUPercent float64 `json:"cpu_percent,omitempty"`
	// MemoryPercent is the memory usage percentage
	MemoryPercent float64 `json:"memory_percent,omitempty"`
	// MemoryUsage is the memory usage in bytes
	MemoryUsage uint64 `json:"memory_usage,omitempty"`
	// IOStats contains I/O statistics
	IOStats *ProcessIOStats `json:"io_stats,omitempty"`
	// Status is the process status (running, sleeping, etc.)
	Status string `json:"status,omitempty"`
	// Threads is the number of threads
	Threads int `json:"threads,omitempty"`
}

// ProcessIOStats contains I/O statistics for a process
type ProcessIOStats struct {
	// ReadBytes is the number of bytes read
	ReadBytes uint64 `json:"read_bytes"`
	// WriteBytes is the number of bytes written
	WriteBytes uint64 `json:"write_bytes"`
	// ReadOperations is the number of read operations
	ReadOperations uint64 `json:"read_operations,omitempty"`
	// WriteOperations is the number of write operations
	WriteOperations uint64 `json:"write_operations,omitempty"`
}

// NetworkInfo represents information about the network
type NetworkInfo struct {
	// Interfaces is a list of network interfaces
	Interfaces []NetworkInterface `json:"interfaces,omitempty"`
	// Connections is a list of network connections
	Connections []NetworkConnection `json:"connections,omitempty"`
	// RoutingTable contains routing information
	RoutingTable []Route `json:"routing_table,omitempty"`
	// DNSConfig contains DNS configuration
	DNSConfig *DNSConfig `json:"dns_config,omitempty"`
	// FirewallRules contains firewall rules
	FirewallRules []FirewallRule `json:"firewall_rules,omitempty"`
}

// NetworkConnection represents a network connection
type NetworkConnection struct {
	// Protocol is the connection protocol (TCP, UDP, etc.)
	Protocol string `json:"protocol"`
	// LocalAddress is the local address
	LocalAddress string `json:"local_address"`
	// LocalPort is the local port
	LocalPort int `json:"local_port"`
	// RemoteAddress is the remote address
	RemoteAddress string `json:"remote_address"`
	// RemotePort is the remote port
	RemotePort int `json:"remote_port"`
	// State is the connection state (ESTABLISHED, LISTEN, etc.)
	State string `json:"state,omitempty"`
	// PID is the process ID that owns this connection
	PID int `json:"pid,omitempty"`
	// ProcessName is the name of the process that owns this connection
	ProcessName string `json:"process_name,omitempty"`
}

// Route represents a routing table entry
type Route struct {
	// Destination is the destination network
	Destination string `json:"destination"`
	// Gateway is the gateway address
	Gateway string `json:"gateway"`
	// Interface is the interface name
	Interface string `json:"interface"`
	// Metric is the route metric
	Metric int `json:"metric,omitempty"`
}

// DNSConfig represents DNS configuration
type DNSConfig struct {
	// Nameservers is a list of nameservers
	Nameservers []string `json:"nameservers"`
	// SearchDomains is a list of search domains
	SearchDomains []string `json:"search_domains,omitempty"`
	// ResolvConf is the content of the resolv.conf file
	ResolvConf string `json:"resolv_conf,omitempty"`
}

// FirewallRule represents a firewall rule
type FirewallRule struct {
	// Chain is the rule chain
	Chain string `json:"chain"`
	// Rule is the rule string
	Rule string `json:"rule"`
	// Target is the rule target
	Target string `json:"target,omitempty"`
	// Protocol is the protocol
	Protocol string `json:"protocol,omitempty"`
	// Source is the source address/network
	Source string `json:"source,omitempty"`
	// Destination is the destination address/network
	Destination string `json:"destination,omitempty"`
}

// ServiceInfo contains information about system services
type ServiceInfo struct {
	// SystemdServices is a list of systemd services
	SystemdServices []Service `json:"systemd_services,omitempty"`
	// OtherServices is a list of non-systemd services
	OtherServices []Service `json:"other_services,omitempty"`
	// SystemdActive is whether systemd is active
	SystemdActive bool `json:"systemd_active,omitempty"`
	// ServiceCount is the total number of services
	ServiceCount int `json:"service_count"`
}

// Service represents a system service
type Service struct {
	// Name is the service name
	Name string `json:"name"`
	// Status is the service status (running, stopped, etc.)
	Status string `json:"status"`
	// Enabled is whether the service is enabled at boot
	Enabled bool `json:"enabled,omitempty"`
	// StartTime is when the service was started
	StartTime time.Time `json:"start_time,omitempty"`
	// Description is the service description
	Description string `json:"description,omitempty"`
	// RecentLogs contains recent log entries
	RecentLogs []LogEntry `json:"recent_logs,omitempty"`
}

// LogEntry represents a log entry
type LogEntry struct {
	// Timestamp is the log entry timestamp
	Timestamp time.Time `json:"timestamp"`
	// Message is the log message
	Message string `json:"message"`
	// Level is the log level (info, warning, error, etc.)
	Level string `json:"level,omitempty"`
	// Source is the log source
	Source string `json:"source,omitempty"`
}

// DockerContainer represents a Docker container
type DockerContainer struct {
	// ID is the container ID
	ID string `json:"id"`
	// Name is the container name
	Name string `json:"name"`
	// Image is the container image
	Image string `json:"image"`
	// Status is the container status
	Status string `json:"status"`
	// Created is when the container was created
	Created time.Time `json:"created"`
	// Ports contains port mappings
	Ports []DockerPort `json:"ports,omitempty"`
	// NetworkMode is the network mode
	NetworkMode string `json:"network_mode,omitempty"`
	// Stats contains container statistics
	Stats *DockerContainerStats `json:"stats,omitempty"`
}

// DockerPort represents a port mapping
type DockerPort struct {
	// PrivatePort is the private port
	PrivatePort int `json:"private_port"`
	// PublicPort is the public port
	PublicPort int `json:"public_port"`
	// Type is the port type (tcp, udp)
	Type string `json:"type"`
}

// DockerContainerStats contains container statistics
type DockerContainerStats struct {
	// CPUPercent is the CPU usage percentage
	CPUPercent float64 `json:"cpu_percent"`
	// MemoryUsage is the memory usage in bytes
	MemoryUsage uint64 `json:"memory_usage"`
	// MemoryLimit is the memory limit in bytes
	MemoryLimit uint64 `json:"memory_limit"`
	// MemoryPercent is the memory usage percentage
	MemoryPercent float64 `json:"memory_percent"`
	// NetworkRx is the network received bytes
	NetworkRx uint64 `json:"network_rx,omitempty"`
	// NetworkTx is the network transmitted bytes
	NetworkTx uint64 `json:"network_tx,omitempty"`
	// BlockRead is the block devices read bytes
	BlockRead uint64 `json:"block_read,omitempty"`
	// BlockWrite is the block devices written bytes
	BlockWrite uint64 `json:"block_write,omitempty"`
}

// DockerImage represents a Docker image
type DockerImage struct {
	// ID is the image ID
	ID string `json:"id"`
	// Repository is the image repository
	Repository string `json:"repository,omitempty"`
	// Tag is the image tag
	Tag string `json:"tag,omitempty"`
	// Size is the image size in bytes
	Size uint64 `json:"size"`
	// Created is when the image was created
	Created time.Time `json:"created"`
}

// DockerVolume represents a Docker volume
type DockerVolume struct {
	// Name is the volume name
	Name string `json:"name"`
	// Driver is the volume driver
	Driver string `json:"driver"`
	// Mountpoint is the volume mountpoint
	Mountpoint string `json:"mountpoint"`
	// Size is the volume size in bytes
	Size uint64 `json:"size,omitempty"`
}

// DockerNetwork represents a Docker network
type DockerNetwork struct {
	// ID is the network ID
	ID string `json:"id"`
	// Name is the network name
	Name string `json:"name"`
	// Driver is the network driver
	Driver string `json:"driver"`
	// Scope is the network scope
	Scope string `json:"scope,omitempty"`
	// Containers is a list of container IDs connected to this network
	Containers []string `json:"containers,omitempty"`
}

// KubernetesInfo contains information about Kubernetes
type KubernetesInfo struct {
	// Version is the Kubernetes version
	Version string `json:"version"`
	// CurrentNamespace is the current namespace
	CurrentNamespace string `json:"current_namespace"`
	// Nodes is a list of nodes
	Nodes []KubernetesNode `json:"nodes,omitempty"`
	// Pods is a list of pods
	Pods []KubernetesPod `json:"pods,omitempty"`
	// Services is a list of services
	Services []KubernetesService `json:"services,omitempty"`
	// Deployments is a list of deployments
	Deployments []KubernetesDeployment `json:"deployments,omitempty"`
	// Events is a list of recent events
	Events []KubernetesEvent `json:"events,omitempty"`
}

// KubernetesNode represents a Kubernetes node
type KubernetesNode struct {
	// Name is the node name
	Name string `json:"name"`
	// Status is the node status
	Status string `json:"status"`
	// Roles are the node roles
	Roles []string `json:"roles,omitempty"`
	// KubeletVersion is the kubelet version
	KubeletVersion string `json:"kubelet_version,omitempty"`
	// Capacity contains the node capacity
	Capacity map[string]string `json:"capacity,omitempty"`
	// Allocatable contains the allocatable resources
	Allocatable map[string]string `json:"allocatable,omitempty"`
}

// KubernetesPod represents a Kubernetes pod
type KubernetesPod struct {
	// Name is the pod name
	Name string `json:"name"`
	// Namespace is the pod namespace
	Namespace string `json:"namespace"`
	// Status is the pod status
	Status string `json:"status"`
	// Node is the node running this pod
	Node string `json:"node,omitempty"`
	// IP is the pod IP
	IP string `json:"ip,omitempty"`
	// Containers is the number of containers
	Containers int `json:"containers"`
	// RestartCount is the number of restarts
	RestartCount int `json:"restart_count,omitempty"`
	// Age is the pod age as a human-readable string
	Age string `json:"age,omitempty"`
}

// KubernetesService represents a Kubernetes service
type KubernetesService struct {
	// Name is the service name
	Name string `json:"name"`
	// Namespace is the service namespace
	Namespace string `json:"namespace"`
	// Type is the service type
	Type string `json:"type"`
	// ClusterIP is the cluster IP
	ClusterIP string `json:"cluster_ip"`
	// ExternalIP is the external IP
	ExternalIP string `json:"external_ip,omitempty"`
	// Ports contains port mappings
	Ports []KubernetesPort `json:"ports,omitempty"`
	// Selector is the service selector
	Selector map[string]string `json:"selector,omitempty"`
}

// KubernetesPort represents a Kubernetes port mapping
type KubernetesPort struct {
	// Name is the port name
	Name string `json:"name,omitempty"`
	// Port is the port number
	Port int `json:"port"`
	// TargetPort is the target port
	TargetPort int `json:"target_port"`
	// Protocol is the port protocol
	Protocol string `json:"protocol,omitempty"`
	// NodePort is the node port
	NodePort int `json:"node_port,omitempty"`
}

// KubernetesDeployment represents a Kubernetes deployment
type KubernetesDeployment struct {
	// Name is the deployment name
	Name string `json:"name"`
	// Namespace is the deployment namespace
	Namespace string `json:"namespace"`
	// Replicas is the number of replicas
	Replicas int `json:"replicas"`
	// ReadyReplicas is the number of ready replicas
	ReadyReplicas int `json:"ready_replicas"`
	// Strategy is the deployment strategy
	Strategy string `json:"strategy,omitempty"`
	// Age is the deployment age as a human-readable string
	Age string `json:"age,omitempty"`
}

// KubernetesEvent represents a Kubernetes event
type KubernetesEvent struct {
	// Type is the event type
	Type string `json:"type"`
	// Reason is the event reason
	Reason string `json:"reason"`
	// Object is the object involved
	Object string `json:"object"`
	// Message is the event message
	Message string `json:"message"`
	// Time is when the event occurred
	Time time.Time `json:"time"`
}

// ProxmoxInfo contains information about Proxmox
type ProxmoxInfo struct {
	// Version is the Proxmox version
	Version string `json:"version,omitempty"`
	// Host is the Proxmox host
	Host string `json:"host"`
	// Status is the cluster status
	Status string `json:"status,omitempty"`
	// Nodes is a list of nodes
	Nodes []ProxmoxNode `json:"nodes,omitempty"`
	// VMs is a list of virtual machines
	VMs []ProxmoxVM `json:"vms,omitempty"`
	// Containers is a list of containers
	Containers []ProxmoxContainer `json:"containers,omitempty"`
	// Storage is a list of storage
	Storage []ProxmoxStorage `json:"storage,omitempty"`
}

// ProxmoxNode represents a Proxmox node
type ProxmoxNode struct {
	// Name is the node name
	Name string `json:"name"`
	// Status is the node status
	Status string `json:"status"`
	// CPUUsage is the CPU usage percentage
	CPUUsage float64 `json:"cpu_usage,omitempty"`
	// MemoryUsage is the memory usage percentage
	MemoryUsage float64 `json:"memory_usage,omitempty"`
	// MemoryTotal is the total memory in bytes
	MemoryTotal uint64 `json:"memory_total,omitempty"`
	// Uptime is the node uptime in seconds
	Uptime uint64 `json:"uptime,omitempty"`
}

// ProxmoxVM represents a Proxmox virtual machine
type ProxmoxVM struct {
	// ID is the VM ID
	ID string `json:"id"`
	// Name is the VM name
	Name string `json:"name"`
	// Status is the VM status
	Status string `json:"status"`
	// Node is the node running this VM
	Node string `json:"node,omitempty"`
	// CPU is the number of CPUs
	CPU int `json:"cpu"`
	// CPUUsage is the CPU usage percentage
	CPUUsage float64 `json:"cpu_usage,omitempty"`
	// Memory is the memory in bytes
	Memory uint64 `json:"memory"`
	// MemoryUsage is the memory usage percentage
	MemoryUsage float64 `json:"memory_usage,omitempty"`
	// Disk is the disk size in bytes
	Disk uint64 `json:"disk,omitempty"`
	// DiskUsage is the disk usage percentage
	DiskUsage float64 `json:"disk_usage,omitempty"`
	// Uptime is the VM uptime in seconds
	Uptime uint64 `json:"uptime,omitempty"`
}

// ProxmoxContainer represents a Proxmox container
type ProxmoxContainer struct {
	// ID is the container ID
	ID string `json:"id"`
	// Name is the container name
	Name string `json:"name"`
	// Status is the container status
	Status string `json:"status"`
	// Node is the node running this container
	Node string `json:"node,omitempty"`
	// CPU is the number of CPUs
	CPU int `json:"cpu"`
	// CPUUsage is the CPU usage percentage
	CPUUsage float64 `json:"cpu_usage,omitempty"`
	// Memory is the memory in bytes
	Memory uint64 `json:"memory"`
	// MemoryUsage is the memory usage percentage
	MemoryUsage float64 `json:"memory_usage,omitempty"`
	// Disk is the disk size in bytes
	Disk uint64 `json:"disk,omitempty"`
	// DiskUsage is the disk usage percentage
	DiskUsage float64 `json:"disk_usage,omitempty"`
	// Uptime is the container uptime in seconds
	Uptime uint64 `json:"uptime,omitempty"`
}

// ProxmoxStorage represents Proxmox storage
type ProxmoxStorage struct {
	// Name is the storage name
	Name string `json:"name"`
	// Type is the storage type
	Type string `json:"type"`
	// Status is the storage status
	Status string `json:"status"`
	// Total is the total size in bytes
	Total uint64 `json:"total"`
	// Used is the used size in bytes
	Used uint64 `json:"used"`
	// UsedPercent is the used percentage
	UsedPercent float64 `json:"used_percent"`
	// Available is the available size in bytes
	Available uint64 `json:"available"`
}

// CephInfo contains information about Ceph
type CephInfo struct {
	// Health is the Ceph health status
	Health string `json:"health"`
	// HealthDetail is a detailed health description
	HealthDetail string `json:"health_detail,omitempty"`
	// MonMap contains monitor information
	MonMap *CephMonMap `json:"mon_map,omitempty"`
	// OSDMap contains OSD information
	OSDMap *CephOSDMap `json:"osd_map,omitempty"`
	// PGMap contains PG information
	PGMap *CephPGMap `json:"pg_map,omitempty"`
	// Pools is a list of pools
	Pools []CephPool `json:"pools,omitempty"`
}

// CephMonMap represents Ceph monitor information
type CephMonMap struct {
	// Monitors is the number of monitors
	Monitors int `json:"monitors"`
	// Quorum is a list of monitors in quorum
	Quorum []string `json:"quorum,omitempty"`
}

// CephOSDMap represents Ceph OSD information
type CephOSDMap struct {
	// NumOSDs is the total number of OSDs
	NumOSDs int `json:"num_osds"`
	// NumUpOSDs is the number of up OSDs
	NumUpOSDs int `json:"num_up_osds"`
	// NumInOSDs is the number of in OSDs
	NumInOSDs int `json:"num_in_osds"`
	// Full is whether the cluster is full
	Full bool `json:"full"`
	// NearFull is whether the cluster is near full
	NearFull bool `json:"near_full"`
}

// CephPGMap represents Ceph PG information
type CephPGMap struct {
	// NumPGs is the number of PGs
	NumPGs int `json:"num_pgs"`
	// BytesTotal is the total bytes
	BytesTotal uint64 `json:"bytes_total"`
	// BytesUsed is the used bytes
	BytesUsed uint64 `json:"bytes_used"`
	// BytesAvailable is the available bytes
	BytesAvailable uint64 `json:"bytes_available"`
}

// CephPool represents a Ceph pool
type CephPool struct {
	// Name is the pool name
	Name string `json:"name"`
	// ID is the pool ID
	ID int `json:"id"`
	// Size is the pool size
	Size int `json:"size,omitempty"`
	// UsedBytes is the used bytes
	UsedBytes uint64 `json:"used_bytes,omitempty"`
	// UsedPercent is the used percentage
	UsedPercent float64 `json:"used_percent,omitempty"`
}

// MetricSummary represents aggregated metrics from various sources
type MetricSummary struct {
	// CPUSummary contains CPU metrics
	CPUSummary *MetricData `json:"cpu_summary,omitempty"`
	// MemorySummary contains memory metrics
	MemorySummary *MetricData `json:"memory_summary,omitempty"`
	// DiskSummary contains disk metrics
	DiskSummary *MetricData `json:"disk_summary,omitempty"`
	// NetworkSummary contains network metrics
	NetworkSummary *MetricData `json:"network_summary,omitempty"`
	// LoadSummary contains load metrics
	LoadSummary *MetricData `json:"load_summary,omitempty"`
}

// MetricData represents a collection of metrics
type MetricData struct {
	// Values is a map of metric names to values
	Values map[string]float64 `json:"values"`
	// History contains historical data points
	History []MetricPoint `json:"history,omitempty"`
	// Min is the minimum value
	Min float64 `json:"min,omitempty"`
	// Max is the maximum value
	Max float64 `json:"max,omitempty"`
	// Avg is the average value
	Avg float64 `json:"avg,omitempty"`
	// Unit is the unit of measurement
	Unit string `json:"unit,omitempty"`
}

// MetricPoint represents a single data point in a time series
type MetricPoint struct {
	// Timestamp is the point timestamp
	Timestamp time.Time `json:"timestamp"`
	// Value is the point value
	Value float64 `json:"value"`
}

// Alert represents an alert or warning detected during collection
type Alert struct {
	// Severity is the alert severity (info, warning, error, critical)
	Severity string `json:"severity"`
	// Source is the alert source
	Source string `json:"source"`
	// Message is the alert message
	Message string `json:"message"`
	// Timestamp is when the alert was detected
	Timestamp time.Time `json:"timestamp"`
	// ResourceID is the affected resource ID
	ResourceID string `json:"resource_id,omitempty"`
	// Metadata contains additional metadata
	Metadata map[string]string `json:"metadata,omitempty"`
}
