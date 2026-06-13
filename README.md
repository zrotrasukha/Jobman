# Jobman

I built this after finishing another project [movapi](https://github.com/zrotrasukha/movapi) when I wanted to start applying to jobs aggressively. I needed something to track applications, but more than that I wanted a real reason to dig into goroutines and background workers, tasty topics. And Jobman gave me reason.
The backend is designed to be consumed by both a CLI client and a web frontend. It's not finished the pinging system is in progress, the digest is queued, and the TUI hasn't been started but the foundation is solid and I'm genuinely enjoying building it.

---

## What it does

Track job applications through their lifecycle: Applied → Interviewing → Offered / Rejected / Ghosted.

- **Staleness worker** — runs every 6 hours, automatically marks applications as Ghosted when they've gone cold. Done.
- **Pinging system** — reminds you about upcoming interviews on a halving schedule. In progress.
- **Weekly digest** — aggregate analytics about your job search delivered on a schedule. Queued.

---

## Stack

Go, PostgreSQL, pgx/v5, chi, golang-migrate, testcontainers-go.

---

## Running it

```bash
# dependencies: Go 1.22+, PostgreSQL, Docker (for integration tests)

git clone https://github.com/zrotrasukha/jobman
cd jobman

make migrate/up
make run/api
```

Environment variables:

```
JOBMAN_DB_DSN
JOBMAN_PORT                 (default 4000)
JOBMAN_ENV                  (development/production)
JOBMAN_SMTP_HOST
JOBMAN_SMTP_PORT
JOBMAN_SMTP_USERNAME
JOBMAN_SMTP_PASSWORD
JOBMAN_SMTP_SENDER
JOBMAN_STALENESS_INTERVAL   (default 6h)
```

---

## API

```
POST   /v1/users                        Register
PUT    /v1/users/activate               Activate account
POST   /v1/tokens/authentication        Login

POST   /v1/applications                 Create application
GET    /v1/applications                 List (search, filter, paginate, sort)
GET    /v1/applications/:id             Get one
PATCH  /v1/applications/:id             Update
DELETE /v1/applications/:id             Delete
```

---

## Testing

Handler tests use mocks for fast isolated feedback. Data layer tests run against real PostgreSQL via testcontainers — one container per test run, shared across all tests, truncated between cases.

```bash
make test/unit          # mock-based handler tests
make test/integration   # requires Docker
```

---

## Status

- [x] CRUD with full-text search, pagination, sorting
- [x] Optimistic locking
- [x] Staleness detection worker
- [x] User auth with email verification
- [ ] Pinging system
- [ ] Weekly digest
- [ ] TUI / CLI client
- [ ] Deployment
