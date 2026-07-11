package data

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zrotrasukha/jobman/internal/validator"
)

// time windows
const (
	Window7d = "7d"
	Window1m = "1m"
	Window1y = "1y"
)

type DigestModel struct {
	pool *pgxpool.Pool
}

type DigestModelInterface interface {
	GetDigest(userId int64, from, to time.Time) (*Digest, error)
}

type Digest struct {
	Window      string       `json:"window"`
	From        time.Time    `json:"from"`
	To          time.Time    `json:"to"`
	Funnel      *Funnel      `json:"funnel"`
	GhostCohort *GhostCohort `json:"all_time_ghost_cohort,omitempty"`
}

type Funnel struct {
	Applied      int `json:"applied"`
	Replied      int `json:"replied"`
	Interviewing int `json:"interviewing"`
	Offered      int `json:"offered"`
	Rejected     int `json:"rejected"`
	Ghosted      int `json:"ghosted"`
}

func (m DigestModel) GetDigest(userId int64, from, to time.Time) (*Digest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := m.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	funnel, err := m.GetFunnel(ctx, tx, userId, from, to)
	if err != nil {
		return nil, err
	}

	ghostCohort, err := m.GetGhostCohort(ctx, tx, userId)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	var digest Digest
	digest.From = from
	digest.To = to
	digest.Funnel = funnel
	digest.GhostCohort = ghostCohort

	return &digest, nil

}

func (m DigestModel) GetFunnel(ctx context.Context, tx pgx.Tx, userId int64, from, to time.Time) (*Funnel, error) {
	// Use RepeatableRead read-only transaction for snapshot consistency

	var funnel Funnel

	// 1. Applied: count of applications created in the window
	appliedQuery := `SELECT COUNT(*) FROM applications
										WHERE users_id = $1
										AND applied_at >= $2
										AND applied_at < $3`

	err := tx.QueryRow(ctx, appliedQuery, userId, from, to).Scan(&funnel.Applied)
	if err != nil {
		return nil, err
	}

	// 2. Replied: count of applications where last_communication was set in the window
	repliesQuery := `SELECT COUNT(*) FROM applications
										WHERE users_id = $1
										AND last_communication IS NOT NULL
										AND last_communication >= $2
										AND last_communication < $3`

	err = tx.QueryRow(ctx, repliesQuery, userId, from, to).Scan(&funnel.Replied)
	if err != nil {
		return nil, err
	}

	// Shared transition query using status_history
	transitionQuery := `SELECT COUNT(*) FROM applications AS a
											INNER JOIN status_history AS sh ON a.id = sh.application_id
											WHERE a.users_id = $1
											AND sh.status = $2
											AND sh.changed_at >= $3
											AND sh.changed_at <= $4`

	// 3. Interviewing transitions
	err = tx.QueryRow(ctx, transitionQuery, userId, StatusInterviewing, from, to).Scan(&funnel.Interviewing)
	if err != nil {
		return nil, err
	}

	// 4. Offered transitions
	err = tx.QueryRow(ctx, transitionQuery, userId, StatusOffered, from, to).Scan(&funnel.Offered)
	if err != nil {
		return nil, err
	}

	// 5. Rejected transitions
	err = tx.QueryRow(ctx, transitionQuery, userId, StatusRejected, from, to).Scan(&funnel.Rejected)
	if err != nil {
		return nil, err
	}

	// 6. Ghosted transitions
	err = tx.QueryRow(ctx, transitionQuery, userId, StatusGhosted, from, to).Scan(&funnel.Ghosted)
	if err != nil {
		return nil, err
	}

	return &funnel, nil
}

func ValidateDigest(v *validator.Validator, window string) {
	v.CheckField(
		window == "7d" || window == "1m" || window == "1y",
		"window", "must be one of 7d, 1m, 1y or all",
	)
}

type GhostCohort struct {
	Matured        int   `json:"matured"`
	MaturedGhosted int   `json:"matured_ghosted"`
	Rate           *Rate `json:"rate"` // nil if matured == 0 — no verdict-able data yet
}

func (m DigestModel) GetGhostCohort(ctx context.Context, tx pgx.Tx, userId int64) (*GhostCohort, error) {
	query := `SELECT
								COUNT(*) FILTER (
										WHERE
												status IN ('Rejected', 'Offered', 'Ghosted')
												OR (
														stale_after IS NOT NULL
														AND stale_after < NOW()
												)
								) AS matured,
								COUNT(*) FILTER (
										WHERE
												status = 'Ghosted' -- Removed the dangling 'AND'
								) AS matured_ghosted -- Removed the trailing comma
						FROM
								applications
						WHERE
								users_id = $1;`

	var cohort GhostCohort
	err := tx.QueryRow(ctx, query, userId).Scan(&cohort.Matured, &cohort.MaturedGhosted)
	if err != nil {
		return nil, err // counts can legitimately be 0, so no sql.ErrNoRows special-case needed
	}

	if cohort.Matured > 0 {
		rate := Rate(float64(cohort.MaturedGhosted) / float64(cohort.Matured) * 100)
		cohort.Rate = &rate
	}

	return &cohort, nil
}
