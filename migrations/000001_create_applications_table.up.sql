CREATE TABLE IF NOT EXISTS applications (
    id bigserial PRIMARY KEY,
    company_name text NOT NULL,
    role_title text NOT NULL,
    status text NOT NULL,
    applied_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    last_communication timestamptz,
    notes text,
    version integer NOT NULL DEFAULT 1
);

