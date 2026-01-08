CREATE TABLE IF NOT EXISTS public.states (
    state_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO public.states (state_id, name)
VALUES
    (1, 'in creation'),
    (2, 'failed'),
    (3, 'created'),
    (4, 'published')

ON CONFLICT DO NOTHING;