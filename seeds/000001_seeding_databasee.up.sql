INSERT INTO public.applications (company_name, role_title, status, applied_at, updated_at, last_communication, notes, version, interview_at, stale_after)
VALUES
-- GROUP 1: Brand New Applications (Active, within 30-day window)
('Google', 'Software Engineer', 'applied', NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days', NULL, 'Applied via referral', 1, NULL, NOW() - INTERVAL '2 days' + INTERVAL '30 days'),
('Meta', 'Production Engineer', 'applied', NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days', NULL, 'Cold outreach on LinkedIn', 1, NULL, NOW() - INTERVAL '5 days' + INTERVAL '30 days'),
('Netflix', 'Senior Backend Engineer', 'applied', NOW() - INTERVAL '12 days', NOW() - INTERVAL '12 days', NOW() - INTERVAL '12 days', 'Recruiter reached out', 1, NULL, NOW() - INTERVAL '12 days' + INTERVAL '30 days'),
('Apple', 'ICT3 Engineer', 'applied', NOW() - INTERVAL '20 days', NOW() - INTERVAL '20 days', NULL, 'Submitted resume on portal', 1, NULL, NOW() - INTERVAL '20 days' + INTERVAL '30 days'),
-- GROUP 2: Applications Past 30 Days (Stale / Waiting to be Ghosted by Worker)
('Stripe', 'Full Stack Developer', 'applied', NOW() - INTERVAL '35 days', NOW() - INTERVAL '35 days', NULL, 'No response yet', 1, NULL, NOW() - INTERVAL '35 days' + INTERVAL '30 days'),
('Uber', 'Systems Engineer', 'applied', NOW() - INTERVAL '40 days', NOW() - INTERVAL '40 days', NOW() - INTERVAL '40 days', 'Initial screening completed, then silence', 1, NULL, NOW() - INTERVAL '40 days' + INTERVAL '30 days'),
('Airbnb', 'Frontend Engineer', 'applied', NOW() - INTERVAL '31 days', NOW() - INTERVAL '31 days', NULL, 'Automated confirmation received', 1, NULL, NOW() - INTERVAL '31 days' + INTERVAL '30 days'),
-- GROUP 3: Active Interviews (Upcoming, within grace period)
('Amazon', 'Software Development Engineer II', 'Interviewing', NOW() - INTERVAL '10 days', NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days', 'Technical phone screen booked', 2, NOW() + INTERVAL '3 days', NOW() + INTERVAL '3 days' + INTERVAL '5 days'),
('Microsoft', 'Cloud Architect', 'Interviewing', NOW() - INTERVAL '15 days', NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day', 'Loop scheduled', 3, NOW() + INTERVAL '1 day', NOW() + INTERVAL '1 day' + INTERVAL '5 days'),
('Databricks', 'Distributed Systems Engineer', 'Interviewing', NOW() - INTERVAL '4 days', NOW() - INTERVAL '4 days', NOW() - INTERVAL '4 days', 'Onsite next week', 1, NOW() + INTERVAL '6 days', NOW() + INTERVAL '6 days' + INTERVAL '5 days'),
('Snowflake', 'Database Engineer', 'Interviewing', NOW() - INTERVAL '8 days', NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days', 'Hiring Manager chat done', 2, NOW() + INTERVAL '2 days', NOW() + INTERVAL '2 days' + INTERVAL '5 days'),
-- GROUP 4: Missed/Past Interview Grace Period (Stale / Waiting to be Ghosted by Worker)
('Coinbase', 'Security Engineer', 'Interviewing', NOW() - INTERVAL '25 days', NOW() - INTERVAL '12 days', NOW() - INTERVAL '12 days', 'Interview was 6 days ago, no follow-up', 2, NOW() - INTERVAL '6 days', NOW() - INTERVAL '6 days' + INTERVAL '5 days'),
('Figma', 'Product Designer', 'Interviewing', NOW() - INTERVAL '20 days', NOW() - INTERVAL '10 days', NOW() - INTERVAL '10 days', 'Technical panel was last week', 2, NOW() - INTERVAL '7 days', NOW() - INTERVAL '7 days' + INTERVAL '5 days'),
('Slack', 'Desktop Engineer', 'Interviewing', NOW() - INTERVAL '30 days', NOW() - INTERVAL '15 days', NOW() - INTERVAL '15 days', 'Final round completed 10 days ago', 3, NOW() - INTERVAL '10 days', NOW() - INTERVAL '10 days' + INTERVAL '5 days'),
-- GROUP 5: Already Manually or Process-Archived States (Should NOT be touched by the Staleness Worker)
('Canva', 'Full Stack Engineer', 'rejected', NOW() - INTERVAL '50 days', NOW() - INTERVAL '45 days', NOW() - INTERVAL '45 days', 'Received rejection email', 2, NULL, NOW() - INTERVAL '50 days' + INTERVAL '30 days'),
('Linear', 'Product Engineer', 'offered', NOW() - INTERVAL '14 days', NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day', 'Received written offer details!', 4, NOW() - INTERVAL '4 days', NOW() - INTERVAL '4 days' + INTERVAL '5 days'),
('Vercel', 'DevRel Engineer', 'rejected', NOW() - INTERVAL '45 days', NOW() - INTERVAL '40 days', NOW() - INTERVAL '40 days', 'Application closed by company', 2, NOW() - INTERVAL '42 days', NOW() - INTERVAL '42 days' + INTERVAL '5 days'),
('Supabase', 'Database Advocate', 'withdrawn', NOW() - INTERVAL '18 days', NOW() - INTERVAL '16 days', NOW() - INTERVAL '16 days', 'Withdrew due to location constraints', 2, NULL, NOW() - INTERVAL '18 days' + INTERVAL '30 days'),
-- GROUP 6: Edge Cases / Pre-Existing Ghosted States
('Palantir', 'Forward Deployed Engineer', 'ghosted', NOW() - INTERVAL '60 days', NOW() - INTERVAL '30 days', NULL, 'Automatically ghosted in previous run', 2, NULL, NOW() - INTERVAL '60 days' + INTERVAL '30 days'),
('Notion', 'Backend Infrastructure', 'ghosted', NOW() - INTERVAL '55 days', NOW() - INTERVAL '20 days', NOW() - INTERVAL '25 days', 'Interview window passed and flagged', 3, NOW() - INTERVAL '25 days', NOW() - INTERVAL '25 days' + INTERVAL '5 days');

