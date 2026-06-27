package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SearchInput struct {
	Query string `json:"query"`
	Limit int    `json:"limit,omitempty"`
}

func SearchCodebase(ctx context.Context, req *mcp.CallToolRequest, input SearchInput) (
	*mcp.CallToolResult, any, error,
) {
	result := "result...."

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result},
		},
	}, nil, nil
}
