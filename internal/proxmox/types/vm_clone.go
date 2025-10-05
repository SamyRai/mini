package types

// VMCloneRequest represents a request to clone a VM
type VMCloneRequest struct {
	NewID       int    `json:"newid"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Full        bool   `json:"full"` // Full clone (not linked)
	Pool        string `json:"pool,omitempty"`
	Storage     string `json:"storage,omitempty"`
	Target      string `json:"target,omitempty"` // Target node
}

// VMMigrateRequest represents a request to migrate a VM
type VMMigrateRequest struct {
	Target    string `json:"target"`               // Target node
	Online    bool   `json:"online,omitempty"`     // Online migration
	WithLocal bool   `json:"with-local,omitempty"` // Migrate local disks
	Force     bool   `json:"force,omitempty"`      // Force migration
}
