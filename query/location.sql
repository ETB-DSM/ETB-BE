-- name: CreateLocation :one
INSERT INTO locations (user_id, device_id, latitude, longitude, recorded_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
