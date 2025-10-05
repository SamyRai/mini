package resources

// DockerInfo represents Docker system information.
// Example:
//
//	{
//	  "version": "20.10.12",
//	  "containers": 10,
//	  "running": 8,
//	  "paused": 0,
//	  "stopped": 2,
//	  "images": 25,
//	  "server_version": "20.10.12",
//	  "storage_driver": "overlay2",
//	  "logging_driver": "json-file",
//	  "volume_plugins": ["local", "nfs"]
//	}
type DockerInfo struct {
	// Version is the Docker client version
	Version string `json:"version,omitempty"`
	// Containers is the total number of containers
	Containers int `json:"containers,omitempty"`
	// Running is the number of running containers
	Running int `json:"running,omitempty"`
	// Paused is the number of paused containers
	Paused int `json:"paused,omitempty"`
	// Stopped is the number of stopped containers
	Stopped int `json:"stopped,omitempty"`
	// Images is the number of Docker images
	Images int `json:"images,omitempty"`
	// ServerVersion is the Docker server version
	ServerVersion string `json:"server_version,omitempty"`
	// StorageDriver is the storage driver being used
	StorageDriver string `json:"storage_driver,omitempty"`
	// LoggingDriver is the logging driver being used
	LoggingDriver string `json:"logging_driver,omitempty"`
	// VolumePlugins is a list of volume plugins installed
	VolumePlugins []string `json:"volume_plugins,omitempty"`
}

// NewDockerInfo creates a new DockerInfo with the given parameters.
func NewDockerInfo() *DockerInfo {
	return &DockerInfo{}
}

// DockerError represents an error response from Docker.
// Example:
//
//	{"error": "Docker not found on system"}
type DockerError struct {
	// Error is the error message
	Error string `json:"error"`
}

// NewDockerError creates a new DockerError with the given error message.
func NewDockerError(message string) *DockerError {
	return &DockerError{
		Error: message,
	}
}

// DockerRawOutput represents raw output from Docker.
// Example:
//
//	{"raw": "CONTAINER ID   IMAGE     COMMAND   CREATED   STATUS   PORTS   NAMES"}
type DockerRawOutput struct {
	// Raw is the raw output string
	Raw string `json:"raw"`
}

// NewDockerRawOutput creates a new DockerRawOutput with the given raw output.
func NewDockerRawOutput(raw string) *DockerRawOutput {
	return &DockerRawOutput{
		Raw: raw,
	}
}
