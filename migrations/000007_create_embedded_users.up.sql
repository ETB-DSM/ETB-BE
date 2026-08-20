CREATE TABLE embedded_users (
    user_id        TEXT PRIMARY KEY,
    name           TEXT NOT NULL,
    guardian_name  TEXT NOT NULL,
    guardian_phone TEXT NOT NULL,
    device_id      TEXT NOT NULL UNIQUE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
