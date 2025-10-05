// Package util contains helper functions for metrics and collections.
package util

import (
	"go/ast"
	"strings"
)

func CalculateStructComplexity(st *ast.StructType) int {
	if st == nil || st.Fields == nil {
		return 0
	}
	complexity := 0
	for _, field := range st.Fields.List {
		// count named fields; embedded fields have Names==nil -> count as 1
		if len(field.Names) == 0 {
			complexity++
			continue
		}
		complexity += len(field.Names)
	}
	return complexity
}

func CalculateFunctionComplexity(fn *ast.FuncDecl) int {
	if fn == nil {
		return 0
	}
	complexity := 1 // base path
	ast.Inspect(fn, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.CaseClause, *ast.CommClause:
			complexity++
		}
		return true
	})
	return complexity
}

func CalculateLinesOfCode(src string) int {
	if src == "" {
		return 0
	}
	return len(strings.Split(src, "\n"))
}

func UniqStrings(in []string) []string {
	if len(in) <= 1 {
		return in
	}
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, s := range in {
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}
