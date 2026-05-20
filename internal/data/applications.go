package data

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zrotrasukha/jobman/internal/validator"
)

type JobApplication struct {
	ID                int64      `json:"id"`
	CompanyName       string     `json:"company_name"`
	RoleTitle         string     `json:"role_title"`
	Status            Status     `json:"status"`
	AppliedAt         *time.Time `json:"applied_at"`
	UpdatedAt         time.Time  `json:"updated_at"` // cannot be null
	LastCommunication *time.Time `json:"last_communication"`
	Notes             string     `json:"notes"`
	Version           int32      `json:"version"` // needed for optimistic locking
}

type JobApplicationModelInterface interface {
	Insert(jobApp *JobApplication) error
	Get(id int64) (*JobApplication, error)
}

type JobApplicationModel struct {
	pool *pgxpool.Pool
}

func (m JobApplicationModel) Insert(jobApp *JobApplication) error {
	query := `INSERT INTO applications (company_name, role_title, status, notes)
						VALUES ($1, $2, $3, $4) RETURNING id, version, applied_at, updated_at`

	args := []any{
		jobApp.CompanyName,
		jobApp.RoleTitle,
		jobApp.Status,
		jobApp.Notes,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.pool.QueryRow(ctx, query, args...).Scan(
		&jobApp.ID,
		&jobApp.Version,
		&jobApp.AppliedAt,
		&jobApp.UpdatedAt,
	)
}

func (m JobApplicationModel) Get(id int64) (*JobApplication, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, company_name, role_title, status, applied_at, updated_at, last_communication, notes, version
						FROM applications
						WHERE id = $1
	`

	var jobApp JobApplication

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.pool.QueryRow(ctx, query, id).Scan(
		&jobApp.ID,
		&jobApp.CompanyName,
		&jobApp.RoleTitle,
		&jobApp.Status,
		&jobApp.AppliedAt,
		&jobApp.UpdatedAt,
		&jobApp.LastCommunication,
		&jobApp.Notes,
		&jobApp.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &jobApp, nil
}

func ValidateJobApplication(v *validator.Validator, jobApp *JobApplication) {
	v.CheckField(jobApp.CompanyName != "", "company_name", "must be provided")
	v.CheckField(len(jobApp.CompanyName) <= 200, "company_name", "must not be more than 200 bytes long")

	v.CheckField(jobApp.RoleTitle != "", "role_title", "must be provided")
	v.CheckField(len(jobApp.RoleTitle) <= 200, "role_title", "must not be more than 200 bytes long")

	v.CheckField(jobApp.Status != "", "status", "must be provided")
	v.CheckField(len(jobApp.Status) <= 200, "status", "must not be more than 200 bytes long")
	v.CheckField(jobApp.Status.IsValid(), "status", "must be a valid status")

	v.CheckField(jobApp.AppliedAt != nil, "applied_at", "must be provided")
	v.CheckField(len(jobApp.Notes) <= 8000, "notes", "must not be more than 1000 bytes long")
}
