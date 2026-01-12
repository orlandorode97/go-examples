package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	_ "github.com/go-sql-driver/mysql"
)

type Row map[string]any

type QueryInput struct {
	SQL string `json:"sql"`
}

type QueryOutput struct {
	Rows []Row `json:"rows"`
}

func Query(ctx context.Context, req *mcp.CallToolRequest, input QueryInput) (
	*mcp.CallToolResult, QueryOutput, error) {
	db, _ := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/database")

	defer db.Close()

	rows, _ := db.QueryContext(ctx, input.SQL)
	defer rows.Close()

	cols, _ := rows.Columns()
	results := make([]Row, 0)

	for rows.Next() {
		values := make([]any, len(cols))
		pointers := make([]any, len(cols))

		for i := range values {
			pointers[i] = &values[i]
		}
		if err := rows.Scan(pointers...); err != nil {
			return nil, QueryOutput{}, fmt.Errorf("failed to scan row: %w", err)
		}
		row := make(Row)
		for i, col := range cols {
			val := values[i]
			if b, ok := val.([]byte); ok {
				val = string(b)
			}
			row[col] = val
		}
		results = append(results, row)
	}

	return nil, QueryOutput{Rows: results}, nil
}

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "mcp-mysql",
		Version: "0.0.1",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "query",
		Description: "Run a SQL query on the local MySQL database",
	}, Query)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}

}
