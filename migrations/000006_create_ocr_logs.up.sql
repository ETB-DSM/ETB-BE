CREATE TABLE ocr_logs (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    destination_id  UUID NOT NULL REFERENCES destinations(id) ON DELETE CASCADE,
    recognized_text TEXT NOT NULL,
    target_text     TEXT NOT NULL,
    matched         BOOLEAN NOT NULL,
    confidence      REAL NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
