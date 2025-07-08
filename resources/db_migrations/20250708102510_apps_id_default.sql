-- +goose Up
ALTER TABLE apps ALTER COLUMN id SET DEFAULT gen_random_uuid();

-- +goose Down
ALTER TABLE apps ALTER COLUMN id DROP DEFAULT;
