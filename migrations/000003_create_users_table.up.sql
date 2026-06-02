CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    created_at timestamptz NOT NULL DEFAULT now(),
    name text NOT NULL,
    email citext NOT NULL UNIQUE,
    password_hash bytea NOT NULL,
    activated bool NOT NULL DEFAULT FALSE,
    version integer NOT NULL DEFAULT 1
);

