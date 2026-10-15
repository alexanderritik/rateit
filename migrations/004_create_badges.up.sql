CREATE TABLE IF NOT EXISTS badges (
    id         SERIAL PRIMARY KEY,
    device_id  VARCHAR(64) NOT NULL,
    badge_id   VARCHAR(32) NOT NULL,
    earned_at  TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(device_id, badge_id)
);

CREATE INDEX IF NOT EXISTS idx_badges_device_id ON badges(device_id);
