-- +goose Up
CREATE TABLE apps ("id" text NOT NULL, "installedAt" timestamp NOT NULL);

-- +goose Down
DROP DATABASE apps;
