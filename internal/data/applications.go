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
	Update(jobApp *JobApplication) error
	Delete(id int64) error
}

type JobApplicationModel struct {
	pool *pgxpool.Pool
}

func (m JobApplicationModel) Insert(jobApp *JobApplication) error {
	query := `INSERT INTO applications (company_name, role_title, applied_at, status, notes)
						VALUES ($1, $2, $3, $4, $5) RETURNING id, version, applied_at, updated_at`

	args := []any{
		jobApp.CompanyName,
		jobApp.RoleTitle,
		jobApp.AppliedAt,
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

func (m JobApplicationModel) Update(jobApp *JobApplication) error {
	query := `UPDATE applications
	SET company_name = $1, role_title = $2, status = $3, applied_at = $4, last_communication = $5, notes = $6, updated_at = now(), version = version + 1
	WHERE id = $7 AND version = $8
	RETURNING version, updated_at`

	args := []any{
		jobApp.CompanyName,
		jobApp.RoleTitle,
		jobApp.Status,
		jobApp.AppliedAt,
		jobApp.LastCommunication,
		jobApp.Notes,
		jobApp.ID,
		jobApp.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.pool.QueryRow(ctx, query, args...).Scan(&jobApp.Version, &jobApp.UpdatedAt)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m JobApplicationModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM applications WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrRecordNotFound
	}

	return err
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
