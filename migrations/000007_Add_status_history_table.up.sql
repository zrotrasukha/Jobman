CREATE TABLE IF NOT EXISTS status_history (
    id bigserial PRIMARY KEY,
    application_id bigint NOT NULL REFERENCES applications (id) ON DELETE CASCADE,
    status text NOT NULL,
    changed_at timestamptz NOT NULL DEFAULT NOW()
);

