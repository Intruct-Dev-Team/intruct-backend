CREATE TABLE IF NOT EXISTS public.quizes_options(
       option_id SERIAL PRIMARY KEY NOT NULL,
       quiz_id INTEGER NOT NULL REFERENCES quizes(quiz_id),
       content TEXT NOT NULL,
       is_answer BOOLEAN NOT NULL DEFAULT FALSE,
       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);