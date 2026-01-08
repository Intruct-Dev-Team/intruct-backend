CREATE TABLE IF NOT EXISTS public.languages(
       language_id SERIAL PRIMARY KEY NOT NULL,
       name TEXT NOT NULL,
       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO public.languages (language_id, name)
VALUES
    (1, 'English'),
    (2, 'Russian'),
    (3, 'Serbian')

ON CONFLICT DO NOTHING;