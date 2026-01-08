CREATE TABLE IF NOT EXISTS public.lessons(
       lesson_id SERIAL PRIMARY KEY NOT NULL,
       course_id INTEGER NOT NULL REFERENCES courses(course_id),
       module_id INTEGER NOT NULL REFERENCES modules(module_id),
       serial_number INTEGER NOT NULL,
       title TEXT NOT NULL,
       description TEXT NOT NULL DEFAULT '',
       content TEXT NOT NULL,
       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
       CONSTRAINT unique_course_lesson UNIQUE (course_id, serial_number)
);