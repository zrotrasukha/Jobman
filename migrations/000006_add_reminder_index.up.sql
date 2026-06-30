-- composite index is being used for fast retreival of upcoming interviews
CREATE INDEX IF NOT EXISTS idx_applications_reminder ON applications (users_id, interview_at)
WHERE
    status = 'Interviewing'
    AND interview_at IS NOT NULL;