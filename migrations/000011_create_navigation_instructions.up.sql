CREATE TABLE navigation_instructions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id      UUID NOT NULL REFERENCES navigation_sessions(id) ON DELETE CASCADE,
    action          TEXT NOT NULL,
    distance_meters INT,
    message         TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_nav_instructions_session_id ON navigation_instructions(session_id, created_at DESC);
