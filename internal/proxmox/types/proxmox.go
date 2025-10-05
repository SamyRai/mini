package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Node represents a Proxmox node
type Node struct {
	Node   string `json:"node"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Level  string `json:"level"`
	ID     string `json:"id"`
}

// NodeStatus represents the status of a Proxmox node
type NodeStatus struct {
	Node    string    `json:"node"`
	Status  string    `json:"status"`
	Uptime  int64     `json:"uptime"`
	CPU     float64   `json:"cpu"`
	Memory  int64     `json:"memory"`
	MaxMem  int64     `json:"maxmem"`
	Disk    int64     `json:"disk"`
	MaxDisk int64     `json:"maxdisk"`
	LoadAvg []float64 `json:"loadavg"`
}

// VM represents a virtual machine
type VM struct {
	VMID   int    `json:"vmid"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Node   string `json:"node"`
	Type   string `json:"type"`
}

// VMStatus represents the status of a virtual machine
type VMStatus struct {
	VMID      int     `json:"vmid"`
	Status    string  `json:"status"`
	Uptime    int64   `json:"uptime"`
	CPU       float64 `json:"cpu"`
	Memory    int64   `json:"memory"`
	MaxMem    int64   `json:"maxmem"`
	Disk      int64   `json:"disk"`
	MaxDisk   int64   `json:"maxdisk"`
	NetIn     int64   `json:"netin"`
	NetOut    int64   `json:"netout"`
	DiskRead  int64   `json:"diskread"`
	DiskWrite int64   `json:"diskwrite"`
}

// VMConfig represents the configuration of a virtual machine
type VMConfig struct {
	VMID         int               `json:"vmid"`
	Name         string            `json:"name"`
	Memory       int               `json:"memory"`
	Cores        int               `json:"cores"`
	OSType       string            `json:"ostype"`
	SCSIHw       string            `json:"scsihw"`
	Net0         string            `json:"net0"`
	Agent        string            `json:"agent"`
	VGA          string            `json:"vga"`
	Boot         string            `json:"boot"`
	BootDisk     string            `json:"bootdisk"`
	OnBoot       bool              `json:"onboot"`
	Start        bool              `json:"start"`
	IPConfig0    string            `json:"ipconfig0,omitempty"`
	Nameserver   string            `json:"nameserver,omitempty"`
	SearchDomain string            `json:"searchdomain,omitempty"`
	Extra        map[string]string `json:"-"`
}

// UnmarshalJSON implements custom JSON unmarshaling for VMConfig with type safety
func (v *VMConfig) UnmarshalJSON(data []byte) error {
	// Create a temporary struct with flexible field types
	type Alias VMConfig
	aux := &struct {
		Memory interface{} `json:"memory"`
		OnBoot interface{} `json:"onboot"`
		Start  interface{} `json:"start"`
		*Alias
	}{
		Alias: (*Alias)(v),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Handle memory field with type safety
	if aux.Memory != nil {
		if err := v.setMemory(aux.Memory); err != nil {
			return fmt.Errorf("invalid memory value: %w", err)
		}
	}

	// Handle onboot field with type safety
	if aux.OnBoot != nil {
		if err := v.setOnBoot(aux.OnBoot); err != nil {
			return fmt.Errorf("invalid onboot value: %w", err)
		}
	}

	// Handle start field with type safety
	if aux.Start != nil {
		if err := v.setStart(aux.Start); err != nil {
			return fmt.Errorf("invalid start value: %w", err)
		}
	}

	return nil
}

// setMemory safely sets the memory value with type checking
func (v *VMConfig) setMemory(value interface{}) error {
	switch mem := value.(type) {
	case string:
		if mem != "" {
			if parsed, err := strconv.Atoi(mem); err == nil {
				v.Memory = parsed
			} else {
				return fmt.Errorf("cannot parse memory string: %s", mem)
			}
		}
	case float64:
		v.Memory = int(mem)
	case int:
		v.Memory = mem
	default:
		return fmt.Errorf("unsupported memory type: %T", value)
	}
	return nil
}

// setOnBoot safely sets the onboot value with type checking
func (v *VMConfig) setOnBoot(value interface{}) error {
	switch onboot := value.(type) {
	case bool:
		v.OnBoot = onboot
	case float64:
		v.OnBoot = onboot != 0
	case int:
		v.OnBoot = onboot != 0
	case string:
		v.OnBoot = onboot == "1" || onboot == "true"
	default:
		return fmt.Errorf("unsupported onboot type: %T", value)
	}
	return nil
}

// setStart safely sets the start value with type checking
func (v *VMConfig) setStart(value interface{}) error {
	switch start := value.(type) {
	case bool:
		v.Start = start
	case float64:
		v.Start = start != 0
	case int:
		v.Start = start != 0
	case string:
		v.Start = start == "1" || start == "true"
	default:
		return fmt.Errorf("unsupported start type: %T", value)
	}
	return nil
}

// VMCreateRequest represents a request to create a VM
type VMCreateRequest struct {
	VMID         int    `json:"vmid"`
	Name         string `json:"name"`
	Memory       int    `json:"memory"`
	Cores        int    `json:"cores"`
	OSType       string `json:"ostype"`
	SCSIHw       string `json:"scsihw"`
	Net0         string `json:"net0"`
	Agent        string `json:"agent"`
	VGA          string `json:"vga"`
	Boot         string `json:"boot"`
	BootDisk     string `json:"bootdisk"`
	OnBoot       bool   `json:"onboot"`
	Start        bool   `json:"start"`
	Bios         string `json:"bios,omitempty"`    // BIOS type (ovmf for UEFI)
	Serial0      string `json:"serial0,omitempty"` // Serial console
	IPConfig0    string `json:"ipconfig0,omitempty"`
	Nameserver   string `json:"nameserver,omitempty"`
	SearchDomain string `json:"searchdomain,omitempty"`
	// Storage and ISO configuration
	Storage string `json:"storage,omitempty"`
	IDE0    string `json:"ide0,omitempty"`  // ISO image for boot
	SATA0   string `json:"sata0,omitempty"` // Primary disk
	SCSI0   string `json:"scsi0,omitempty"` // Alternative disk
	// Cloud-init configuration
	CIUser     string `json:"ciuser,omitempty"`
	CIPassword string `json:"cipassword,omitempty"`
	CIUserData string `json:"ciuserdata,omitempty"`
	CINetwork  string `json:"cinetwork,omitempty"`
	CIDNS      string `json:"cidns,omitempty"`
	CISearch   string `json:"cisearch,omitempty"`
	Hostname   string `json:"hostname,omitempty"`
	SSHKeys    string `json:"sshkeys,omitempty"`
}

// Storage represents a Proxmox storage
type Storage struct {
	Storage string `json:"storage"`
	Type    string `json:"type"`
	Content string `json:"content"`
	Active  int    `json:"active"`
	Used    int64  `json:"used"`
	Avail   int64  `json:"avail"`
	Total   int64  `json:"total"`
}

// StorageContent represents content in a storage
type StorageContent struct {
	VolID   string `json:"volid"`
	Format  string `json:"format"`
	Size    int64  `json:"size"`
	CTime   int64  `json:"ctime"`
	Path    string `json:"path"`
	Content string `json:"content"`
}

// Image represents a VM image
type Image struct {
	VolID   string    `json:"volid"`
	Format  string    `json:"format"`
	Size    int64     `json:"size"`
	CTime   time.Time `json:"ctime"`
	Path    string    `json:"path"`
	Content string    `json:"content"`
	Name    string    `json:"name"`
	Storage string    `json:"storage"`
}

// UbuntuDeployConfig represents configuration for Ubuntu deployment
type UbuntuDeployConfig struct {
	Name       string   `json:"name"`
	IP         string   `json:"ip"`
	Cores      int      `json:"cores"`
	Memory     int      `json:"memory"`
	Disk       int      `json:"disk"`
	Version    string   `json:"version"` // 22.04, 24.04
	CIPassword string   `json:"cipassword,omitempty"`
	SSHKeys    []string `json:"ssh_keys,omitempty"`
	UserData   string   `json:"user_data,omitempty"`
	CloudInit  bool     `json:"cloud_init"`
}

// TemplateDeployConfig represents configuration for template deployment
type TemplateDeployConfig struct {
	Name     string `json:"name"`
	IP       string `json:"ip"`
	Template int    `json:"template"`
	Full     bool   `json:"full"`
}

// CustomDeployConfig represents configuration for custom deployment
type CustomDeployConfig struct {
	Name     string            `json:"name"`
	IP       string            `json:"ip"`
	Cores    int               `json:"cores"`
	Memory   int               `json:"memory"`
	Disk     int               `json:"disk"`
	Config   map[string]string `json:"config"`
	Image    string            `json:"image,omitempty"`
	Template int               `json:"template,omitempty"`
}

// AuthConfig represents authentication configuration
type AuthConfig struct {
	Proxmox struct {
		Host       string `yaml:"host" json:"host"`
		User       string `yaml:"user" json:"user"`
		Password   string `yaml:"password" json:"password"`
		TokenName  string `yaml:"token_name" json:"token_name"`
		TokenValue string `yaml:"token_value" json:"token_value"`
		VerifySSL  bool   `yaml:"verify_ssl" json:"verify_ssl"`
		Timeout    int    `yaml:"timeout" json:"timeout"`
		Node       string `yaml:"node" json:"node"`
	} `yaml:"proxmox" json:"proxmox"`
	// Optional global post-install command to run on first boot (one-line)
	PostInstall string `yaml:"post_install" json:"post_install"`
}

// VMConfigItem represents a VM configuration item in MainConfig
type VMConfigItem struct {
	ID     int    `yaml:"id" json:"id"`
	Name   string `yaml:"name" json:"name"`
	IP     string `yaml:"ip" json:"ip"`
	Cores  int    `yaml:"cores" json:"cores"`
	Memory int    `yaml:"memory" json:"memory"`
	Disk   int    `yaml:"disk" json:"disk"`
}

// ImageSource represents an image source configuration
type ImageSource struct {
	Name        string `yaml:"name" json:"name"`
	URL         string `yaml:"url" json:"url"`
	Version     string `yaml:"version" json:"version"`
	Description string `yaml:"description" json:"description"`
	Type        string `yaml:"type" json:"type"` // ubuntu, centos, debian, etc.
}

// MainConfig represents the main application configuration
type MainConfig struct {
	VMs  []VMConfigItem `yaml:"vms" json:"vms"`
	User struct {
		Name     string `yaml:"name" json:"name"`
		Password string `yaml:"password" json:"password"`
	} `yaml:"user" json:"user"`
	SSHKeys []string `yaml:"ssh_keys" json:"ssh_keys"`
	Network struct {
		Gateway    string `yaml:"gateway" json:"gateway"`
		Netmask    int    `yaml:"netmask" json:"netmask"`
		Nameserver string `yaml:"nameserver" json:"nameserver"`
		Bridge     string `yaml:"bridge" json:"bridge"`
	} `yaml:"network" json:"network"`
	TemplateID int `yaml:"template_id" json:"template_id"`
	Images     struct {
		Sources        []ImageSource `yaml:"sources" json:"sources"`
		DefaultStorage string        `yaml:"default_storage" json:"default_storage"`
	} `yaml:"images" json:"images"`
	Proxmox struct {
		Storage struct {
			Local         string `yaml:"local" json:"local"`
			ISOStorage    string `yaml:"iso_storage" json:"iso_storage"`
			SnippetsStore string `yaml:"snippets_storage" json:"snippets_storage"`
		} `yaml:"storage" json:"storage"`
		Network struct {
			Bridge string `yaml:"bridge" json:"bridge"`
			Model  string `yaml:"model" json:"model"`
		} `yaml:"network" json:"network"`
		VMDefaults struct {
			ScsiController string `yaml:"scsi_controller" json:"scsi_controller"`
			VGA            string `yaml:"vga" json:"vga"`
			Serial         string `yaml:"serial" json:"serial"`
			BootOrder      string `yaml:"boot_order" json:"boot_order"`
			BootDisk       string `yaml:"boot_disk" json:"boot_disk"`
			Agent          string `yaml:"agent" json:"agent"`
		} `yaml:"vm_defaults" json:"vm_defaults"`
	} `yaml:"proxmox" json:"proxmox"`
}

// NetworkInterface represents a network interface configuration
type NetworkInterface struct {
	Iface      string            `json:"iface"`
	Type       string            `json:"type"`
	Bridge     string            `json:"bridge,omitempty"`
	VLAN       int               `json:"vlan,omitempty"`
	Address    string            `json:"address,omitempty"`
	Netmask    string            `json:"netmask,omitempty"`
	Gateway    string            `json:"gateway,omitempty"`
	BondMode   string            `json:"bond_mode,omitempty"`
	BondSlaves string            `json:"bond_slaves,omitempty"`
	Comments   string            `json:"comments,omitempty"`
	Extra      map[string]string `json:"-"`
}

// NetworkInterfaceConfig represents configuration for creating/updating network interfaces
type NetworkInterfaceConfig struct {
	Type       string            `json:"type"` // bridge, bond, vlan, etc.
	Bridge     string            `json:"bridge,omitempty"`
	VLAN       int               `json:"vlan,omitempty"`
	Address    string            `json:"address,omitempty"`
	Netmask    string            `json:"netmask,omitempty"`
	Gateway    string            `json:"gateway,omitempty"`
	BondMode   string            `json:"bond_mode,omitempty"`
	BondSlaves string            `json:"bond_slaves,omitempty"`
	Comments   string            `json:"comments,omitempty"`
	Extra      map[string]string `json:"-"`
}

// FirewallRule represents a firewall rule
type FirewallRule struct {
	Pos     int    `json:"pos"`
	Action  string `json:"action"` // ACCEPT, DROP, REJECT
	Type    string `json:"type"`   // in, out, fwbr
	Source  string `json:"source,omitempty"`
	Dest    string `json:"dest,omitempty"`
	Proto   string `json:"proto,omitempty"`  // tcp, udp, icmp, etc.
	Sport   string `json:"sport,omitempty"`  // source port
	Dport   string `json:"dport,omitempty"`  // destination port
	Macro   string `json:"macro,omitempty"`  // DNAT, SNAT, etc.
	Target  string `json:"target,omitempty"` // target for DNAT/SNAT
	Log     string `json:"log,omitempty"`    // log level
	Comment string `json:"comment,omitempty"`
	Enable  int    `json:"enable"` // 0 or 1
}

// FirewallOptions represents firewall options
type FirewallOptions struct {
	Enable          int    `json:"enable"` // 0 or 1
	LogLevel        string `json:"log_level,omitempty"`
	LogRateLimit    string `json:"log_rate_limit,omitempty"`
	PolicyIn        string `json:"policy_in,omitempty"`  // ACCEPT, DROP, REJECT
	PolicyOut       string `json:"policy_out,omitempty"` // ACCEPT, DROP, REJECT
	PolicyFwd       string `json:"policy_fwd,omitempty"` // ACCEPT, DROP, REJECT
	SmurfProtection int    `json:"smurf_protection"`     // 0 or 1
	TCPFlags        string `json:"tcp_flags,omitempty"`
	NDP             int    `json:"ndp"`      // 0 or 1
	IPFilter        int    `json:"ipfilter"` // 0 or 1
	Radv            int    `json:"radv"`     // 0 or 1
}

// FirewallAlias represents a firewall alias (IP/network group)
type FirewallAlias struct {
	Name    string `json:"name"`
	CIDR    string `json:"cidr,omitempty"`
	Comment string `json:"comment,omitempty"`
}

// FirewallGroup represents a firewall group
type FirewallGroup struct {
	Group   string `json:"group"`
	Comment string `json:"comment,omitempty"`
}

// CLIError represents a CLI error
type CLIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *CLIError) Error() string {
	return e.Message
}

// NewCLIError creates a new CLI error
func NewCLIError(code int, message string, details ...string) *CLIError {
	err := &CLIError{
		Code:    code,
		Message: message,
	}
	if len(details) > 0 {
		err.Details = details[0]
	}
	return err
}
