package main

import (
	"context"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	client := mcp.NewClient(&mcp.Implementation{
		Name:    "mcp-mysql",
		Version: "0.0.1",
	}, nil)

	session, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: "http://localhost:5005"}, nil)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer session.Close()

	result, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "query",
		Arguments: map[string]any{
			"sql": "SELECT * FROM `order` LIMIT 1",
		},
	})

	if err != nil {
		log.Printf("Failed to get query: %s\n", err.Error())
		return
	}

	// Print the result.
	for _, content := range result.Content {
		if textContent, ok := content.(*mcp.TextContent); ok {
			fmt.Println(textContent.Text)
		}
	}

}
