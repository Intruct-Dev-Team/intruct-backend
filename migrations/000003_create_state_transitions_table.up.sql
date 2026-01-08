CREATE TABLE IF NOT EXISTS public.state_transitions (
    transition_id SERIAL PRIMARY KEY,
    state_machine_id INTEGER NOT NULL REFERENCES state_machines(state_machine_id),
    current_state_id INTEGER NOT NULL REFERENCES states(state_id),
    next_state_id INTEGER NOT NULL REFERENCES states(state_id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (state_machine_id, current_state_id, next_state_id)
);

INSERT INTO public.state_transitions (transition_id, state_machine_id, current_state_id, next_state_id)
VALUES
    -- course

    (1, 1, 1, 2),
    (2, 1, 1, 3),
    (3, 1, 3, 4)

ON CONFLICT DO NOTHING;