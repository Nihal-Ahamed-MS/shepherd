package main

import (
	"log"

	"shepherd/core/tools"
	"shepherd/core/types"
)

func main() {
	log.Println("Starting shepherd....")

	tools.IndexTesting(types.IndexCodebaseInput{
		Path: "/Users/nihalahamed/Desktop/kimchi-cli",
	})

	// server := mcp.NewServer(&mcp.Implementation{
	// 	Name:    "Shepherd",
	// 	Version: "0.0.1",
	// }, nil)

	// mcp.AddTool(server, &mcp.Tool{
	// 	Name:        "search_codebase",
	// 	Description: "Semantic + BM25 hybrid search over inexed files",
	// }, tools.SearchCodebase)

	// mcp.AddTool(server, &mcp.Tool{
	// 	Name:        "index_codebase",
	// 	Description: "Returns the absolute path of the given codebase directory",
	// }, tools.IndexCodebase)

	// if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
	// 	log.Println("Error:", err)
	// }
}
