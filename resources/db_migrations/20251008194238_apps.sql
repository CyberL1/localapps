-- +goose Up
CREATE TABLE apps (
    id INTEGER PRIMARY KEY,
    app_id TEXT NOT NULL,
    name TEXT NOT NULL,
    icon TEXT NOT NULL DEFAULT '',
    installed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    parts TEXT NOT NULL CHECK (json_valid(parts))
);

-- +goose Down
SELECT 'down SQL query';
DROP TABLE apps;
