CREATE TABLE IF NOT EXISTS public.ratings (
    rating_id  SERIAL PRIMARY KEY,
    course_id  INTEGER NOT NULL REFERENCES courses(course_id),
    user_id    INTEGER NOT NULL REFERENCES users(user_id),
    rating     INTEGER NOT NULL CHECK (rating BETWEEN 1 AND 5),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (course_id, user_id)
);
