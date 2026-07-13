# Phase 5 — Weekly Digest

- [ ] Status: Not started Depends on: Auth (done), applications schema with status/timestamps (done), Redis caching (done, ready to reuse) Supersedes: `reminder_dispatches` plan (deprecated — per-interview drip emails were assessed as poor UX; the halving "urgency" concept may resurface here instead, as a way to flag soon-due interviews within the digest rather than firing separate emails)

---

## Design principle

The digest is honest, not flattering. No softened language, no hiding bad numbers. If the user's reply rate is 8%, the digest says 8%, plainly. Any auto-generated commentary must be derived from the user's actual numbers (observational), never generic encouragement — generic positivity undermines the credibility of an otherwise honest product.

---

## What the digest shows

### 1. Funnel (primary view)

Not flat bullet points — a funnel, since the relationships between stages matter more than any single count in isolation.

```
Applied:        20
Replied:         6   (30%)
Interviewing:    3   (15%)
Offered:         1   (5%)
Rejected:        4   (20%)
Ghosted:        12   (60%)
```

Definitions (decided, exact):

- **Applied** — count of applications created in the window
- **Replied (self-reported)** — `last_communication IS NOT NULL`, set within the window. Labeled explicitly as self-reported in the API response since this is manually entered by the user, not system-observed — see Open Questions for full reasoning.
- **Interviewing / Offered / Rejected / Ghosted** — count of applications that _transitioned into_ that status during the window, per `status_history`. Not a snapshot of current status. This is a deliberate choice — the digest answers "what happened this week," not "what does my pipeline look like right now."

### 2. Ghosting rate, called out explicitly

Don't bury this inside the funnel as just another row — surface it as its own headline number. This is the most "honest" stat available and the one most existing job-tracking tools don't show you plainly.

```
Ghosting rate this week: 60%
(12 of 20 applications received no response at all)
```

### 3. Week-over-week trend

Not just a snapshot — compare current window against the prior window of the same length.

```
Applications sent: 20 (▼ 5 from last week)
Reply rate: 30% (▲ 12% from last week)
```

### 4. Average time-to-reply / time-to-ghost

Texture beyond counts — two users with identical funnel numbers can have very different waiting experiences.

```
Average time to first reply: 5 days
Average time to ghosting: 22 days
```

Computed from `applied_at` → `last_communication` (for replies) and `applied_at` → `stale_after` or the point status flipped to `Ghosted` (for ghosting). Needs `status_history` or equivalent to compute accurately — see Open Questions, this may not be fully buildable in v1 without schema changes.

### 5. (v2) Top N longest-silent applications

Spicier, more emotionally resonant than an aggregate percentage. Optional, build after v1 ships.

```
Still no word from:
- Acme Corp (applied 34 days ago)
- Initech (applied 29 days ago)
- Globex (applied 27 days ago)
```

### 6. (v2) One observational insight line

Auto-generated, but strictly derived from the user's own numbers — never generic. Example of correct vs incorrect framing:

```
Correct (data-driven):
"Your reply rate this week (40%) is double your 4-week average (18%)."

Incorrect (generic, do not build this):
"Keep going, you've got this!"
```

This may be a good candidate for an LLM call later (genuinely useful AI feature, distinct from a chatbot bolted on for its own sake) — feed it the computed numbers, ask for one short factual observation, not encouragement. Flagged as v2, not blocking v1.

---

## What got cut from the original five-item list, and why

Nothing was actually cut — the original list (applications sent, reply backs, interviews planned, offers, rejections) is fully represented in the funnel. The change was structural: presenting them as a connected funnel with percentages rather than five disconnected bullet points, since the gaps between stages are the actually meaningful signal, not the raw counts alone.

---

## Schema considerations

### What's already available, no migration needed

- `applications.status`
- `applications.applied_at`
- `applications.last_communication`
- `applications.stale_after`
- `applications.users_id`

### What's likely missing — status history

DECIDED: build `status_history` now (see Open Questions above for full reasoning and schema). The funnel is transition-based, not current-state, so this table is required, not optional, for v1.

---

## Computation strategy: on-demand vs cached

This is where Redis genuinely earns its keep — unlike the earlier discussion about caching individual `GET /v1/applications/:id` lookups (correctly rejected as unnecessary), the digest involves window functions and aggregations across a user's full application history, which is a meaningfully more expensive query than a primary key lookup.

```
GET /v1/digest
    → check Redis: digest:user:{id}:{window}
    → hit → return cached
    → miss → run aggregation queries (in a transaction) → cache with TTL → return

TTL: 1 hour is reasonable. Digest data doesn't need to be real-time; a user
checking it twice in the same hour should get a consistent, fast answer rather
than recomputing every time.
```

Cache invalidation: not needed beyond TTL expiry. Unlike the token cache (which had a real staleness bug around `Activated` status), digest data becoming slightly stale within an hour is acceptable and matches how digests are typically consumed (checked occasionally, not real-time).

**Decided: no background worker for pre-warming the cache.** A worker that proactively recomputes the digest every hour was considered and rejected — this only makes sense when computation is expensive enough that doing it inside a request would cause noticeable latency. The digest's aggregation queries are a handful of COUNT/AVG queries over one user's application history — not expensive enough to justify a dedicated worker. Cache-aside alone, computed on first request after expiry, is the right level of complexity here. This mirrors the same reasoning that correctly ruled out caching individual application lookups earlier — match the mechanism to the actual cost of computation, not to what sounds more sophisticated.

**Decided: wrap aggregation queries in a transaction.** Not for the write guarantees (this is read-only), but for read consistency — if the staleness worker flips an application's status mid-computation, separate un-transacted queries could see different snapshots of the data and produce a funnel that doesn't internally add up (e.g., ghosting rate computed from a different moment than the funnel counts). Use `RepeatableRead` isolation so every query inside the transaction sees the same snapshot:

```go
tx, err := pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
defer tx.Rollback(ctx)

funnel := computeFunnel(tx, userID, window)
ghostingRate := computeGhostingRate(tx, userID, window)
trend := computeTrend(tx, userID, window)

tx.Commit(ctx) // read-only; commit just closes the transaction cleanly
```

At current scale the odds of an actual race are low, but the correctness discipline costs nothing extra and is worth doing properly from the start.

**Decided: support explicit manual refresh via query param**, distinct from the (rejected) worker-based refresh idea. This lets the user bypass the cache on demand rather than waiting out the TTL, and is transparent about cache state in the response:

```
GET /v1/digest              → cache-aside, may return up to 1h stale data
GET /v1/digest?refresh=true → bypass cache, recompute, overwrite cache entry
```

Response includes a `cached: bool` field so the CLI can display something like `(cached, refreshed 12 min ago)` vs `(just refreshed)` — small transparency touch, cheap to build, keeps the "honest" theme of this feature consistent even down to the caching mechanics.

---

## API surface

```
GET /v1/digest?window=7d
GET /v1/digest?window=7d&refresh=true
```

Query param for window size (default 7 days, weekly). Could later support `30d`, `all-time`, etc. Keep it simple for v1 — just support the weekly case, generalize later only if there's an actual need.

`refresh=true` bypasses the cache, recomputes, and overwrites the cached entry — see caching section above.

Response shape (draft):

```json
{
  "window": "7d",
  "from": "2026-06-24",
  "to": "2026-07-01",
  "cached": false,
  "funnel": {
    "applied": 20,
    "replied": 6,
    "interviewing": 3,
    "offered": 1,
    "rejected": 4,
    "ghosted": 12
  },
  "funnel_basis": "transitions_during_window",
  "ghosting_rate": 0.60,
  "trend": {
    "applications_delta": -5,
    "reply_rate_delta": 0.12
  },
  "averages": {
    "time_to_reply_days": 5,
    "time_to_ghost_days": 22
  }
}
```

Notes on the shape:

- `replied_self_reported` — field named explicitly to signal this metric is a proxy based on manual user input, not system-observed, consistent with the "honest about the data's own limitations" decision above.
- `funnel_basis` — states plainly that counts are transition-based (occurred during this window), not a snapshot of current status. Makes the digest's counting methodology transparent rather than implicit.
- All other funnel fields (`interviewing`, `offered`, `rejected`, `ghosted`) are system-observed via `status_history` and don't need the same caveat.

---

## Open Questions — Decided

**Funnel counting basis — DECIDED: transition-based, not current-state.** "Interviewing: 3" means 3 applications transitioned into `Interviewing` during the window, not 3 currently sitting in that status. This answers "what happened this week," which is the correct frame for a _weekly_ digest. Requires `status_history` (see below). This also clarifies why "Replied" and "Interviewing" aren't actually redundant despite both correlating with forward progress — not every reply leads to an interview (some are outright rejections), so both remain meaningful funnel stages. The real fix for the redundancy concern is presenting them as a funnel with conversion percentages, not as flat disconnected counts.

**status_history — DECIDED: build it now, as a table.**

```sql
-- 000008_create_status_history.up.sql

CREATE TABLE IF NOT EXISTS status_history (
    id bigserial PRIMARY KEY,
    application_id bigint NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    status text NOT NULL,
    changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS status_history_application_id_idx
    ON status_history(application_id);
```

```sql
-- 000008_create_status_history.down.sql

DROP INDEX IF EXISTS status_history_application_id_idx;
DROP TABLE IF EXISTS status_history;
```

One row per transition. Written from `UpdateApplicationHandler` whenever the incoming status differs from the application's current status:

```go
if input.Status != nil && *input.Status != application.Status {
    app.models.StatusHistory.Insert(application.ID, *input.Status)
}
```

Backfill consideration: existing applications (test data, anything created before this migration) will have zero history rows. Consider inserting one synthetic row per existing application (`status = current status, changed_at = applied_at`) so the digest doesn't show a false gap for pre-existing records. Not a hard blocker, but worth doing before relying on this table for real digest output.

**"Replied" definition — DECIDED: keep `last_communication` as the proxy, but label it honestly as self-reported.**

This field is manually entered by the user, which is a permanent limitation, not something more schema can fix — the underlying data source is human input, unlike `status`, which the system observes directly through its own write path. Building more machinery around this (email parsing, IMAP integration — already ruled out earlier in this project) would be solving a data-quality problem with engineering, which doesn't actually fix data quality; only better user habits do.

Decisions:

- `last_communication IS NOT NULL`, set within the window, remains the proxy for v1. No new field, no new table.
- The digest output should label this metric honestly rather than presenting it with the same authority as system-observed facts: `"Replied (self-reported): 6"` rather than just `"Replied: 6"`. Consistent with the digest's stated "brutal honesty" principle — honest about the data's own limitations, not just the numbers themselves.
- Small UX nudge (not a schema change): when a user updates status to `Interviewing` or beyond, prompt them to also set `last_communication` if it's still null, since a reply almost always precedes those transitions. Improves data quality without adding complexity.

**Window boundaries — DECIDED: rolling 7 days from request time.** Confirmed, no change from the original recommendation. Avoids calendar-week timezone edge cases.

---

## TODO

### Schema

- [ ] Migration `000008` — `status_history` table + index (decided: build now)
- [ ] `StatusHistory` model, `Insert(applicationID, status)` method
- [ ] Wire into `UpdateApplicationHandler` — insert a row whenever incoming status differs from current status
- [ ] Backfill consideration: synthetic history rows for pre-existing applications (one row per app, status = current, changed_at = applied_at)
- [ ] If yes: migration, model, write path wired into `UpdateApplicationHandler` (insert a row every time status actually changes)

### Aggregation queries

- [ ] Funnel query — counts per status within window (basis depends on status_history decision above)
- [ ] Ghosting rate calculation
- [ ] Week-over-week trend — query current window + prior window, compute deltas
- [ ] Average time-to-reply — `applied_at` to `last_communication`, only for applications where it's set
- [ ] Average time-to-ghost — `applied_at` to the point status became `Ghosted` (needs status_history for accuracy)

### API

- [ ] `GET /v1/digest` handler, scoped to authenticated user
- [ ] `?window=` query param, default 7d
- [ ] `?refresh=true` query param — bypasses cache, forces recompute
- [ ] Aggregation queries wrapped in a single transaction, RepeatableRead isolation, for internally consistent snapshot
- [ ] Redis caching wired in (cache-aside, same pattern as token caching, 1 hour TTL, graceful fallback if Redis unavailable)
- [ ] Response includes `cached: bool` field

### Testing

- [ ] Integration test: seed applications across various statuses and dates, verify funnel counts are correct
- [ ] Integration test: verify week-over-week trend math against a known two-window dataset
- [ ] Integration test: verify ghosting rate calculation
- [ ] Unit test: average time-to-reply / time-to-ghost calculation logic

### v2 (after v1 ships and is validated)

- [ ] Top N longest-silent applications
- [ ] Observational insight line (possibly LLM-generated, strictly data-derived)
- [ ] Support for window sizes beyond 7d

---

## Order to build in

```
1. status_history migration + StatusHistory model + write-path wiring into
   UpdateApplicationHandler (decided: build now, not optional)
2. Backfill synthetic history rows for any pre-existing applications
3. Funnel query (transition-based, using status_history), tested in
   isolation against seeded data
4. Ghosting rate + trend + averages, layered on top of the funnel query
5. GET /v1/digest handler with ?refresh= support, no caching yet — verify
   correctness first
6. Wrap aggregation queries in a RepeatableRead transaction
7. Redis caching layer on top, once correctness is confirmed
8. v2 features only after v1 is being used and validated
```

---

# New Plan 
Good structure so far — the model/handler split is clean. Here are the three additional window helpers, matching your function's style:

```go
// currentMonthWindow returns the boundaries of the current calendar month,
// from day 1 00:00:00 to the last day 23:59:59.999999999, local time.
func (app *application) currentMonthWindow(t time.Time) (time.Time, time.Time) {
    from := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
    // first day of NEXT month, then step back one nanosecond conceptually —
    // simpler to just go to day 1 of next month and treat `to` as exclusive upstream,
    // but keeping your inclusive style for consistency:
    to := from.AddDate(0, 1, 0).Add(-time.Nanosecond)
    return from, to
}

// currentYearWindow returns the boundaries of the current calendar year,
// from Jan 1 00:00:00 to Dec 31 23:59:59.999999999, local time.
func (app *application) currentYearWindow(t time.Time) (time.Time, time.Time) {
    from := time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, t.Location())
    to := from.AddDate(1, 0, 0).Add(-time.Nanosecond)
    return from, to
}

// allTimeWindow returns a window covering everything from a fixed epoch
// (before Jobman could possibly have any data) to now. There is no
// "previous window" for all-time — see note in ListDigestHandler.
func (app *application) allTimeWindow(t time.Time) (time.Time, time.Time) {
    from := time.Date(2020, time.January, 1, 0, 0, 0, 0, t.Location())
    return from, t
}
```

**Note on `currentMonthWindow`** — using `AddDate(0, 1, 0)` correctly handles variable month lengths (28/30/31 days) without you needing to hardcode day counts. Same idea in `currentYearWindow` with `AddDate(1, 0, 0)` correctly handling leap years automatically.

---

## Wiring these into your handler

```go
func (app *application) ListDigestHandler(w http.ResponseWriter, r *http.Request) {
    user := app.ContextGetUser(r)

    qs := r.URL.Query()
    window := app.readString(qs, "window", "7d")

    v := validator.New()
    v.CheckField(
        window == "7d" || window == "1mo" || window == "1y" || window == "all",
        "window", "must be one of: 7d, 1mo, 1y, all",
    )
    if !v.Valid() {
        app.failedValidationResponse(w, r, v.Errors)
        return
    }

    now := time.Now()
    var from, to time.Time

    switch window {
    case "7d":
        from, to = app.currentWeekWindow(now)
    case "1mo":
        from, to = app.currentMonthWindow(now)
    case "1y":
        from, to = app.currentYearWindow(now)
    case "all":
        from, to = app.allTimeWindow(now)
    }

    funnel, err := app.models.Digest.GetDigest(user.Id, from, to)
    if err != nil {
        app.serverErrResponse(w, r, err)
        return
    }

    response := envelop{
        "window":       window,
        "from":         from.Format("2006-01-02"),
        "to":           to.Format("2006-01-02"),
        "cached":       false,
        "funnel":       funnel,
        "funnel_basis": "transitions_during_window",
    }

    err = app.writeJSON(w, http.StatusOK, response, nil)
    if err != nil {
        app.serverErrResponse(w, r, err)
    }
}
```

---

## The trend/ghosting_rate/averages fields you haven't added yet

Right now `GetDigest` only computes `Funnel`. To get the rest of the JSON shape you designed earlier, you need three more things layered on top — not inside `GetDigest`, but as separate small functions the handler composes:

**1. Ghosting rate — pure math, no new query, derive it from the funnel you already have:**

```go
func ghostingRate(f *data.Funnel) float64 {
    if f.Applied == 0 {
        return 0
    }
    return float64(f.Ghosted) / float64(f.Applied)
}
```

**2. Trend — call `GetDigest` a second time against the previous window, diff the two:**

```go
func (m DigestModel) GetPreviousWindow(from, to time.Time) (time.Time, time.Time) {
    duration := to.Sub(from)
    prevTo := from
    prevFrom := from.Add(-duration)
    return prevFrom, prevTo
}
```

Note — this only makes sense for `7d`, `1mo`, `1y`. For `all`, there's no "previous all-time" to compare against — skip trend computation entirely in that case:

```go
var trend *Trend
if window != "all" {
    prevFrom, prevTo := app.models.Digest.GetPreviousWindow(from, to)
    prevFunnel, err := app.models.Digest.GetDigest(user.Id, prevFrom, prevTo)
    if err != nil {
        app.serverErrResponse(w, r, err)
        return
    }
    trend = &Trend{
        ApplicationsDelta: funnel.Applied - prevFunnel.Applied,
        ReplyRateDelta:    replyRate(funnel) - replyRate(prevFunnel),
    }
}
```

**3. Averages (time-to-reply, time-to-ghost)** — this is a separate query entirely, not derivable from counts. You'll need something like:

```sql
SELECT AVG(EXTRACT(EPOCH FROM (last_communication - applied_at)) / 86400)
FROM applications
WHERE users_id = $1
AND last_communication IS NOT NULL
AND applied_at >= $2 AND applied_at < $3
```

This is its own method (`GetAverages` or similar) — worth building as a separate step once the funnel + trend are confirmed correct, per your original build order in the plan doc.

---

Your `GetDigest` and handler are solid as the funnel foundation. The pattern to hold onto: **funnel** is one query set, **trend** is the same funnel query run twice and diffed, **ghosting rate** is pure arithmetic on the funnel you already have, and **averages** is a genuinely separate query. Build and test each layer independently before combining them, exactly as your original plan doc's build order specified.
