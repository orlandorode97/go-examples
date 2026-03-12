package main

import (
	"context"
	"sql-table/backend/db"
)

type App struct {
	ctx context.Context
	db  *db.Client
}

func NewApp() *App {
	return &App{
		db: db.NewClient(),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

type ConnectionConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

func (a *App) Connect(config ConnectionConfig) error {
	return a.db.Connect(db.Config{
		Host:     config.Host,
		Port:     config.Port,
		Username: config.Username,
		Password: config.Password,
	})
}

func (a *App) Disconnect() error {
	return a.db.Close()
}

func (a *App) IsConnected() bool {
	return a.db.IsConnected()
}

func (a *App) ListDatabases() ([]string, error) {
	return a.db.ListDatabases()
}

func (a *App) ListTables(database string) ([]string, error) {
	return a.db.ListTables(database)
}

func (a *App) ExecuteQuery(query string) (*db.QueryResult, error) {
	return a.db.ExecuteQuery(query)
}
