// Package loader loads Go packages using go/packages with typed syntax.
package loader

import (
	"context"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/packages"
)

// Package represents a loaded Go package with syntax and file info.
type Package struct {
	ID         string
	PkgPath    string
	Dir        string
	GoFiles    []string
	Fset       *token.FileSet
	Syntax     []*ast.File
	ModuleDir  string // absolute dir of the main module root
	ModulePath string // module path of the main module (e.g., example.com/mod)
}

// LoadAll loads packages under the given directory with the provided pattern (use "./..." to recurse).
// If pattern is empty, it defaults to "./...".
func LoadAll(ctx context.Context, dir string, pattern string) ([]*Package, error) {
	if pattern == "" {
		pattern = "./..."
	}
	cfg := &packages.Config{
		Mode: packages.LoadSyntax | packages.NeedModule,
		Dir:  dir,
	}
	pkgs, err := packages.Load(cfg, pattern)
	if err != nil {
		return nil, err
	}
	// Find main module root if available (assume first main module encountered)
	var moduleDir string
	var modulePath string
	for _, p := range pkgs {
		if p.Module != nil && p.Module.Main {
			moduleDir = p.Module.Dir
			modulePath = p.Module.Path
			break
		}
	}
	out := make([]*Package, 0, len(pkgs))
	for _, p := range pkgs {
		out = append(out, &Package{
			ID:         p.ID,
			PkgPath:    p.PkgPath,
			Dir:        p.Dir,
			GoFiles:    p.CompiledGoFiles,
			Fset:       p.Fset,
			Syntax:     p.Syntax,
			ModuleDir:  moduleDir,
			ModulePath: modulePath,
		})
	}
	return out, nil
}
