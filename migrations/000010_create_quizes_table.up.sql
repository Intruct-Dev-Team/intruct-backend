CREATE TABLE IF NOT EXISTS public.quizes(
       quiz_id SERIAL PRIMARY KEY NOT NULL,
       lesson_id INTEGER NOT NULL REFERENCES lessons(lesson_id),
       serial_number INTEGER NOT NULL,
       question TEXT NOT NULL,
       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
       CONSTRAINT unique_lesson_quiz UNIQUE (lesson_id, serial_number)
);