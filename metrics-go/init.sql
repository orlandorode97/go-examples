CREATE EXTENSION IF NOT EXISTS timescaledb;

CREATE TABLE IF NOT EXISTS cpu_metrics (
    time            TIMESTAMPTZ     NOT NULL,
    host            TEXT            NOT NULL,
    core            TEXT            NOT NULL,
    usage_percent   DOUBLE PRECISION NOT NULL,
    user_percent    DOUBLE PRECISION NOT NULL,
    system_percent  DOUBLE PRECISION NOT NULL,
    idle_percent    DOUBLE PRECISION NOT NULL
);

SELECT create_hypertable('cpu_metrics', 'time', if_not_exists => TRUE);

CREATE INDEX IF NOT EXISTS idx_cpu_metrics_host_time ON cpu_metrics (host, time DESC);

CREATE TABLE IF NOT EXISTS mem_metrics (
    time            TIMESTAMPTZ     NOT NULL,
    host            TEXT            NOT NULL,
    total_bytes     BIGINT          NOT NULL,
    available_bytes BIGINT          NOT NULL,
    used_bytes      BIGINT          NOT NULL,
    used_percent    DOUBLE PRECISION NOT NULL
);

SELECT create_hypertable('mem_metrics', 'time', if_not_exists => TRUE);

CREATE INDEX IF NOT EXISTS idx_mem_metrics_host_time ON mem_metrics (host, time DESC);
