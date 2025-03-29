-- +goose Up
CREATE TABLE apps ("id" text NOT NULL, "installed_at" timestamp NOT NULL);

-- +goose Down
DROP TABLE apps;
