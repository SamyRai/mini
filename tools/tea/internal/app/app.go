// Package app composes the analyzer pipeline.
package app

import (
	"context"

	"mini-mcp/tools/tea/internal/analyzer"
	"mini-mcp/tools/tea/internal/graph"
	"mini-mcp/tools/tea/internal/loader"
	"mini-mcp/tools/tea/internal/model"
	"mini-mcp/tools/tea/internal/report"
)

// Run executes the analysis pipeline for the given root dir and writes output to outPath.
func Run(ctx context.Context, rootDir, outPath string) (*model.Analysis, error) {
    pkgs, err := loader.LoadAll(ctx, rootDir, "./...")
    if err != nil {
        return nil, err
    }
    a := analyzer.Analyze(pkgs)

    // Build cycles
    a.CyclicDependencies = graph.FindCycles(a.PackageDependencies)

    if err := report.Write(outPath, a); err != nil {
        return nil, err
    }
    return a, nil
}
