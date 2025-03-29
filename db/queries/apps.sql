-- name: ListApps :many
SELECT * FROM apps ORDER BY installed_at;

-- name: GetApp :one
SELECT * FROM apps WHERE id = $1;

-- name: CreateApp :one
INSERT INTO apps (id, installed_at, name, parts, icon) VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteApp :exec
DELETE FROM apps WHERE id = $1;
