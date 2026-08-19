-- name: CreateUser :one
INSERT INTO users (email, nickname, password_hash, email_verified)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: UpdateEmailVerified :exec
UPDATE users SET email_verified = $1 WHERE email = $2;
