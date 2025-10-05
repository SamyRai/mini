package resources

// CommandDoc represents documentation for a command.
// Example:
//
//	{
//	  "description": "Docker container management",
//	  "usage": "docker [OPTIONS] COMMAND",
//	  "url": "https://docs.docker.com/engine/reference/commandline/cli/"
//	}
type CommandDoc struct {
	// Description is a brief description of the command
	Description string `json:"description"`
	// Usage is the command usage syntax
	Usage string `json:"usage"`
	// URL is a link to the command's documentation
	URL string `json:"url"`
}

// NewCommandDoc creates a new CommandDoc with the given parameters.
func NewCommandDoc(description, usage, url string) *CommandDoc {
	return &CommandDoc{
		Description: description,
		Usage:       usage,
		URL:         url,
	}
}

// CommandDocCollection represents a collection of command documentation.
// Example:
//
//	{
//	  "docker": {
//	    "description": "Docker container management",
//	    "usage": "docker [OPTIONS] COMMAND",
//	    "url": "https://docs.docker.com/engine/reference/commandline/cli/"
//	  },
//	  "git": {
//	    "description": "Distributed version control system",
//	    "usage": "git [--version] [--help] [-C <path>] [-c <name>=<value>]",
//	    "url": "https://git-scm.com/docs"
//	  }
//	}
type CommandDocCollection map[string]CommandDoc

// NewCommandDocCollection creates a new empty CommandDocCollection.
func NewCommandDocCollection() CommandDocCollection {
	return make(CommandDocCollection)
}

// Add adds a command documentation to the collection.
func (c CommandDocCollection) Add(name string, doc CommandDoc) {
	c[name] = doc
}

// Get retrieves a command documentation from the collection.
func (c CommandDocCollection) Get(name string) (CommandDoc, bool) {
	doc, ok := c[name]
	return doc, ok
}
