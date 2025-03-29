-- +goose Up
ALTER TABLE apps ADD COLUMN icon text NOT NULL;

-- +goose Down
ALTER TABLE apps DROP COLUMN icon;
