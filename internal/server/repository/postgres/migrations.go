package postgres

const (
	createCounterMetricsTable = `CREATE TABLE IF NOT EXISTS counter_metrics
(
    id          INT                 GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    name        VARCHAR(64)         UNIQUE NOT NULL,
    delta       BIGINT              NOT NULL,
    value       BIGINT              NOT NULL,
    updated_at  TIMESTAMPTZ(6)      DEFAULT NOW()
);`
	createGaugeMetricsTable = `CREATE TABLE IF NOT EXISTS gauge_metrics
(
    id          INT                 GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    name        VARCHAR(64)         UNIQUE NOT NULL,
    value       DOUBLE PRECISION    NOT NULL,
    updated_at  TIMESTAMPTZ(6)      DEFAULT NOW()
);`
)
