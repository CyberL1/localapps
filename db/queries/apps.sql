-- name: ListApps :many
SELECT * FROM apps;

-- name: GetApp :one
SELECT * FROM apps
WHERE id = $1;

-- name: CreateApp :one
INSERT INTO apps (id) VALUES ($1)
RETURNING *;

-- name: DeleteApp :exec
DELETE FROM apps WHERE id = $1;
