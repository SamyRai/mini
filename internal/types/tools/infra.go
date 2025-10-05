package tools

import (
	"mini-mcp/internal/shared/validation"
)

// ProxmoxArgs represents arguments for proxmox_status.
// Example:
//
//	{"host": "proxmox.example.com", "user": "root", "password": "secret"}
type ProxmoxArgs struct {
	// Host is the hostname or IP address of the Proxmox server
	Host string `json:"host"`
	// User is the username for Proxmox authentication
	User string `json:"user"`
	// Password is the password for Proxmox authentication
	Password string `json:"password"`
}

// Validate checks if the Proxmox arguments are valid.
func (args *ProxmoxArgs) Validate() error {
	result := validation.ValidateProxmoxArgs(args.Host, args.User, args.Password)
	if !result.IsValid {
		return result.Errors[0] // Return first error for backward compatibility
	}
	return nil
}

// NewProxmoxArgs creates a new ProxmoxArgs with the given host, user, and password.
func NewProxmoxArgs(host, user, password string) *ProxmoxArgs {
	return &ProxmoxArgs{
		Host:     host,
		User:     user,
		Password: password,
	}
}

// CephArgs represents arguments for ceph_status.
type CephArgs struct {
	// CephArgs is intentionally empty as the ceph_status tool doesn't require
	// arguments. The struct exists to provide type safety and future extensibility.
}

// Validate checks if the Ceph arguments are valid.
func (args *CephArgs) Validate() error {
	// No validation needed as there are no required fields
	return nil
}

// NewCephArgs creates a new CephArgs.
func NewCephArgs() *CephArgs {
	return &CephArgs{}
}

// SystemInfoArgs represents arguments for system_info tool.
// Example:
//
//	{"detailed": true, "include_cloud": true}
type SystemInfoArgs struct {
	// Detailed indicates whether to include detailed system information
	Detailed bool `json:"detailed"`
	// IncludeProcesses indicates whether to include process information
	IncludeProcesses bool `json:"include_processes"`
	// IncludeCloud indicates whether to attempt cloud provider detection
	IncludeCloud bool `json:"include_cloud"`
	// IncludeNetwork indicates whether to include detailed network information
	IncludeNetwork bool `json:"include_network"`
	// IncludeIOStats indicates whether to include disk I/O statistics
	IncludeIOStats bool `json:"include_io_stats"`
}

// Validate checks if the SystemInfo arguments are valid.
func (args *SystemInfoArgs) Validate() error {
	// No validation needed as all fields have default values
	return nil
}

// NewSystemInfoArgs creates a new SystemInfoArgs with the given parameters.
func NewSystemInfoArgs(detailed, includeProcesses, includeCloud, includeNetwork, includeIOStats bool) *SystemInfoArgs {
	return &SystemInfoArgs{
		Detailed:         detailed,
		IncludeProcesses: includeProcesses,
		IncludeCloud:     includeCloud,
		IncludeNetwork:   includeNetwork,
		IncludeIOStats:   includeIOStats,
	}
}

// DockerInfoArgs represents arguments for docker_info tool.
// Example:
//
//	{"include_containers": true, "include_stats": true}
type DockerInfoArgs struct {
	// IncludeContainers indicates whether to include container information
	IncludeContainers bool `json:"include_containers"`
	// IncludeStats indicates whether to include container statistics
	IncludeStats bool `json:"include_stats"`
	// IncludeImages indicates whether to include image information
	IncludeImages bool `json:"include_images"`
	// IncludeVolumes indicates whether to include volume information
	IncludeVolumes bool `json:"include_volumes"`
	// IncludeNetworks indicates whether to include network information
	IncludeNetworks bool `json:"include_networks"`
}

// Validate checks if the DockerInfo arguments are valid.
func (args *DockerInfoArgs) Validate() error {
	// No validation needed as all fields have default values
	return nil
}

// NewDockerInfoArgs creates a new DockerInfoArgs with the given parameters.
func NewDockerInfoArgs(includeContainers, includeStats, includeImages, includeVolumes, includeNetworks bool) *DockerInfoArgs {
	return &DockerInfoArgs{
		IncludeContainers: includeContainers,
		IncludeStats:      includeStats,
		IncludeImages:     includeImages,
		IncludeVolumes:    includeVolumes,
		IncludeNetworks:   includeNetworks,
	}
}

// KubernetesInfoArgs represents arguments for kubernetes_info tool.
// Example:
//
//	{"namespace": "default", "include_pods": true, "include_services": true}
type KubernetesInfoArgs struct {
	// Namespace is the Kubernetes namespace to query
	Namespace string `json:"namespace"`
	// IncludePods indicates whether to include pod information
	IncludePods bool `json:"include_pods"`
	// IncludeServices indicates whether to include service information
	IncludeServices bool `json:"include_services"`
	// IncludeDeployments indicates whether to include deployment information
	IncludeDeployments bool `json:"include_deployments"`
	// IncludeNodes indicates whether to include node information
	IncludeNodes bool `json:"include_nodes"`
	// IncludeEvents indicates whether to include recent events
	IncludeEvents bool `json:"include_events"`
	// KubeconfigPath is the path to the kubeconfig file
	KubeconfigPath string `json:"kubeconfig_path,omitempty"`
}

// Validate checks if the KubernetesInfo arguments are valid.
func (args *KubernetesInfoArgs) Validate() error {
	// Namespace is optional, defaults to "default" if not specified
	if args.Namespace == "" {
		args.Namespace = "default"
	}
	return nil
}

// NewKubernetesInfoArgs creates a new KubernetesInfoArgs with the given parameters.
func NewKubernetesInfoArgs(namespace string, includePods, includeServices, includeDeployments, includeNodes, includeEvents bool, kubeconfigPath string) *KubernetesInfoArgs {
	return &KubernetesInfoArgs{
		Namespace:          namespace,
		IncludePods:        includePods,
		IncludeServices:    includeServices,
		IncludeDeployments: includeDeployments,
		IncludeNodes:       includeNodes,
		IncludeEvents:      includeEvents,
		KubeconfigPath:     kubeconfigPath,
	}
}

// CloudInfoArgs represents arguments for cloud_info tool.
// Example:
//
//	{"provider": "aws", "region": "us-east-1"}
type CloudInfoArgs struct {
	// Provider is the cloud provider (aws, gcp, azure, etc.)
	Provider string `json:"provider,omitempty"`
	// Region is the cloud region
	Region string `json:"region,omitempty"`
	// IncludeInstances indicates whether to include instance information
	IncludeInstances bool `json:"include_instances"`
	// IncludeStorages indicates whether to include storage information
	IncludeStorages bool `json:"include_storages"`
	// IncludeNetworking indicates whether to include networking information
	IncludeNetworking bool `json:"include_networking"`
	// IncludeDatabases indicates whether to include database information
	IncludeDatabases bool `json:"include_databases"`
	// AwsAccessKey is the AWS access key
	AwsAccessKey string `json:"aws_access_key,omitempty"`
	// AwsSecretKey is the AWS secret key
	AwsSecretKey string `json:"aws_secret_key,omitempty"`
	// AzureClientId is the Azure client ID
	AzureClientId string `json:"azure_client_id,omitempty"`
	// AzureClientSecret is the Azure client secret
	AzureClientSecret string `json:"azure_client_secret,omitempty"`
	// AzureTenantId is the Azure tenant ID
	AzureTenantId string `json:"azure_tenant_id,omitempty"`
	// GcpCredentialsPath is the path to GCP credentials file
	GcpCredentialsPath string `json:"gcp_credentials_path,omitempty"`
}

// Validate checks if the CloudInfo arguments are valid.
func (args *CloudInfoArgs) Validate() error {
	// Provider is optional, auto-detection will be attempted if not specified
	return nil
}

// NewCloudInfoArgs creates a new CloudInfoArgs with the given parameters.
func NewCloudInfoArgs(provider, region string, includeInstances, includeStorages, includeNetworking, includeDatabases bool) *CloudInfoArgs {
	return &CloudInfoArgs{
		Provider:          provider,
		Region:            region,
		IncludeInstances:  includeInstances,
		IncludeStorages:   includeStorages,
		IncludeNetworking: includeNetworking,
		IncludeDatabases:  includeDatabases,
	}
}

// ProcessInfoArgs represents arguments for process_info tool.
// Example:
//
//	{"sort": "cpu", "limit": 10}
type ProcessInfoArgs struct {
	// Sort is the field to sort processes by (cpu, memory, io, name)
	Sort string `json:"sort,omitempty"`
	// Limit is the maximum number of processes to return
	Limit int `json:"limit,omitempty"`
	// FilterUser filters processes by user
	FilterUser string `json:"filter_user,omitempty"`
	// FilterName filters processes by name
	FilterName string `json:"filter_name,omitempty"`
	// IncludeCPU indicates whether to include CPU usage
	IncludeCPU bool `json:"include_cpu"`
	// IncludeMemory indicates whether to include memory usage
	IncludeMemory bool `json:"include_memory"`
	// IncludeIO indicates whether to include I/O statistics
	IncludeIO bool `json:"include_io"`
}

// Validate checks if the ProcessInfo arguments are valid.
func (args *ProcessInfoArgs) Validate() error {
	// Set defaults if not specified
	if args.Limit <= 0 {
		args.Limit = 10
	}
	if args.Sort == "" {
		args.Sort = "cpu"
	}
	return nil
}

// NewProcessInfoArgs creates a new ProcessInfoArgs with the given parameters.
func NewProcessInfoArgs(sort string, limit int, filterUser, filterName string, includeCPU, includeMemory, includeIO bool) *ProcessInfoArgs {
	return &ProcessInfoArgs{
		Sort:          sort,
		Limit:         limit,
		FilterUser:    filterUser,
		FilterName:    filterName,
		IncludeCPU:    includeCPU,
		IncludeMemory: includeMemory,
		IncludeIO:     includeIO,
	}
}

// NetworkInfoArgs represents arguments for network_info tool.
// Example:
//
//	{"include_connections": true, "include_routing": true}
type NetworkInfoArgs struct {
	// IncludeInterfaces indicates whether to include interface information
	IncludeInterfaces bool `json:"include_interfaces"`
	// IncludeConnections indicates whether to include connection information
	IncludeConnections bool `json:"include_connections"`
	// IncludeRouting indicates whether to include routing information
	IncludeRouting bool `json:"include_routing"`
	// IncludeDNS indicates whether to include DNS information
	IncludeDNS bool `json:"include_dns"`
	// IncludeFirewall indicates whether to include firewall rules
	IncludeFirewall bool `json:"include_firewall"`
	// FilterInterface filters by interface name
	FilterInterface string `json:"filter_interface,omitempty"`
}

// Validate checks if the NetworkInfo arguments are valid.
func (args *NetworkInfoArgs) Validate() error {
	// No validation needed as all fields have default values
	return nil
}

// NewNetworkInfoArgs creates a new NetworkInfoArgs with the given parameters.
func NewNetworkInfoArgs(includeInterfaces, includeConnections, includeRouting, includeDNS, includeFirewall bool, filterInterface string) *NetworkInfoArgs {
	return &NetworkInfoArgs{
		IncludeInterfaces:  includeInterfaces,
		IncludeConnections: includeConnections,
		IncludeRouting:     includeRouting,
		IncludeDNS:         includeDNS,
		IncludeFirewall:    includeFirewall,
		FilterInterface:    filterInterface,
	}
}

// ServiceInfoArgs represents arguments for service_info tool.
// Example:
//
//	{"include_systemd": true, "include_logs": true}
type ServiceInfoArgs struct {
	// IncludeSystemd indicates whether to include systemd services
	IncludeSystemd bool `json:"include_systemd"`
	// IncludeLogs indicates whether to include recent logs
	IncludeLogs bool `json:"include_logs"`
	// IncludeStatus indicates whether to include service status
	IncludeStatus bool `json:"include_status"`
	// FilterService filters by service name
	FilterService string `json:"filter_service,omitempty"`
	// LogLines is the number of log lines to include (if IncludeLogs is true)
	LogLines int `json:"log_lines,omitempty"`
}

// Validate checks if the ServiceInfo arguments are valid.
func (args *ServiceInfoArgs) Validate() error {
	// Set defaults if not specified
	if args.LogLines <= 0 {
		args.LogLines = 20
	}
	return nil
}

// NewServiceInfoArgs creates a new ServiceInfoArgs with the given parameters.
func NewServiceInfoArgs(includeSystemd, includeLogs, includeStatus bool, filterService string, logLines int) *ServiceInfoArgs {
	return &ServiceInfoArgs{
		IncludeSystemd: includeSystemd,
		IncludeLogs:    includeLogs,
		IncludeStatus:  includeStatus,
		FilterService:  filterService,
		LogLines:       logLines,
	}
}

// InfrastructureInfoArgs represents arguments for infrastructure_info tool which collects comprehensive information.
// Example:
//
//	{"collection_level": "detailed", "include_all": true}
type InfrastructureInfoArgs struct {
	// CollectionLevel determines how much detail to collect (basic, standard, detailed, comprehensive)
	CollectionLevel string `json:"collection_level,omitempty"`

	// Include flags for different types of information
	IncludeSystem     bool `json:"include_system"`
	IncludeProcesses  bool `json:"include_processes"`
	IncludeNetwork    bool `json:"include_network"`
	IncludeServices   bool `json:"include_services"`
	IncludeDocker     bool `json:"include_docker"`
	IncludeKubernetes bool `json:"include_kubernetes"`
	IncludeCloud      bool `json:"include_cloud"`
	IncludeProxmox    bool `json:"include_proxmox"`
	IncludeCeph       bool `json:"include_ceph"`
	IncludeNomad      bool `json:"include_nomad"`

	// IncludeAll is a shortcut to include all available information
	IncludeAll bool `json:"include_all"`

	// TimeInterval is the time interval in seconds for metrics collection
	TimeInterval int `json:"time_interval,omitempty"`

	// Authentication details for various services
	ProxmoxAuth    *ProxmoxArgs        `json:"proxmox_auth,omitempty"`
	KubernetesAuth *KubernetesInfoArgs `json:"kubernetes_auth,omitempty"`
	CloudAuth      *CloudInfoArgs      `json:"cloud_auth,omitempty"`
}

// Validate checks if the InfrastructureInfo arguments are valid.
func (args *InfrastructureInfoArgs) Validate() error {
	// Set defaults if not specified
	if args.CollectionLevel == "" {
		args.CollectionLevel = "standard"
	}

	// If IncludeAll is set, enable all include flags
	if args.IncludeAll {
		args.IncludeSystem = true
		args.IncludeProcesses = true
		args.IncludeNetwork = true
		args.IncludeServices = true
		args.IncludeDocker = true
		args.IncludeKubernetes = true
		args.IncludeCloud = true
		args.IncludeProxmox = true
		args.IncludeCeph = true
		args.IncludeNomad = true
	}

	// Default to at least include system info
	if !args.IncludeSystem && !args.IncludeProcesses && !args.IncludeNetwork &&
		!args.IncludeServices && !args.IncludeDocker && !args.IncludeKubernetes &&
		!args.IncludeCloud && !args.IncludeProxmox && !args.IncludeCeph && !args.IncludeNomad {
		args.IncludeSystem = true
	}

	// Set default time interval
	if args.TimeInterval <= 0 {
		args.TimeInterval = 5
	}

	return nil
}

// NewInfrastructureInfoArgs creates a new InfrastructureInfoArgs with sensible defaults.
func NewInfrastructureInfoArgs() *InfrastructureInfoArgs {
	return &InfrastructureInfoArgs{
		CollectionLevel:   "standard",
		IncludeSystem:     true,
		IncludeProcesses:  true,
		IncludeNetwork:    true,
		IncludeServices:   false,
		IncludeDocker:     false,
		IncludeKubernetes: false,
		IncludeCloud:      false,
		IncludeProxmox:    false,
		IncludeCeph:       false,
		IncludeNomad:      false,
		IncludeAll:        false,
		TimeInterval:      5,
	}
}
