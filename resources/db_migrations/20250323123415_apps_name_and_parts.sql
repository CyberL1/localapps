-- +goose Up
ALTER TABLE apps ADD COLUMN name text NOT NULL;
ALTER TABLE apps ADD COLUMN parts jsonb;

-- +goose Down
ALTER TABLE apps DROP COLUMN name;
ALTER TABLE apps DROP COLUMN parts;
