// Package model defines the core analysis data structures.
package model

// Core domain types for the analyzer. Keep this package dumb, with no deps on
// go/ast or go/packages to maintain clean boundaries and single ownership.

// Analysis is the top-level result structure produced by the analyzer.
type Analysis struct {
	PackageDependencies map[string][]string       `json:"packageDependencies"`
	InterfaceUsage      map[string][]string       `json:"interfaceUsage"`
	StructComplexity    map[string]int            `json:"structComplexity"`
	LayerViolations     []string                  `json:"layerViolations"`
	CyclicDependencies  [][]string                `json:"cyclicDependencies"`
	CodeMetrics         map[string]PackageMetrics `json:"codeMetrics"`
}

// PackageMetrics captures basic metrics per package.
type PackageMetrics struct {
	LinesOfCode          int     `json:"linesOfCode"`
	CyclomaticComplexity int     `json:"cyclomaticComplexity"`
	Dependencies         int     `json:"dependencies"`
	Abstractions         int     `json:"abstractions"`
	ConcreteTypes        int     `json:"concreteTypes"`
	AbstractionRatio     float64 `json:"abstractionRatio"`
}

// Graph models a simple language graph that can be extended later.
// Nodes are package directories (absolute paths). Edges are package deps.
type Graph struct {
	// Adj maps a node to its outbound neighbors.
	Adj map[string][]string `json:"adj"`
}
