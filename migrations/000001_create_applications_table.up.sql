CREATE TABLE IF NOT EXISTS applications (
    id bigserial PRIMARY KEY,
    company_name text NOT NULL,
    orle_title text NOT NULL,
    status text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    last_communication text NOT NULL,
    notes text NOT NULL,
    version int NOT NULL DEFAULT 1
);

