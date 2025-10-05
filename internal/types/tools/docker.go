package tools

import (
	"mini-mcp/internal/shared/validation"
)

// DockerComposeUpArgs represents arguments for docker_compose_up.
// Example:
//
//	{"path": "docker-compose.yml", "detached": true, "context": "mycontext", "host": "tcp://remote:2376"}
type DockerComposeUpArgs struct {
	// Path is the path to the docker-compose.yml file
	Path string `json:"path"`
	// Detached indicates whether to run in detached mode
	Detached bool `json:"detached"`
	// Context is the docker context to use
	Context string `json:"context,omitempty"`
	// Host is the docker host to connect to
	Host string `json:"host,omitempty"`
	// DockerPath is the path to the docker binary
	DockerPath string `json:"docker_path,omitempty"`
	// ComposeFile is the path to the compose file
	ComposeFile string `json:"compose_file,omitempty"`
}

// Validate checks if the Docker Compose Up arguments are valid.
func (args *DockerComposeUpArgs) Validate() error {
	// Validate path
	if err := validation.StringRequired("path", args.Path); err != nil {
		return err
	}
	
	// Validate path is safe
	if err := validation.Path("path", args.Path); err != nil {
		return err
	}
	
	return nil
}

// NewDockerComposeUpArgs creates a new DockerComposeUpArgs with the given path and detached flag.
func NewDockerComposeUpArgs(path string, detached bool) *DockerComposeUpArgs {
	return &DockerComposeUpArgs{
		Path:     path,
		Detached: detached,
	}
}

// DockerComposeDownArgs represents arguments for docker_compose_down.
// Example:
//
//	{"path": "docker-compose.yml", "remove_volumes": true, "context": "mycontext", "host": "tcp://remote:2376"}
type DockerComposeDownArgs struct {
	// Path is the path to the docker-compose.yml file
	Path string `json:"path"`
	// RemoveVolumes indicates whether to remove volumes
	RemoveVolumes bool `json:"remove_volumes"`
	// Context is the docker context to use
	Context string `json:"context,omitempty"`
	// Host is the docker host to connect to
	Host string `json:"host,omitempty"`
	// DockerPath is the path to the docker binary
	DockerPath string `json:"docker_path,omitempty"`
	// ComposeFile is the path to the compose file
	ComposeFile string `json:"compose_file,omitempty"`
}

// Validate checks if the Docker Compose Down arguments are valid.
func (args *DockerComposeDownArgs) Validate() error {
	// Validate path
	if err := validation.StringRequired("path", args.Path); err != nil {
		return err
	}
	
	// Validate path is safe
	if err := validation.Path("path", args.Path); err != nil {
		return err
	}
	
	return nil
}

// NewDockerComposeDownArgs creates a new DockerComposeDownArgs with the given path and removeVolumes flag.
func NewDockerComposeDownArgs(path string, removeVolumes bool) *DockerComposeDownArgs {
	return &DockerComposeDownArgs{
		Path:          path,
		RemoveVolumes: removeVolumes,
	}
}

// DockerSwarmInfoArgs represents arguments for docker_swarm_info.
type DockerSwarmInfoArgs struct {
	// Context is the docker context to use
	Context string `json:"context,omitempty"`
	// Host is the docker host to connect to
	Host string `json:"host,omitempty"`
	// DockerPath is the path to the docker binary
	DockerPath string `json:"docker_path,omitempty"`
}

// Validate checks if the Docker Swarm Info arguments are valid.
func (args *DockerSwarmInfoArgs) Validate() error {
	// No validation needed as there are no required fields
	return nil
}

// NewDockerSwarmInfoArgs creates a new DockerSwarmInfoArgs.
func NewDockerSwarmInfoArgs() *DockerSwarmInfoArgs {
	return &DockerSwarmInfoArgs{}
}
