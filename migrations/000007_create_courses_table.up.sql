CREATE TABLE IF NOT EXISTS public.courses(
       course_id SERIAL PRIMARY KEY NOT NULL,
       owner_id INTEGER NOT NULL REFERENCES users(user_id),
       title TEXT UNIQUE NOT NULL,
       description TEXT NOT NULL DEFAULT '',
       language_id INTEGER NOT NULL REFERENCES languages(language_id),
       state_machine_item_id INTEGER NOT NULL REFERENCES state_machines_items(item_id),
       is_public BOOLEAN NOT NULL DEFAULT FALSE,
       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);