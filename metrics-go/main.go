package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shirou/gopsutil/v4/cpu"
)

// CPUMetric holds a single CPU reading for one core.
type CPUMetric struct {
	Time          time.Time
	Host          string
	Core          string
	UsagePercent  float64
	UserPercent   float64
	SystemPercent float64
	IdlePercent   float64
}

// collect scrapes per-core CPU metrics and returns one CPUMetric per core,
// plus an extra "total" entry averaged across all cores.
func collect(ctx context.Context, host string) ([]CPUMetric, error) {
	ts := time.Now()

	perCore, err := cpu.PercentWithContext(ctx, 0, true)
	if err != nil {
		return nil, fmt.Errorf("cpu.Percent: %w", err)
	}

	times, err := cpu.TimesWithContext(ctx, true)
	if err != nil {
		return nil, fmt.Errorf("cpu.Times: %w", err)
	}

	metrics := make([]CPUMetric, 0, len(perCore)+1)

	for i, pct := range perCore {
		m := CPUMetric{
			Time:         ts,
			Host:         host,
			Core:         fmt.Sprintf("cpu%d", i),
			UsagePercent: pct,
		}
		if i < len(times) {
			if total := times[i].Total(); total > 0 {
				m.UserPercent = times[i].User / total * 100
				m.SystemPercent = times[i].System / total * 100
				m.IdlePercent = times[i].Idle / total * 100
			}
		}
		metrics = append(metrics, m)
	}

	// Total across all cores
	if totalPct, err := cpu.PercentWithContext(ctx, 0, false); err == nil && len(totalPct) > 0 {
		metrics = append(metrics, CPUMetric{
			Time:         ts,
			Host:         host,
			Core:         "total",
			UsagePercent: totalPct[0],
		})
	}

	return metrics, nil
}

// flush writes a batch of CPUMetrics to TimescaleDB using the COPY protocol.
func flush(ctx context.Context, pool *pgxpool.Pool, metrics []CPUMetric) error {
	if len(metrics) == 0 {
		return nil
	}

	cols := []string{"time", "host", "core", "usage_percent", "user_percent", "system_percent", "idle_percent"}

	_, err := pool.CopyFrom(ctx,
		pgx.Identifier{"cpu_metrics"},
		cols,
		pgx.CopyFromSlice(len(metrics), func(i int) ([]any, error) {
			m := metrics[i]
			return []any{m.Time, m.Host, m.Core, m.UsagePercent, m.UserPercent, m.SystemPercent, m.IdlePercent}, nil
		}),
	)
	return err
}

func main() {
	ctx := context.Background()

	host, _ := os.Hostname()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "password"),
		getEnv("DB_NAME", "metricsdb"),
	)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer pool.Close()

	log.Printf("agent started — host=%s", host)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		metrics, err := collect(ctx, host)
		if err != nil {
			log.Printf("collect error: %v", err)
			continue
		}

		if err := flush(ctx, pool, metrics); err != nil {
			log.Printf("flush error: %v", err)
			continue
		}

		log.Printf("flushed %d cpu metrics", len(metrics))
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
