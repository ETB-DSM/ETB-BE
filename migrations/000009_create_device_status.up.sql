CREATE TABLE device_status (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id  TEXT NOT NULL,
    battery    INT NOT NULL,
    lidar_ok   BOOL NOT NULL DEFAULT TRUE,
    camera_ok  BOOL NOT NULL DEFAULT TRUE,
    gps_ok     BOOL NOT NULL DEFAULT TRUE,
    network_ok BOOL NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_device_status_device_id_created_at ON device_status(device_id, created_at DESC);
