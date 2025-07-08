-- name: ListApps :many
SELECT * FROM apps ORDER BY installed_at;

-- name: GetAppById :one
SELECT * FROM apps WHERE id = $1;

-- name: GetAppByAppId :one
SELECT * FROM apps WHERE app_id = $1;

-- name: CreateApp :one
INSERT INTO apps (app_id, installed_at, name, parts, icon) VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateApp :one
UPDATE apps SET name = $2, parts = $3, icon = $4 WHERE app_id = $1
RETURNING *;

-- name: DeleteApp :exec
DELETE FROM apps WHERE id = $1;
