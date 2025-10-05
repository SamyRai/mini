package tools

import (
	"mini-mcp/internal/shared/validation"
)

// GitCloneArgs represents arguments for git_clone.
// Example:
//
//	{"repo": "https://github.com/user/repo.git", "path": "./repo", "branch": "main"}
type GitCloneArgs struct {
	// Repo is the Git repository URL or path
	Repo string `json:"repo"`
	// Path is the local path where the repository will be cloned
	Path string `json:"path"`
	// Branch is the branch to clone (optional)
	Branch string `json:"branch,omitempty"`
}

// Validate checks if the Git clone arguments are valid.
func (args *GitCloneArgs) Validate() error {
	// Validate repository URL
	if err := validation.StringRequired("repo", args.Repo); err != nil {
		return err
	}
	
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

// NewGitCloneArgs creates a new GitCloneArgs with the given repository, path, and optional branch.
func NewGitCloneArgs(repo, path, branch string) *GitCloneArgs {
	return &GitCloneArgs{
		Repo:   repo,
		Path:   path,
		Branch: branch,
	}
}
