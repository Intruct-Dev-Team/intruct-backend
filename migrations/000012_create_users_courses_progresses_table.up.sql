CREATE TABLE IF NOT EXISTS public.users_courses_progresses(
       user_id INTEGER NOT NULL REFERENCES users(user_id),
       course_id INTEGER NOT NULL REFERENCES courses(course_id),
       current_lesson_id INTEGER NOT NULL REFERENCES lessons(lesson_id),
       is_finished BOOLEAN NOT NULL DEFAULT FALSE,
       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
       PRIMARY KEY (user_id, course_id)
);