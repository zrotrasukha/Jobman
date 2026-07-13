-- Seed data to test 7d, 1m, and 1y windows for user 1 (Shivang)
-- Deletes existing applications for user 1 to ensure a clean state for testing
DELETE FROM applications WHERE users_id = 1;

-- ============================================================================
-- GROUP A: 7d window (Current Week, current month, current year)
-- ============================================================================

-- A1: Applied 1 hour ago, Replied 1 hour ago (counts for 7d applied, 7d replied)
WITH ins AS (
  INSERT INTO applications (company_name, role_title, status, applied_at, last_communication, users_id)
  VALUES ('7d Company A', 'Software Engineer', 'Applied', NOW() - INTERVAL '1 hour', NOW() - INTERVAL '1 hour', 1)
  RETURNING id
)
INSERT INTO status_history (application_id, status, changed_at)
SELECT id, 'Applied', NOW() - INTERVAL '1 hour' FROM ins;

-- A2: Applied 1 day ago, Interviewing 12 hours ago (counts for 7d applied, 7d interviewing)
WITH ins AS (
  INSERT INTO applications (company_name, role_title, status, applied_at, users_id)
  VALUES ('7d Company B', 'Frontend Developer', 'Interviewing', NOW() - INTERVAL '1 day', 1)
  RETURNING id
)
INSERT INTO status_history (application_id, status, changed_at)
SELECT id, 'Applied', NOW() - INTERVAL '1 day' FROM ins
UNION ALL
SELECT id, 'Interviewing', NOW() - INTERVAL '12 hours' FROM ins;

-- A3: Applied 2 days ago, Interviewing 1.5 days ago, Offered 1 day ago, Replied 1 day ago (counts for 7d applied, 7d replied, 7d interviewing, 7d offered)
WITH ins AS (
  INSERT INTO applications (company_name, role_title, status, applied_at, last_communication, users_id)
  VALUES ('7d Company C', 'Backend Engineer', 'Offered', NOW() - INTERVAL '2 days', NOW() - INTERVAL '1 day', 1)
  RETURNING id
)
INSERT INTO status_history (application_id, status, changed_at)
SELECT id, 'Applied', NOW() - INTERVAL '2 days' FROM ins
UNION ALL
SELECT id, 'Interviewing', NOW() - INTERVAL '1.5 days' FROM ins
UNION ALL
SELECT id, 'Offered', NOW() - INTERVAL '1 day' FROM ins;

-- A4: Applied 2 days ago, Rejected 1 day ago (counts for 7d applied, 7d rejected)
WITH ins AS (
  INSERT INTO applications (company_name, role_title, status, applied_at, users_id)
  VALUES ('7d Company D', 'QA Engineer', 'Rejected', NOW() - INTERVAL '2 days', 1)
  RETURNING id
)
INSERT INTO status_history (application_id, status, changed_at)
SELECT id, 'Applied', NOW() - INTERVAL '2 days' FROM ins
UNION ALL
SELECT id, 'Rejected', NOW() - INTERVAL '1 day' FROM ins;

-- A5: Applied 2 days ago, Ghosted 1 day ago (counts for 7d applied, 7d ghosted)
WITH ins AS (
  INSERT INTO applications (company_name, role_title, status, applied_at, users_id)
  VALUES ('7d Company E', 'Data Scientist', 'Ghosted', NOW() - INTERVAL '2 days', 1)
  RETURNING id
)
INSERT INTO status_history (application_id, status, changed_at)
SELECT id, 'Applied', NOW() - INTERVAL '2 days' FROM ins
UNION ALL
SELECT id, 'Ghosted', NOW() - INTERVAL '1 day' FROM ins;


-- ============================================================================
-- GROUP B: 1m window but NOT 7d (Earlier in current month, current year)
-- ============================================================================

-- B1: Applied 5 days ago (counts for 1m applied)
WITH ins AS (
  INSERT INTO applications (company_name, role_title, status, applied_at, users_id)
  VALUES ('1m Company A', 'Mobile Engineer', 'Applied', NOW() - INTERVAL '5 days', 1)
  RETURNING id
)
INSERT INTO status_history (application_id, status, changed_at)
SELECT id, 'Applied', NOW() - INTERVAL '5 days' FROM ins;

-- B2: Applied 6 days ago, Interviewing 5 days ago, Replied 5 days ago (counts for 1m applied, 1m replied, 1m interviewing)
WITH ins AS (
  INSERT INTO applications (company_name, role_title, status, applied_at, last_communication, users_id)
  VALUES ('1m Company B', 'DevOps Engineer', 'Interviewing', NOW() - INTERVAL '6 days', NOW() - INTERVAL '5 days', 1)
  RETURNING id
)
INSERT INTO status_history (application_id, status, changed_at)
SELECT id, 'Applied', NOW() - INTERVAL '6 days' FROM ins
UNION ALL
SELECT id, 'Interviewing', NOW() - INTERVAL '5 days' FROM ins;

-- B3: Applied 6 days ago, Interviewing 5 days ago, Offered 4 days ago (counts for 1m applied, 1m interviewing, 1m offered)
WITH ins AS (
  INSERT INTO applications (company_name, role_title, status, applied_at, users_id)
  VALUES ('1m Company C', 'Security Specialist', 'Offered', NOW() - INTERVAL '6 days', 1)
  RETURNING id
)
INSERT INTO status_history (application_id, status, changed_at)
SELECT id, 'Applied', NOW() - INTERVAL '6 days' FROM ins
UNION ALL
SELECT id, 'Interviewing', NOW() - INTERVAL '5 days' FROM ins
UNION ALL
SELECT id, 'Offered', NOW() - INTERVAL '4 days' FROM ins;

-- B4: Applied 6 days ago, Rejected 4 days ago (counts for 1m applied, 1m rejected)
WITH ins AS (
  INSERT INTO applications (company_name, role_title, status, applied_at, users_id)
  VALUES ('1m Company D', 'Product Manager', 'Rejected', NOW() - INTERVAL '6 days', 1)
  RETURNING id
)
INSERT INTO status_history (application_id, status, changed_at)
SELECT id, 'Applied', NOW() - INTERVAL '6 days' FROM ins
UNION ALL
SELECT id, 'Rejected', NOW() - INTERVAL '4 days' FROM ins;

-- B5: Applied 6 days ago, Ghosted 4 days ago (counts for 1m applied, 1m ghosted)
WITH ins AS (
  INSERT INTO applications (company_name, role_title, status, applied_at, users_id)
  VALUES ('1m Company E', 'Designer', 'Ghosted', NOW() - INTERVAL '6 days', 1)
  RETURNING id
)
INSERT INTO status_history (application_id, status, changed_at)
SELECT id, 'Applied', NOW() - INTERVAL '6 days' FROM ins
UNION ALL
SELECT id, 'Ghosted', NOW() - INTERVAL '4 days' FROM ins;


-- ============================================================================
-- GROUP C: 1y window but NOT 1m (Earlier in current year)
-- ============================================================================

-- C1: Applied 15 days ago (counts for 1y applied)
WITH ins AS (
  INSERT INTO applications (company_name, role_title, status, applied_at, users_id)
  VALUES ('1y Company A', 'Solutions Architect', 'Applied', NOW() - INTERVAL '15 days', 1)
  RETURNING id
)
INSERT INTO status_history (application_id, status, changed_at)
SELECT id, 'Applied', NOW() - INTERVAL '15 days' FROM ins;

-- C2: Applied 20 days ago, Interviewing 18 days ago, Replied 18 days ago (counts for 1y applied, 1y replied, 1y interviewing)
WITH ins AS (
  INSERT INTO applications (company_name, role_title, status, applied_at, last_communication, users_id)
  VALUES ('1y Company B', 'Cloud Engineer', 'Interviewing', NOW() - INTERVAL '20 days', NOW() - INTERVAL '18 days', 1)
  RETURNING id
)
INSERT INTO status_history (application_id, status, changed_at)
SELECT id, 'Applied', NOW() - INTERVAL '20 days' FROM ins
UNION ALL
SELECT id, 'Interviewing', NOW() - INTERVAL '18 days' FROM ins;

-- C3: Applied 20 days ago, Interviewing 18 days ago, Offered 17 days ago (counts for 1y applied, 1y interviewing, 1y offered)
WITH ins AS (
  INSERT INTO applications (company_name, role_title, status, applied_at, users_id)
  VALUES ('1y Company C', 'Database Administrator', 'Offered', NOW() - INTERVAL '20 days', 1)
  RETURNING id
)
INSERT INTO status_history (application_id, status, changed_at)
SELECT id, 'Applied', NOW() - INTERVAL '20 days' FROM ins
UNION ALL
SELECT id, 'Interviewing', NOW() - INTERVAL '18 days' FROM ins
UNION ALL
SELECT id, 'Offered', NOW() - INTERVAL '17 days' FROM ins;

-- C4: Applied 20 days ago, Rejected 17 days ago (counts for 1y applied, 1y rejected)
WITH ins AS (
  INSERT INTO applications (company_name, role_title, status, applied_at, users_id)
  VALUES ('1y Company D', 'Systems Analyst', 'Rejected', NOW() - INTERVAL '20 days', 1)
  RETURNING id
)
INSERT INTO status_history (application_id, status, changed_at)
SELECT id, 'Applied', NOW() - INTERVAL '20 days' FROM ins
UNION ALL
SELECT id, 'Rejected', NOW() - INTERVAL '17 days' FROM ins;

-- C5: Applied 20 days ago, Ghosted 17 days ago (counts for 1y applied, 1y ghosted)
WITH ins AS (
  INSERT INTO applications (company_name, role_title, status, applied_at, users_id)
  VALUES ('1y Company E', 'Fullstack Engineer', 'Ghosted', NOW() - INTERVAL '20 days', 1)
  RETURNING id
)
INSERT INTO status_history (application_id, status, changed_at)
SELECT id, 'Applied', NOW() - INTERVAL '20 days' FROM ins
UNION ALL
SELECT id, 'Ghosted', NOW() - INTERVAL '17 days' FROM ins;
