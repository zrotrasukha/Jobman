-- Seed data for new statuses: Selected, RoundCleared, Declined.
-- Deletes existing applications with these names to prevent duplicates if run multiple times.

-- ============================================================================
-- USER 2 (two@gmail.com)
-- ============================================================================
DELETE FROM status_history WHERE application_id IN (SELECT id FROM applications WHERE users_id = 2 AND company_name IN ('Seed Company Selected', 'Seed Company RoundCleared', 'Seed Company Declined'));
DELETE FROM applications WHERE users_id = 2 AND company_name IN ('Seed Company Selected', 'Seed Company RoundCleared', 'Seed Company Declined');

-- 1. Selected transition for User 2
WITH ins_selected AS (
  INSERT INTO applications (company_name, role_title, status, applied_at, last_communication, users_id, updated_at)
  VALUES ('Seed Company Selected', 'Senior Staff Engineer', 'Selected', NOW() - INTERVAL '2 hours', NOW() - INTERVAL '2 hours', 2, NOW())
  RETURNING id
)
INSERT INTO status_history (application_id, status, changed_at)
SELECT id, 'Applied', NOW() - INTERVAL '4 hours' FROM ins_selected
UNION ALL
SELECT id, 'Interviewing', NOW() - INTERVAL '3 hours' FROM ins_selected
UNION ALL
SELECT id, 'Offered', NOW() - INTERVAL '2 hours' FROM ins_selected
UNION ALL
SELECT id, 'Selected', NOW() - INTERVAL '1 hour' FROM ins_selected;

-- 2. RoundCleared transition for User 2
WITH ins_cleared AS (
  INSERT INTO applications (company_name, role_title, status, applied_at, last_communication, users_id, updated_at)
  VALUES ('Seed Company RoundCleared', 'Systems Programmer', 'RoundCleared', NOW() - INTERVAL '2 hours', NOW() - INTERVAL '2 hours', 2, NOW())
  RETURNING id
)
INSERT INTO status_history (application_id, status, changed_at)
SELECT id, 'Applied', NOW() - INTERVAL '3 hours' FROM ins_cleared
UNION ALL
SELECT id, 'Interviewing', NOW() - INTERVAL '2 hours' FROM ins_cleared
UNION ALL
SELECT id, 'RoundCleared', NOW() - INTERVAL '1 hour' FROM ins_cleared;

-- 3. Declined transition for User 2
WITH ins_declined AS (
  INSERT INTO applications (company_name, role_title, status, applied_at, last_communication, users_id, updated_at)
  VALUES ('Seed Company Declined', 'Cloud Architect', 'Declined', NOW() - INTERVAL '2 hours', NOW() - INTERVAL '2 hours', 2, NOW())
  RETURNING id
)
INSERT INTO status_history (application_id, status, changed_at)
SELECT id, 'Applied', NOW() - INTERVAL '4 hours' FROM ins_declined
UNION ALL
SELECT id, 'Interviewing', NOW() - INTERVAL '3 hours' FROM ins_declined
UNION ALL
SELECT id, 'Offered', NOW() - INTERVAL '2 hours' FROM ins_declined
UNION ALL
SELECT id, 'Declined', NOW() - INTERVAL '1 hour' FROM ins_declined;


-- ============================================================================
-- USER 1 (Shivang - one@gmail.com)
-- ============================================================================
DELETE FROM status_history WHERE application_id IN (SELECT id FROM applications WHERE users_id = 1 AND company_name IN ('Seed Company Selected', 'Seed Company RoundCleared', 'Seed Company Declined'));
DELETE FROM applications WHERE users_id = 1 AND company_name IN ('Seed Company Selected', 'Seed Company RoundCleared', 'Seed Company Declined');

-- 1. Selected transition for User 1
WITH ins_selected AS (
  INSERT INTO applications (company_name, role_title, status, applied_at, last_communication, users_id, updated_at)
  VALUES ('Seed Company Selected', 'Senior Staff Engineer', 'Selected', NOW() - INTERVAL '2 hours', NOW() - INTERVAL '2 hours', 1, NOW())
  RETURNING id
)
INSERT INTO status_history (application_id, status, changed_at)
SELECT id, 'Applied', NOW() - INTERVAL '4 hours' FROM ins_selected
UNION ALL
SELECT id, 'Interviewing', NOW() - INTERVAL '3 hours' FROM ins_selected
UNION ALL
SELECT id, 'Offered', NOW() - INTERVAL '2 hours' FROM ins_selected
UNION ALL
SELECT id, 'Selected', NOW() - INTERVAL '1 hour' FROM ins_selected;

-- 2. RoundCleared transition for User 1
WITH ins_cleared AS (
  INSERT INTO applications (company_name, role_title, status, applied_at, last_communication, users_id, updated_at)
  VALUES ('Seed Company RoundCleared', 'Systems Programmer', 'RoundCleared', NOW() - INTERVAL '2 hours', NOW() - INTERVAL '2 hours', 1, NOW())
  RETURNING id
)
INSERT INTO status_history (application_id, status, changed_at)
SELECT id, 'Applied', NOW() - INTERVAL '3 hours' FROM ins_cleared
UNION ALL
SELECT id, 'Interviewing', NOW() - INTERVAL '2 hours' FROM ins_cleared
UNION ALL
SELECT id, 'RoundCleared', NOW() - INTERVAL '1 hour' FROM ins_cleared;

-- 3. Declined transition for User 1
WITH ins_declined AS (
  INSERT INTO applications (company_name, role_title, status, applied_at, last_communication, users_id, updated_at)
  VALUES ('Seed Company Declined', 'Cloud Architect', 'Declined', NOW() - INTERVAL '2 hours', NOW() - INTERVAL '2 hours', 1, NOW())
  RETURNING id
)
INSERT INTO status_history (application_id, status, changed_at)
SELECT id, 'Applied', NOW() - INTERVAL '4 hours' FROM ins_declined
UNION ALL
SELECT id, 'Interviewing', NOW() - INTERVAL '3 hours' FROM ins_declined
UNION ALL
SELECT id, 'Offered', NOW() - INTERVAL '2 hours' FROM ins_declined
UNION ALL
SELECT id, 'Declined', NOW() - INTERVAL '1 hour' FROM ins_declined;
