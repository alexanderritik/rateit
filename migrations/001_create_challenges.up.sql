CREATE TABLE IF NOT EXISTS challenges (
    id          SERIAL PRIMARY KEY,
    play_date   DATE UNIQUE NOT NULL,
    segments    JSONB NOT NULL,
    mistake_id  INT NOT NULL,
    explanation TEXT NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);
