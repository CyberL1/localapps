-- +goose Up
ALTER TABLE apps ADD COLUMN app_id text NOT NULL;

-- +goose Down
ALTER TABLE apps DROP COLUMN app_id;
