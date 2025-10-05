package types

// VMConsole represents VM console access information
type VMConsole struct {
	User       string `json:"user"`
	Ticket     string `json:"ticket"`
	Port       string `json:"port,omitempty"`
	VNCProxy   string `json:"vncproxy,omitempty"`
	SPICEProxy string `json:"spiceproxy,omitempty"`
	TermProxy  string `json:"termproxy,omitempty"`
	WebSocket  string `json:"websocket,omitempty"`
	Cert       string `json:"cert,omitempty"`
	// SPICE-specific fields
	Host             string `json:"host,omitempty"`
	TLSPort          int    `json:"tls-port,omitempty"`
	Password         string `json:"password,omitempty"`
	Type             string `json:"type,omitempty"`
	CA               string `json:"ca,omitempty"`
	Proxy            string `json:"proxy,omitempty"`
	Title            string `json:"title,omitempty"`
	HostSubject      string `json:"host-subject,omitempty"`
	DeleteFile       int    `json:"delete-this-file,omitempty"`
	ReleaseCursor    string `json:"release-cursor,omitempty"`
	SecureAttention  string `json:"secure-attention,omitempty"`
	ToggleFullscreen string `json:"toggle-fullscreen,omitempty"`
}

// VMConsoleRequest represents a request to get VM console access
type VMConsoleRequest struct {
	VMID        int    `json:"vmid"`
	WebSocket   bool   `json:"websocket,omitempty"`    // Use WebSocket instead of VNC
	ConsoleType string `json:"console_type,omitempty"` // "vnc", "spice", "serial" - auto-detected if not specified
}
