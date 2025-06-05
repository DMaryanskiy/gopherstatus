CREATE TABLE IF NOT EXISTS check_results (
    id SERIAL PRIMARY KEY,
    service_name TEXT,
    url TEXT,
    method TEXT,
    online BOOLEAN,
    response_ms BIGINT,
    checked_at TIMESTAMP,
    error TEXT
);