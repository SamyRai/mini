package types

// VMTemplate represents a VM template
type VMTemplate struct {
	VMID        int    `json:"vmid"`
	Name        string `json:"name"`
	Template    int    `json:"template"`
	Description string `json:"description,omitempty"`
	Pool        string `json:"pool,omitempty"`
}

// VMTemplateCreateRequest represents a request to create a VM template
type VMTemplateCreateRequest struct {
	VMID        int    `json:"vmid"`
	Description string `json:"description,omitempty"`
}

// VMTemplateDeployRequest represents a request to deploy from template
type VMTemplateDeployRequest struct {
	TemplateID int    `json:"template"`
	NewID      int    `json:"newid"`
	Name       string `json:"name,omitempty"`
	Pool       string `json:"pool,omitempty"`
	Storage    string `json:"storage,omitempty"`
	Target     string `json:"target,omitempty"` // Target node
	Full       bool   `json:"full"`             // Full clone
}
