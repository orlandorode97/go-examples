package main

import (
	"context"
	"log"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	ctx := context.Background()
	client := mcp.NewClient(
		&mcp.Implementation{
			Name:    "mcp-client",
			Version: "v1.0.0",
		}, nil)

	transport := &mcp.CommandTransport{
		Command: exec.Command("go", "run", "mcp_server.go"),
	}

	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer session.Close()

	params := &mcp.CallToolParams{
		Name:      "greet",
		Arguments: map[string]any{"name": "you"},
	}

	res, err := session.CallTool(ctx, params)
	if err != nil {
		log.Fatal(err)
	}

	if res.IsError {
		log.Fatal("tool failed")
	}

	for _, c := range res.Content {
		log.Printf(c.(*mcp.TextContent).Text)
	}
}
