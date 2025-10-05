package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	"mini-mcp/tools/tea/internal/app"
)

func main() {
	rootPath := flag.String("path", ".", "Root path of the project to analyze")
	outputPath := flag.String("output", "analysis.json", "Output file path for the analysis")
	flag.Parse()

	rootAbs, _ := filepath.Abs(*rootPath)
	_, err := app.Run(context.Background(), rootAbs, *outputPath)
	if err != nil {
		fmt.Printf("tea: error: %v\n", err)
		return
	}
	fmt.Printf("Analysis completed and written to %s\n", *outputPath)
}
