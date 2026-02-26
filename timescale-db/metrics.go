package main

import "time"

// SensorMetric represents a single sensor reading.
type SensorMetric struct {
	Time        time.Time `json:"time"`
	SensorID    string    `json:"sensor_id"`
	Location    string    `json:"location"`
	Temperature *float64  `json:"temperature,omitempty"`
	Humidity    *float64  `json:"humidity,omitempty"`
	Pressure    *float64  `json:"pressure,omitempty"`
}

// BucketedMetric represents an aggregated (time-bucketed) result from TimescaleDB.
type BucketedMetric struct {
	Bucket       time.Time `json:"bucket"`
	SensorID     string    `json:"sensor_id"`
	AvgTemp      *float64  `json:"avg_temperature,omitempty"`
	AvgHumidity  *float64  `json:"avg_humidity,omitempty"`
	AvgPressure  *float64  `json:"avg_pressure,omitempty"`
	ReadingCount int64     `json:"reading_count"`
}
