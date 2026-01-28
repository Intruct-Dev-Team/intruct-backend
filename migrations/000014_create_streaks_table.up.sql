CREATE TABLE IF NOT EXISTS public.streaks (
    user_id           INTEGER NOT NULL REFERENCES users(user_id) PRIMARY KEY,
    days_streak_count INTEGER NOT NULL DEFAULT 0,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
