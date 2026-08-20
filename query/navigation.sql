-- name: CreateNavigationSession :one
INSERT INTO navigation_sessions (user_id, device_id, destination_id, start_latitude, start_longitude)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetNavigationSession :one
SELECT * FROM navigation_sessions WHERE id = $1 LIMIT 1;

-- name: UpdateNavigationSessionStatus :one
UPDATE navigation_sessions
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: CreateNavigationInstruction :one
INSERT INTO navigation_instructions (session_id, action, distance_meters, message)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetLatestNavigationInstruction :one
SELECT * FROM navigation_instructions
WHERE session_id = $1
ORDER BY created_at DESC
LIMIT 1;
