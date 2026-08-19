-- name: CreateDestination :one
INSERT INTO destinations (user_id, name, latitude, longitude, radius_m, target_text)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListDestinationsByUser :many
SELECT * FROM destinations WHERE user_id = $1 ORDER BY created_at DESC;

-- name: GetDestination :one
SELECT * FROM destinations WHERE id = $1 AND user_id = $2 LIMIT 1;

-- name: DeleteDestination :exec
DELETE FROM destinations WHERE id = $1 AND user_id = $2;
