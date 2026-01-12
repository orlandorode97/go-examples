package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Input struct {
	Name string `json:"name" jsonchema:"the name of the person to greet"`
}

type Output struct {
	Greeting string `json:"greeting" jsonchema:"the greeting to tell to the user"`
}

func SayHi(ctx context.Context, req *mcp.CallToolRequest, input Input) (*mcp.CallToolResult, Output, error) {
	return nil, Output{
		Greeting: "Hi, " + input.Name,
	}, nil
}

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name: "greeter",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "greet",
		Description: "Say Hi",
	}, SayHi)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
