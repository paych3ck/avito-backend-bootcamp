CREATE TABLE IF NOT EXISTS houses (
    id SERIAL PRIMARY KEY,
    address TEXT NOT NULL,
    year INT NOT NULL,
    developer TEXT,
    created_at TIMESTAMP NOT NULL,
    update_at TIMESTAMP NOT NULL
);