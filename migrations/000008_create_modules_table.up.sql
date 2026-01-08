CREATE TABLE IF NOT EXISTS public.modules(
       module_id SERIAL PRIMARY KEY NOT NULL,
       course_id INTEGER NOT NULL REFERENCES courses(course_id),
       serial_number INTEGER NOT NULL,
       title TEXT NOT NULL,
       description TEXT NOT NULL DEFAULT '',
       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
       CONSTRAINT unique_course_module UNIQUE (course_id, serial_number)
);