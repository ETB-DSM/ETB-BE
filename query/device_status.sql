-- name: UpsertDeviceStatus :one
INSERT INTO device_status (device_id, battery, lidar_ok, camera_ok, gps_ok, network_ok)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetLatestDeviceStatus :one
SELECT * FROM device_status
WHERE device_id = $1
ORDER BY created_at DESC
LIMIT 1;
