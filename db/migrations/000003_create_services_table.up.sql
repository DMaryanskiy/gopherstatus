CREATE TABLE IF NOT EXISTS services (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    method TEXT NOT NULL,
    interval INTEGER NOT NULL,
    body TEXT,
    user_id INTEGER REFERENCES users (id) 
);

CREATE TABLE IF NOT EXISTS headers (
    id SERIAL PRIMARY KEY,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    service_id INTEGER REFERENCES services (id)
);
