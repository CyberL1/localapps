-- name: GetConfig :many
SELECT * FROM config;

-- name: SetConfigKey :one
INSERT INTO config (key, value) VALUES (?, ?)
RETURNING *;

-- name: UpdateConfigKey :one
UPDATE config SET value = ? WHERE key = ?
RETURNING *;
