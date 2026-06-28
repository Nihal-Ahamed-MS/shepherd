package tools

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"shepherd/core/retriever"
	"shepherd/core/types"

	bleve "github.com/blevesearch/bleve/v2"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SearchInput struct {
	Query string `json:"query"`
	Limit int    `json:"limit,omitempty"`
}

type CodeSearcher struct {
	index  bleve.Index
	chunks []types.Chunk
}

func NewCodeSearcher(chunks []types.Chunk) (*CodeSearcher, error) {
	index, err := retriever.BuildIndex(chunks)
	if err != nil {
		return nil, fmt.Errorf("build index: %w", err)
	}
	return &CodeSearcher{index: index, chunks: chunks}, nil
}

func (cs *CodeSearcher) SearchCodebase(ctx context.Context, req *mcp.CallToolRequest, input SearchInput) (
	*mcp.CallToolResult, any, error,
) {
	limit := input.Limit
	if limit == 0 {
		limit = 10
	}

	query := bleve.NewMatchQuery(input.Query)
	query.Analyzer = "code"
	searchReq := bleve.NewSearchRequestOptions(query, limit, 0, false)

	result, err := cs.index.Search(searchReq)
	if err != nil {
		return nil, nil, fmt.Errorf("search: %w", err)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "Found %d results\n\n", result.Total)

	for _, hit := range result.Hits {
		idx, err := strconv.Atoi(hit.ID)
		if err != nil || idx >= len(cs.chunks) {
			continue
		}
		chunk := cs.chunks[idx]
		filePath := ""
		if len(chunk.FilePath) > 0 {
			filePath = chunk.FilePath[0]
		}
		fmt.Fprintf(&sb, "File: %s (lines %d-%d, score: %.2f)\n%s\n---\n",
			filePath, chunk.StartLine, chunk.EndLine, hit.Score, chunk.SourceCode)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: sb.String()},
		},
	}, nil, nil
}