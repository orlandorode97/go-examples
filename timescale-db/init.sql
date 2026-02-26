-- Enable TimescaleDB extension
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- Create the sensor_metrics table
CREATE TABLE IF NOT EXISTS sensor_metrics (
    time        TIMESTAMPTZ       NOT NULL,
    sensor_id   TEXT              NOT NULL,
    location    TEXT              NOT NULL,
    temperature DOUBLE PRECISION,
    humidity    DOUBLE PRECISION,
    pressure    DOUBLE PRECISION
);

-- Convert to a hypertable (TimescaleDB magic ✨)
SELECT create_hypertable('sensor_metrics', 'time', if_not_exists => TRUE);

-- Create an index on sensor_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_sensor_metrics_sensor_id
    ON sensor_metrics (sensor_id, time DESC);

-- Seed some sample data
INSERT INTO sensor_metrics (time, sensor_id, location, temperature, humidity, pressure)
VALUES
    (NOW() - INTERVAL '2 hours', 'sensor-001', 'warehouse-a', 22.5, 60.1, 1013.2),
    (NOW() - INTERVAL '90 minutes', 'sensor-001', 'warehouse-a', 23.1, 58.4, 1012.8),
    (NOW() - INTERVAL '1 hour', 'sensor-001', 'warehouse-a', 24.0, 57.9, 1012.5),
    (NOW() - INTERVAL '30 minutes', 'sensor-001', 'warehouse-a', 23.7, 59.2, 1013.0),
    (NOW() - INTERVAL '2 hours', 'sensor-002', 'warehouse-b', 18.2, 72.3, 1014.1),
    (NOW() - INTERVAL '1 hour', 'sensor-002', 'warehouse-b', 19.0, 70.5, 1013.9),
    (NOW() - INTERVAL '30 minutes', 'sensor-002', 'warehouse-b', 18.8, 71.1, 1014.0);
