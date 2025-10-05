package types

import "time"

// VMPerformanceData represents VM performance metrics
type VMPerformanceData struct {
	VMID      int     `json:"vmid"`
	CPU       float64 `json:"cpu"`
	Memory    float64 `json:"memory"`
	MaxMem    float64 `json:"maxmem"`
	Disk      float64 `json:"disk"`
	MaxDisk   float64 `json:"maxdisk"`
	NetIn     float64 `json:"netin"`
	NetOut    float64 `json:"netout"`
	DiskRead  float64 `json:"diskread"`
	DiskWrite float64 `json:"diskwrite"`
	Uptime    int64   `json:"uptime"`
	Time      int64   `json:"time"`
}

// VMLogEntry represents a VM log entry
type VMLogEntry struct {
	Time    time.Time `json:"time"`
	Level   string    `json:"level"`
	Message string    `json:"message"`
	PID     int       `json:"pid,omitempty"`
	User    string    `json:"user,omitempty"`
}

// VMEvent represents a VM event
type VMEvent struct {
	Time    time.Time `json:"time"`
	Type    string    `json:"type"`
	Message string    `json:"message"`
	User    string    `json:"user,omitempty"`
	VMID    int       `json:"vmid,omitempty"`
}

// VMStatistics represents VM statistics request parameters
type VMStatistics struct {
	VMID      int    `json:"vmid"`
	TimeFrame string `json:"timeframe,omitempty"` // hour, day, week, month, year
	Start     string `json:"start,omitempty"`     // Start time
	End       string `json:"end,omitempty"`       // End time
}
