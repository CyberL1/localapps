-- name: GetConfig :many
SELECT * FROM config;

-- name: SetConfigKey :one
INSERT INTO config (key, value) VALUES ($1, $2)
RETURNING *;
