#!/bin/bash

cd /Users/damirmukimov/projects/mcp/mini

# Compile the tea analyzer
go build -o tools/tea/tea tools/tea/main.go

# Run the analysis on the project
tools/tea/tea -path . -output tools/tea/analysis.json

# Format the JSON output for better readability
if command -v jq &> /dev/null; then
    jq '.' tools/tea/analysis.json > tools/tea/analysis_pretty.json
    mv tools/tea/analysis_pretty.json tools/tea/analysis.json
fi
