package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"shepherd/core/chunking"
	"shepherd/core/types"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func IndexCodebase(ctx context.Context, req *mcp.CallToolRequest, input types.IndexCodebaseInput) (
	*mcp.CallToolResult, any, error,
) {
	absPath, err := filepath.Abs(input.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to resolve path: %w", err)
	}

	if _, err := os.Stat(absPath); err != nil {
		return nil, nil, fmt.Errorf("path does not exist: %s", absPath)
	}

	codebase, err := chunking.ParseCodebase(absPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse codebase: %w", err)
	}

	out, err := json.Marshal(codebase)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to serialize AST: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(out)},
		},
	}, nil, nil
}

func IndexTesting(input types.IndexCodebaseInput) {
	absPath, _ := filepath.Abs(input.Path)
	log.Println(absPath)
	absPath, err := filepath.Abs(input.Path)
	if err != nil {
		return
	}

	if _, err := os.Stat(absPath); err != nil {
		return
	}

	chunking.ParseCodebase(absPath)

}
