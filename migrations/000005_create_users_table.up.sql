CREATE TABLE IF NOT EXISTS public.users(
       user_id SERIAL PRIMARY KEY NOT NULL,
       external_uuid UUID UNIQUE NOT NULL,
       email TEXT UNIQUE NOT NULL,
       name TEXT NOT NULL,
       surname TEXT NOT NULL,
       registration_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
       birthdate DATE NOT NULL,
       avatar TEXT NOT NULL DEFAULT '',
       current_day_streak INTEGER NOT NULL DEFAULT 0,
       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);