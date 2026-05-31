ALTER TABLE applications
    ADD COLUMN IF NOT EXISTS interview_at timestamptz NULL,
    ADD COLUMN IF NOT EXISTS stale_after timestamptz NULL DEFAULT (NOW() + interval '30 days')
