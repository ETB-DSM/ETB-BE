-- name: CreateOcrLog :one
INSERT INTO ocr_logs (destination_id, recognized_text, target_text, matched, confidence)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
