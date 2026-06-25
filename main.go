package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SearchInput struct {
	Query string `json:"query"`
	Limit int    `json:"limit,omitempty"`
}

func sampleToolCall(ctx context.Context, req *mcp.CallToolRequest, input SearchInput) (
	*mcp.CallToolResult, any, error,
) {
	result := "result...."

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result},
		},
	}, nil, nil
}

func main() {
	log.Println("Starting shepherd....")

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "Shepherd",
		Version: "0.0.1",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_codebase",
		Description: "Semantic + BM25 hybrid search over inexed files",
	}, sampleToolCall)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Println("Error:", err)
	}
}
