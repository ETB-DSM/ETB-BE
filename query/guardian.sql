-- name: CreateGuardian :one
INSERT INTO guardians (user_id, name, phone)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListGuardiansByUser :many
SELECT * FROM guardians WHERE user_id = $1 ORDER BY created_at ASC;

-- name: GetGuardian :one
SELECT * FROM guardians WHERE id = $1 AND user_id = $2 LIMIT 1;

-- name: DeleteGuardian :exec
DELETE FROM guardians WHERE id = $1 AND user_id = $2;

-- name: CountGuardiansByUser :one
SELECT COUNT(*) FROM guardians WHERE user_id = $1;
