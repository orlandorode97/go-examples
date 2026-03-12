package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

type Client struct {
	db *sql.DB
}

type QueryResult struct {
	Columns []string
	Rows    [][]interface{}
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Connect(cfg Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	c.db = db
	return nil
}

func (c *Client) IsConnected() bool {
	return c.db != nil
}

func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

func (c *Client) ListDatabases() ([]string, error) {
	if c.db == nil {
		return nil, fmt.Errorf("not connected")
	}

	rows, err := c.db.Query("SHOW DATABASES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return nil, err
		}
		databases = append(databases, dbName)
	}

	return databases, nil
}

func (c *Client) ListTables(database string) ([]string, error) {
	if c.db == nil {
		return nil, fmt.Errorf("not connected")
	}

	_, err := c.db.Exec(fmt.Sprintf("USE %s", database))
	if err != nil {
		return nil, err
	}

	rows, err := c.db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}

	return tables, nil
}

func (c *Client) ExecuteQuery(query string) (*QueryResult, error) {
	if c.db == nil {
		return nil, fmt.Errorf("not connected")
	}

	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	var resultRows [][]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		var row []interface{}
		for i, col := range values {
			if col == nil {
				row = append(row, nil)
			} else {
				switch columnTypes[i].DatabaseTypeName() {
				case "DECIMAL", "NUMERIC":
					row = append(row, fmt.Sprintf("%v", col))
				default:
					row = append(row, col)
				}
			}
		}
		resultRows = append(resultRows, row)
	}

	return &QueryResult{
		Columns: columns,
		Rows:    resultRows,
	}, nil
}
