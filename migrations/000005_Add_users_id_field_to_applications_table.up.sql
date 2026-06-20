ALTER TABLE applications
    ADD COLUMN users_id bigint NOT NULL REFERENCES users (id) ON DELETE CASCADE;

CREATE INDEX idx_applications_users_id ON applications (users_id);

