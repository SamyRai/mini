// Package resources contains types for MCP resource responses.
package resources

import "time"

// SystemInfo represents system information.
// Example:
//
//	{
//	  "os": "linux",
//	  "arch": "amd64",
//	  "cpus": 8,
//	  "hostname": "server1",
//	  "memory": "16384 MB",
//	  "kernel": "5.15.0-58-generic",
//	  "uptime": "5 days, 3 hours",
//	  "boot_time": "2023-01-01T00:00:00Z",
//	  "load_average": [1.5, 1.2, 0.8],
//	  "cpu_details": {...},
//	  "memory_details": {...},
//	  "disk_info": [...],
//	  "network_interfaces": [...]
//	}
type SystemInfo struct {
	// OS is the operating system name (e.g., "linux", "darwin", "windows")
	OS string `json:"os"`
	// Arch is the system architecture (e.g., "amd64", "arm64")
	Arch string `json:"arch"`
	// CPUs is the number of CPU cores
	CPUs int `json:"cpus"`
	// Hostname is the system hostname
	Hostname string `json:"hostname,omitempty"`
	// Memory is a string representation of the system memory
	Memory string `json:"memory,omitempty"`
	// Kernel is the kernel version
	Kernel string `json:"kernel,omitempty"`
	// Uptime is the system uptime as a human-readable string
	Uptime string `json:"uptime,omitempty"`
	// BootTime is when the system was last booted
	BootTime time.Time `json:"boot_time,omitempty"`
	// LoadAverage is the 1, 5, and 15-minute load averages
	LoadAverage []float64 `json:"load_average,omitempty"`
	// CPUDetails provides detailed information about the CPU
	CPUDetails *CPUDetails `json:"cpu_details,omitempty"`
	// MemoryDetails provides detailed information about memory usage
	MemoryDetails *MemoryDetails `json:"memory_details,omitempty"`
	// DiskInfo provides information about disk usage
	DiskInfo []DiskInfo `json:"disk_info,omitempty"`
	// NetworkInterfaces provides information about network interfaces
	NetworkInterfaces []NetworkInterface `json:"network_interfaces,omitempty"`
	// CloudInfo provides information about the cloud provider if applicable
	CloudInfo *CloudInfo `json:"cloud_info,omitempty"`
}

// CPUDetails represents detailed information about CPU
type CPUDetails struct {
	// Model is the CPU model name
	Model string `json:"model,omitempty"`
	// Vendor is the CPU vendor
	Vendor string `json:"vendor,omitempty"`
	// Speed is the CPU speed in MHz
	SpeedMHz float64 `json:"speed_mhz,omitempty"`
	// PhysicalCores is the number of physical CPU cores
	PhysicalCores int `json:"physical_cores,omitempty"`
	// LogicalCores is the number of logical CPU cores
	LogicalCores int `json:"logical_cores,omitempty"`
	// Usage is the current CPU usage percentage
	Usage float64 `json:"usage,omitempty"`
	// PerCoreUsage is the usage percentage per core
	PerCoreUsage []float64 `json:"per_core_usage,omitempty"`
}

// MemoryDetails represents detailed information about memory
type MemoryDetails struct {
	// Total is the total memory in bytes
	Total uint64 `json:"total"`
	// Free is the free memory in bytes
	Free uint64 `json:"free"`
	// Used is the used memory in bytes
	Used uint64 `json:"used"`
	// UsedPercent is the percentage of memory used
	UsedPercent float64 `json:"used_percent"`
	// Cached is the cached memory in bytes
	Cached uint64 `json:"cached,omitempty"`
	// Buffers is the buffer memory in bytes
	Buffers uint64 `json:"buffers,omitempty"`
	// Swap is information about swap memory
	Swap *SwapMemory `json:"swap,omitempty"`
}

// SwapMemory represents information about swap memory
type SwapMemory struct {
	// Total is the total swap memory in bytes
	Total uint64 `json:"total"`
	// Free is the free swap memory in bytes
	Free uint64 `json:"free"`
	// Used is the used swap memory in bytes
	Used uint64 `json:"used"`
	// UsedPercent is the percentage of swap memory used
	UsedPercent float64 `json:"used_percent"`
}

// DiskInfo represents information about disk usage
type DiskInfo struct {
	// Path is the mount path
	Path string `json:"path"`
	// Device is the device name
	Device string `json:"device,omitempty"`
	// Total is the total disk space in bytes
	Total uint64 `json:"total"`
	// Free is the free disk space in bytes
	Free uint64 `json:"free"`
	// Used is the used disk space in bytes
	Used uint64 `json:"used"`
	// UsedPercent is the percentage of disk space used
	UsedPercent float64 `json:"used_percent"`
	// FileSystem is the filesystem type
	FileSystem string `json:"file_system,omitempty"`
	// IOStats provides I/O statistics if available
	IOStats *DiskIOStats `json:"io_stats,omitempty"`
}

// DiskIOStats represents disk I/O statistics
type DiskIOStats struct {
	// ReadCount is the number of read operations
	ReadCount uint64 `json:"read_count"`
	// WriteCount is the number of write operations
	WriteCount uint64 `json:"write_count"`
	// ReadBytes is the number of bytes read
	ReadBytes uint64 `json:"read_bytes"`
	// WriteBytes is the number of bytes written
	WriteBytes uint64 `json:"write_bytes"`
	// ReadTime is the time spent reading in milliseconds
	ReadTime uint64 `json:"read_time,omitempty"`
	// WriteTime is the time spent writing in milliseconds
	WriteTime uint64 `json:"write_time,omitempty"`
}

// NetworkInterface represents information about a network interface
type NetworkInterface struct {
	// Name is the interface name
	Name string `json:"name"`
	// HardwareAddr is the MAC address
	HardwareAddr string `json:"hardware_addr,omitempty"`
	// Addresses contains IP addresses assigned to the interface
	Addresses []string `json:"addresses,omitempty"`
	// BytesSent is the number of bytes sent
	BytesSent uint64 `json:"bytes_sent"`
	// BytesRecv is the number of bytes received
	BytesRecv uint64 `json:"bytes_recv"`
	// PacketsSent is the number of packets sent
	PacketsSent uint64 `json:"packets_sent,omitempty"`
	// PacketsRecv is the number of packets received
	PacketsRecv uint64 `json:"packets_recv,omitempty"`
	// MTU is the Maximum Transmission Unit
	MTU int `json:"mtu,omitempty"`
	// Flags are the interface flags
	Flags []string `json:"flags,omitempty"`
}

// CloudInfo represents information about the cloud provider if the system is running in the cloud
type CloudInfo struct {
	// Provider is the cloud provider name (AWS, GCP, Azure, etc.)
	Provider string `json:"provider,omitempty"`
	// Region is the cloud region
	Region string `json:"region,omitempty"`
	// InstanceType is the instance type or size
	InstanceType string `json:"instance_type,omitempty"`
	// InstanceID is the cloud instance ID
	InstanceID string `json:"instance_id,omitempty"`
	// AvailabilityZone is the availability zone
	AvailabilityZone string `json:"availability_zone,omitempty"`
	// PublicIP is the public IP address
	PublicIP string `json:"public_ip,omitempty"`
	// PrivateIP is the private IP address
	PrivateIP string `json:"private_ip,omitempty"`
}

// NewSystemInfo creates a new SystemInfo with the given parameters.
func NewSystemInfo(os, arch string, cpus int, hostname, memory string) *SystemInfo {
	return &SystemInfo{
		OS:       os,
		Arch:     arch,
		CPUs:     cpus,
		Hostname: hostname,
		Memory:   memory,
	}
}

// NewDetailedSystemInfo creates a new SystemInfo with detailed information.
func NewDetailedSystemInfo(
	os, arch string, cpus int, hostname, memory, kernel, uptime string,
	bootTime time.Time, loadAverage []float64,
	cpuDetails *CPUDetails, memDetails *MemoryDetails,
	diskInfo []DiskInfo, netInterfaces []NetworkInterface,
	cloudInfo *CloudInfo,
) *SystemInfo {
	return &SystemInfo{
		OS:                os,
		Arch:              arch,
		CPUs:              cpus,
		Hostname:          hostname,
		Memory:            memory,
		Kernel:            kernel,
		Uptime:            uptime,
		BootTime:          bootTime,
		LoadAverage:       loadAverage,
		CPUDetails:        cpuDetails,
		MemoryDetails:     memDetails,
		DiskInfo:          diskInfo,
		NetworkInterfaces: netInterfaces,
		CloudInfo:         cloudInfo,
	}
}
