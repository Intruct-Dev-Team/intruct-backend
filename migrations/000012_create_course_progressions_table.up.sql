CREATE TABLE IF NOT EXISTS public.course_progressions(
       user_id INTEGER NOT NULL REFERENCES users(user_id),
       course_id INTEGER NOT NULL REFERENCES courses(course_id),
       current_lesson_id INTEGER REFERENCES lessons(lesson_id) ON DELETE SET NULL;,
       finished_lessons_count INTEGER NOT NULL DEFAULT 0,
       is_finished BOOLEAN NOT NULL DEFAULT FALSE,
       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
       PRIMARY KEY (user_id, course_id)
);