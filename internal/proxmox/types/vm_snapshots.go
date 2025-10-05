package types

import "time"

// VMSnapshot represents a VM snapshot
type VMSnapshot struct {
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Time        time.Time `json:"time"`
	VMState     bool      `json:"vmstate"` // Whether to include VM state
	Parent      string    `json:"parent,omitempty"`
	Children    []string  `json:"children,omitempty"`
}

// VMSnapshotCreateRequest represents a request to create a VM snapshot
type VMSnapshotCreateRequest struct {
	VMID        int    `json:"vmid"`
	SnapName    string `json:"snapname"`
	Description string `json:"description,omitempty"`
	VMState     bool   `json:"vmstate"` // Whether to include VM state
}

// VMSnapshotRollbackRequest represents a request to rollback to a snapshot
type VMSnapshotRollbackRequest struct {
	VMID     int    `json:"vmid"`
	SnapName string `json:"snapname"`
	Start    bool   `json:"start,omitempty"` // Whether to start VM after rollback
}
