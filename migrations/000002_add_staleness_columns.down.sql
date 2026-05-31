ALTER TABLE applications
    DROP COLUMN IF EXISTS interview_at,
    DROP COLUMN IF EXISTS stale_after;

