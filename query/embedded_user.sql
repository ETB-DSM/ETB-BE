-- name: CreateEmbeddedUser :one
INSERT INTO embedded_users (user_id, name, guardian_name, guardian_phone, device_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetEmbeddedUserByID :one
SELECT * FROM embedded_users WHERE user_id = $1 LIMIT 1;

-- name: GetEmbeddedUserByDeviceID :one
SELECT * FROM embedded_users WHERE device_id = $1 LIMIT 1;
