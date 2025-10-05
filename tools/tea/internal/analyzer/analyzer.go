// Package analyzer transforms loaded packages into an analysis model.
package analyzer

import (
	"fmt"
	"go/ast"
	"path/filepath"
	"strings"

	"mini-mcp/tools/tea/internal/loader"
	"mini-mcp/tools/tea/internal/model"
	"mini-mcp/tools/tea/internal/util"
)

// Analyze builds the analysis model from loaded packages.
func Analyze(pkgs []*loader.Package) *model.Analysis {
    a := &model.Analysis{
        PackageDependencies: make(map[string][]string),
        InterfaceUsage:      make(map[string][]string),
        StructComplexity:    make(map[string]int),
        LayerViolations:     make([]string, 0),
        CyclicDependencies:  make([][]string, 0),
        CodeMetrics:         make(map[string]model.PackageMetrics),
    }

    var modulePath, moduleDir string
    if len(pkgs) > 0 {
        modulePath = pkgs[0].ModulePath
        moduleDir = pkgs[0].ModuleDir
    }

    for _, p := range pkgs {
        metrics := model.PackageMetrics{}

        // Gather imports and AST metrics
        for i, f := range p.Syntax {
            // Calculate LOC using source positions only if files available; otherwise skip.
            if i < len(p.GoFiles) {
                // We can't read file contents here without IO; keep LOC minimal by counting AST lines as 0.
                _ = p.GoFiles[i] // Use the file path to avoid unused variable warning
            }
            ast.Inspect(f, func(n ast.Node) bool {
                switch x := n.(type) {
                case *ast.ImportSpec:
                    importPath := strings.Trim(x.Path.Value, "\"")
                    metrics.Dependencies++
                    if modulePath != "" && strings.HasPrefix(importPath, modulePath+"/") {
                        rel := strings.TrimPrefix(importPath, modulePath+"/")
                        depDir := filepath.Join(moduleDir, filepath.FromSlash(rel))
                        a.PackageDependencies[p.Dir] = append(a.PackageDependencies[p.Dir], depDir)
                    }
                case *ast.TypeSpec:
                    switch t := x.Type.(type) {
                    case *ast.InterfaceType:
                        metrics.Abstractions++
                        key := fmt.Sprintf("%s:%s", p.Dir, x.Name.Name)
                        a.InterfaceUsage[key] = extractInterfaceMethods(t)
                    case *ast.StructType:
                        metrics.ConcreteTypes++
                        complexity := util.CalculateStructComplexity(t)
                        fullName := fmt.Sprintf("%s.%s", p.Dir, x.Name.Name)
                        a.StructComplexity[fullName] = complexity
                    }
                case *ast.FuncDecl:
                    metrics.CyclomaticComplexity += util.CalculateFunctionComplexity(x)
                }
                return true
            })
        }

        // Abstraction ratio
        if metrics.ConcreteTypes > 0 {
            metrics.AbstractionRatio = float64(metrics.Abstractions) / float64(metrics.ConcreteTypes)
        }

        a.CodeMetrics[p.Dir] = metrics
    }

    // De-duplicate deps
    for k, v := range a.PackageDependencies {
        a.PackageDependencies[k] = util.UniqStrings(v)
    }

    return a
}

// extractInterfaceMethods returns a slice of method names (or embedded types) for an interface
func extractInterfaceMethods(it *ast.InterfaceType) []string {
    methods := make([]string, 0)
    if it == nil || it.Methods == nil {
        return methods
    }
    for _, f := range it.Methods.List {
        if len(f.Names) == 0 {
            methods = append(methods, exprToString(f.Type))
            continue
        }
        for _, name := range f.Names {
            methods = append(methods, name.Name)
        }
    }
    return methods
}

// exprToString renders a minimal representation of an AST expression
func exprToString(e ast.Expr) string {
    switch v := e.(type) {
    case *ast.Ident:
        return v.Name
    case *ast.SelectorExpr:
        return exprToString(v.X) + "." + v.Sel.Name
    default:
        return fmt.Sprintf("%T", e)
    }
}
