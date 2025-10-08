-- name: ListApps :many
SELECT * FROM apps ORDER BY installed_at;

-- name: GetAppById :one
SELECT * FROM apps WHERE id = ?;

-- name: GetAppByAppId :one
SELECT * FROM apps WHERE app_id = ?;

-- name: CreateApp :one
INSERT INTO apps (app_id, installed_at, name, parts, icon) VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateApp :one
UPDATE apps SET name = ?, parts = ?, icon = ? WHERE app_id = ?
RETURNING *;

-- name: DeleteApp :exec
DELETE FROM apps WHERE id = ?;
