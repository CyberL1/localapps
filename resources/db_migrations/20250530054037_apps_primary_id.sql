-- +goose Up
ALTER TABLE apps ADD CONSTRAINT apps_pkey PRIMARY KEY (id);

-- +goose Down
ALTER TABLE apps DROP CONSTRAINT apps_pkey;
