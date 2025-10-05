package types

// VMBackup represents a VM backup
type VMBackup struct {
	VolID   string `json:"volid"`
	Format  string `json:"format"`
	Size    int64  `json:"size"`
	CTime   int64  `json:"ctime"`
	Path    string `json:"path"`
	Content string `json:"content"`
	VMID    int    `json:"vmid"`
	Notes   string `json:"notes,omitempty"`
}

// VMBackupCreateRequest represents a request to create a VM backup
type VMBackupCreateRequest struct {
	VMID       int    `json:"vmid"`
	Storage    string `json:"storage"`
	Mode       string `json:"mode,omitempty"`     // snapshot, suspend, stop
	Compress   string `json:"compress,omitempty"` // lzo, gzip, zstd
	Remove     bool   `json:"remove,omitempty"`   // Remove old backups
	Notes      string `json:"notes,omitempty"`
	MailTo     string `json:"mailto,omitempty"`
	MailPolicy string `json:"mailpolicy,omitempty"` // always, failure
}

// VMBackupRestoreRequest represents a request to restore from backup
type VMBackupRestoreRequest struct {
	VMID    int    `json:"vmid"`
	Storage string `json:"storage"`
	Backup  string `json:"backup"` // Backup file name
	Force   bool   `json:"force,omitempty"`
	Pool    string `json:"pool,omitempty"`
}
