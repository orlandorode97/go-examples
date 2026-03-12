# AGENTS.md - metrics-go

## Project Overview

metrics-go is a CPU and memory metrics collector that scrapes system metrics and stores them in TimescaleDB (PostgreSQL).

## Build Commands

```bash
# Build the application
go build -o metrics-go main.go

# Run the application (requires TimescaleDB)
go run main.go

# Run with environment variables
DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=password DB_NAME=metricsdb go run main.go

# Download dependencies
go mod download

# Tidy dependencies
go mod tidy
```

## Lint & Code Quality

```bash
# Run go vet
go vet ./...

# Format code
go fmt ./...

# Full build check (vet + compile)
go build -o /dev/null ./...

# Run all checks
go vet ./... && go fmt ./... && go build ./...
```

## Database Commands

```bash
# Start TimescaleDB
docker-compose up -d

# Stop TimescaleDB
docker-compose down

# View logs
docker-compose logs -f timescaledb
```

## Code Style Guidelines

### General

- Go 1.25.6+ required
- Use `gofmt` for formatting (standard Go style)
- No comments unless explaining non-obvious logic

### Imports

- Standard library first, then third-party packages
- Group: stdlib → external → blank import (if needed)
- Use aliased imports only when name conflicts (e.g., `pgxpool "github.com/jackc/pgx/v5/pgxpool"`)

```go
import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/shirou/gopsutil/v4/cpu"
    "github.com/shirou/gopsutil/v4/mem"
)
```

### Naming Conventions

- **Variables/Functions**: camelCase (e.g., `collect`, `flushMem`, `pool`)
- **Types/Interfaces**: PascalCase (e.g., `CPUMetric`, `MemMetric`)
- **Constants**: PascalCase or camelCase for unexported (e.g., `maxRetries`, `DefaultPort`)
- **Acronyms**: Preserve case (e.g., `URL`, `HTTP`, `ID` - not `Url`, `Http`, `Id`)

### Types

- Use specific types: `uint64` for bytes, `float64` for percentages
- Use `time.Time` for timestamps
- Use `context.Context` for cancellation
- Avoid `interface{}` unless necessary; use specific types

### Error Handling

- Return errors with `fmt.Errorf("context: %w", err)` for wrapping
- Log errors and continue where appropriate (non-fatal)
- Use `log.Fatalf` only for startup/fatal errors
- Check errors immediately after calls

```go
// Good
if err != nil {
    return nil, fmt.Errorf("cpu.Percent: %w", err)
}

// Good - non-fatal
if err := flush(ctx, pool, metrics); err != nil {
    log.Printf("flush error: %v", err)
    continue
}
```

### Database Patterns

- Use connection pooling via `pgxpool.New`
- Use `COPY` protocol for bulk inserts (`pool.CopyFrom`)
- Close pools with `defer pool.Close()`
- Use parameterized queries (automatic with pgx)

### Testing

```bash
# Run all tests
go test ./...

# Run specific test
go test -run TestCollect ./...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...
```

### Docker

```bash
# Build image
docker build -t metrics-go .

# Run container
docker run -e DB_HOST=timescaledb -e DB_PASSWORD=password metrics-go

# Run with docker-compose (includes DB)
docker-compose up --build
```

### Project Structure

```
.
├── main.go           # Application entry point
├── go.mod            # Dependencies
├── go.sum            # Lock file
├── init.sql          # Database schema
├── Dockerfile        # Container build
├── docker-compose.yaml # Local dev environment
└── AGENTS.md         # This file
```

### Database Schema Conventions

- Table names: `snake_case` (e.g., `cpu_metrics`, `mem_metrics`)
- Columns: `snake_case` (e.g., `usage_percent`, `available_bytes`)
- Use TimescaleDB hypertables for time-series tables
- Always include `time` column as primary partition key
- Add index on `(host, time DESC)` for common queries
