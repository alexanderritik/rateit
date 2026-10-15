CREATE TABLE IF NOT EXISTS streaks (
    device_id       VARCHAR(64) PRIMARY KEY,
    current_streak  INT DEFAULT 0,
    longest_streak  INT DEFAULT 0,
    last_played_date DATE,
    total_points    INT DEFAULT 0,
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);
