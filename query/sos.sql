-- name: CreateSosEvent :one
INSERT INTO sos_events (user_id, device_id, event_type, latitude, longitude, battery)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListSosEventsByUser :many
SELECT * FROM sos_events WHERE user_id = $1 ORDER BY created_at DESC;
