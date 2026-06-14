CREATE TABLE IF NOT EXISTS tokens (
    hash bytea NOT NULL,
    users_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    expiry timestamptz NOT NULL,
    scope text NOT NULL
);

