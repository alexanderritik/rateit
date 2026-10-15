CREATE TABLE IF NOT EXISTS clips (
    id          SERIAL PRIMARY KEY,
    play_date   DATE NOT NULL,
    clip_index  SMALLINT NOT NULL,   -- 1..5, display order
    text        TEXT NOT NULL,
    language    VARCHAR(10) NOT NULL DEFAULT 'en-US',
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(play_date, clip_index)
);

CREATE TABLE IF NOT EXISTS ratings (
    id          SERIAL PRIMARY KEY,
    device_id   VARCHAR(64) NOT NULL,
    clip_id     INT NOT NULL REFERENCES clips(id),
    score       SMALLINT NOT NULL CHECK (score BETWEEN 1 AND 5),
    listened    BOOLEAN NOT NULL DEFAULT false,
    rated_at    TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(device_id, clip_id)
);

CREATE INDEX IF NOT EXISTS idx_ratings_device ON ratings(device_id);
CREATE INDEX IF NOT EXISTS idx_ratings_clip   ON ratings(clip_id);
CREATE INDEX IF NOT EXISTS idx_clips_date     ON clips(play_date);
