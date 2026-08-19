-- name: CreateSosEvent :one
INSERT INTO sos_events (user_id, device_id, event_type, latitude, longitude)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListSosEventsByUser :many
SELECT * FROM sos_events WHERE user_id = $1 ORDER BY created_at DESC;
