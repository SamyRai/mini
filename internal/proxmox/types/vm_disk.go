package types

// VMDisk represents a VM disk configuration
type VMDisk struct {
	ID       string `json:"id"`                 // e.g., "scsi0", "sata0", "virtio0"
	Storage  string `json:"storage"`            // Storage name
	Size     string `json:"size"`               // Size in GB or format like "32G"
	Format   string `json:"format,omitempty"`   // Disk format (raw, qcow2, vmdk)
	Cache    string `json:"cache,omitempty"`    // Cache mode (none, writethrough, writeback)
	Discard  string `json:"discard,omitempty"`  // Discard mode (on, off, ignore)
	SSD      bool   `json:"ssd,omitempty"`      // SSD emulation
	IOThread bool   `json:"iothread,omitempty"` // IO thread
}

// VMDiskAddRequest represents a request to add a disk to a VM
type VMDiskAddRequest struct {
	VMID    int    `json:"vmid"`
	ID      string `json:"id"` // Disk ID (e.g., "scsi1", "sata1")
	Storage string `json:"storage"`
	Size    string `json:"size"`
	Format  string `json:"format,omitempty"`
	Cache   string `json:"cache,omitempty"`
	SSD     bool   `json:"ssd,omitempty"`
}

// VMDiskResizeRequest represents a request to resize a VM disk
type VMDiskResizeRequest struct {
	VMID int    `json:"vmid"`
	ID   string `json:"id"`   // Disk ID
	Size string `json:"size"` // New size
}

// VMDiskMoveRequest represents a request to move a VM disk
type VMDiskMoveRequest struct {
	VMID    int    `json:"vmid"`
	ID      string `json:"id"`               // Disk ID
	Storage string `json:"storage"`          // Target storage
	Delete  bool   `json:"delete,omitempty"` // Delete source after move
	Online  bool   `json:"online,omitempty"` // Online move
}
