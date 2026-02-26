docker exec -it timescaledb psql -U postgres -d sensordb

-- Example queries
SELECT * FROM sensor_metrics ORDER BY time DESC LIMIT 10;

SELECT time_bucket('30 minutes', time) AS bucket,
       sensor_id, AVG(temperature), COUNT(*)
FROM sensor_metrics
GROUP BY bucket, sensor_id
ORDER BY bucket DESC;
