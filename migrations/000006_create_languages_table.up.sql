CREATE TABLE IF NOT EXISTS public.languages(
       language_id SERIAL PRIMARY KEY NOT NULL,
       name TEXT NOT NULL,
       eng_name TEXT NOT NULL,
       iso_code TEXT NOT NULL,
       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO public.languages (id, name, eng_name, iso_code)
VALUES
    (1, 'English', 'English', 'en'),
    (2, 'Srpski', 'Serbian', 'sr'),
    (3, '中文', 'Chinese', 'zh'),
    (4, 'हिन्दी', 'Hindi', 'hi'),
    (5, 'Русский', 'Russian', 'ru'),
    (6, 'Deutsch', 'German', 'de'),
    (7, 'Español', 'Spanish', 'es'),
    (8, 'Português', 'Portuguese', 'pt'),
    (9, 'Français', 'French', 'fr')
ON CONFLICT DO NOTHING;
