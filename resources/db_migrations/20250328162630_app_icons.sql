-- +goose Up
ALTER TABLE apps ADD COLUMN icon text NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE apps DROP COLUMN icon;
