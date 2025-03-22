-- name: GetConfig :many
SELECT * FROM config;

-- name: SetConfigKey :one
INSERT INTO config (key, value) VALUES ($1, $2)
RETURNING *;

-- name: UpdateConfigKey :one
UPDATE config SET value = $2 WHERE key = $1
RETURNING *;
