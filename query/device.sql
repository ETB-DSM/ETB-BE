-- name: CreateDevice :one
INSERT INTO devices (user_id, name)
VALUES ($1, $2)
RETURNING *;

-- name: ListDevicesByUser :many
SELECT * FROM devices WHERE user_id = $1 ORDER BY created_at ASC;

-- name: GetDevice :one
SELECT * FROM devices WHERE id = $1 AND user_id = $2 LIMIT 1;

-- name: DeleteDevice :exec
DELETE FROM devices WHERE id = $1 AND user_id = $2;

-- name: CountDevicesByUser :one
SELECT COUNT(*) FROM devices WHERE user_id = $1;
