-- +goose Up
CREATE TABLE config (
    key TEXT PRIMARY KEY,
    value TEXT
);

-- +goose Down
DROP TABLE config;
