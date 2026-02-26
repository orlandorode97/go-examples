package main

import (
	"context"
	"fmt"
	"time"
)

// InsertMetric inserts a single sensor reading into the hypertable.
func (d *DB) InsertMetric(ctx context.Context, m SensorMetric) error {
	query := `
		INSERT INTO sensor_metrics (time, sensor_id, location, temperature, humidity, pressure)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := d.Pool.Exec(ctx, query,
		m.Time, m.SensorID, m.Location,
		m.Temperature, m.Humidity, m.Pressure,
	)
	return err
}

// GetRecentMetrics returns raw readings for the last N hours.
func (d *DB) GetRecentMetrics(ctx context.Context, hours int) ([]SensorMetric, error) {
	query := `
		SELECT time, sensor_id, location, temperature, humidity, pressure
		FROM sensor_metrics
		WHERE time > NOW() - make_interval(hours => $1)
		ORDER BY time DESC
		LIMIT 500
	`
	rows, err := d.Pool.Query(ctx, query, hours)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []SensorMetric
	for rows.Next() {
		var m SensorMetric
		if err := rows.Scan(&m.Time, &m.SensorID, &m.Location,
			&m.Temperature, &m.Humidity, &m.Pressure); err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}
	return metrics, rows.Err()
}

// GetBucketedMetrics uses TimescaleDB's time_bucket() to aggregate readings.
func (d *DB) GetBucketedMetrics(ctx context.Context, bucketMinutes int, since time.Time) ([]BucketedMetric, error) {
	query := `
		SELECT
			time_bucket($1::INTERVAL, time)  AS bucket,
			sensor_id,
			AVG(temperature)                 AS avg_temperature,
			AVG(humidity)                    AS avg_humidity,
			AVG(pressure)                    AS avg_pressure,
			COUNT(*)                         AS reading_count
		FROM sensor_metrics
		WHERE time >= $2
		GROUP BY bucket, sensor_id
		ORDER BY bucket DESC, sensor_id
	`
	interval := fmt.Sprintf("%d minutes", bucketMinutes)
	rows, err := d.Pool.Query(ctx, query, interval, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []BucketedMetric
	for rows.Next() {
		var b BucketedMetric
		if err := rows.Scan(&b.Bucket, &b.SensorID,
			&b.AvgTemp, &b.AvgHumidity, &b.AvgPressure, &b.ReadingCount); err != nil {
			return nil, err
		}
		results = append(results, b)
	}
	return results, rows.Err()
}
