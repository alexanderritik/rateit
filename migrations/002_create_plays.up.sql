CREATE TABLE IF NOT EXISTS plays (
    id           SERIAL PRIMARY KEY,
    device_id    VARCHAR(64) NOT NULL,
    challenge_id INT NOT NULL REFERENCES challenges(id),
    correct      BOOLEAN NOT NULL,
    time_ms      INT NOT NULL,
    points       INT NOT NULL DEFAULT 0,
    played_at    TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(device_id, challenge_id)
);

CREATE INDEX IF NOT EXISTS idx_plays_challenge_id ON plays(challenge_id);
CREATE INDEX IF NOT EXISTS idx_plays_device_id ON plays(device_id);
