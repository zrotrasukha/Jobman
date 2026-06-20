package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zrotrasukha/jobman/internal/validator"
)

// JobApplication represents a job application record in the database. It includes fields for the company name, role title, application status, timestamps for when the application was made and last updated, and any notes about the application. The Version field is used for optimistic locking to prevent concurrent update conflicts.
type JobApplication struct {
	ID                int64      `json:"id"`
	UserID            int64      `json:"-"`
	CompanyName       string     `json:"company_name"`
	RoleTitle         string     `json:"role_title"`
	AppliedAt         time.Time  `json:"applied_at"`
	Status            Status     `json:"status"`
	UpdatedAt         time.Time  `json:"updated_at"` // cannot be null
	InterviewAt       *time.Time `json:"interview_at"`
	StaleAfter        *time.Time `json:"stale_after"` // it is possible for this to be null if the user doesn't want the application to ever be marked as stale
	LastCommunication *time.Time `json:"last_communication"`
	Notes             string     `json:"notes"`
	Version           int32      `json:"version"` // needed for optimistic locking
}

// JobApplicationModelInterface defines the methods that any implementation of a job application model must provide. This includes methods for inserting a new job application, retrieving a job application by ID, retrieving all job applications with optional search and filtering, updating an existing job application, and deleting a job application by ID.
type JobApplicationModelInterface interface {
	Insert(jobApp *JobApplication) error
	Get(id int64, userID int64) (*JobApplication, error)
	GetAll(searchString string, filters Filters, userID int64) ([]*JobApplication, *Metadata, error)
	Update(jobApp *JobApplication, userID int64) error
	Delete(id int64, userID int64) error
	MarkStaleApplications(ctx context.Context) (int64, error)
}

// JobApplicationModel provides methods for interacting with the job applications table in the database. It uses a connection pool to execute SQL queries and manage database connections efficiently.
type JobApplicationModel struct {
	pool *pgxpool.Pool
}

// Insert adds a new job application record to the database. It takes a JobApplication struct as input and populates its ID, Version, AppliedAt, and UpdatedAt fields based on the values returned from the database after the insert operation. If the insert is successful, it returns nil; otherwise, it returns an error.
func (m JobApplicationModel) Insert(jobApp *JobApplication) error {
	query := `INSERT INTO applications (users_id, company_name, role_title, applied_at, status, notes, interview_at, stale_after)
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, version, applied_at, updated_at, interview_at, stale_after`

	args := []any{
		jobApp.UserID,
		jobApp.CompanyName,
		jobApp.RoleTitle,
		jobApp.AppliedAt,
		jobApp.Status,
		jobApp.Notes,
		jobApp.InterviewAt,
		jobApp.StaleAfter,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.pool.QueryRow(ctx, query, args...).Scan(
		&jobApp.ID,
		&jobApp.Version,
		&jobApp.AppliedAt,
		&jobApp.UpdatedAt,
		&jobApp.InterviewAt,
		&jobApp.StaleAfter,
	)
}

// Get returns the job application with the specified ID. If no matching record is found, it returns a ErrRecordNotFound error.
func (m JobApplicationModel) Get(id int64, userID int64) (*JobApplication, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, users_id, company_name, role_title, status, interview_at, stale_after, applied_at,  updated_at, last_communication, notes, version
						FROM applications
						WHERE users_id = $1 AND id = $2`

	var jobApp JobApplication

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		userID,
		id,
	}
	err := m.pool.QueryRow(ctx, query, args...).Scan(
		&jobApp.ID,
		&jobApp.UserID,
		&jobApp.CompanyName,
		&jobApp.RoleTitle,
		&jobApp.Status,
		&jobApp.InterviewAt,
		&jobApp.StaleAfter,
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

// GetAll returns a list of job applications matching the provided searchString and filters. It also returns metadata about the total number of records and pagination details.
func (m JobApplicationModel) GetAll(searchString string, filters Filters, userID int64) ([]*JobApplication, *Metadata, error) {
	query := fmt.Sprintf(`
					SELECT COUNT(*) OVER() as total_records, id, company_name, role_title, status, interview_at, stale_after, applied_at, updated_at, last_communication, notes, version
					FROM applications
					WHERE users_id = $1 AND
					(to_tsvector('simple', company_name || ' ' || role_title) @@ plainto_tsquery('simple', $2) or $2 = '')
					ORDER BY %s %s, id asc
					LIMIT $3 OFFSET $4`, filters.SortColumn(), filters.SortDirection())

	args := []any{
		userID,
		searchString,
		filters.PageSize,
		filters.Offset(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	totalRecords := 0
	jobApps := []*JobApplication{}

	for rows.Next() {
		var jobApp JobApplication

		err := rows.Scan(
			&totalRecords,
			&jobApp.ID,
			&jobApp.CompanyName,
			&jobApp.RoleTitle,
			&jobApp.Status,
			&jobApp.InterviewAt,
			&jobApp.StaleAfter,
			&jobApp.AppliedAt,
			&jobApp.UpdatedAt,
			&jobApp.LastCommunication,
			&jobApp.Notes,
			&jobApp.Version,
		)
		if err != nil {
			return nil, nil, err
		}

		jobApps = append(jobApps, &jobApp)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return jobApps, metadata, nil
}

// Update modifies the details of an existing job application in the database. It uses optimistic locking to ensure that updates are only applied if the record has not been modified by another process since it was last read. If a version conflict is detected, it returns an ErrEditConflict error.
func (m JobApplicationModel) Update(jobApp *JobApplication, userID int64) error {
	query := `UPDATE applications
						SET company_name = $1, role_title = $2, status = $3, interview_at = $4, stale_after = $5, applied_at = $6, last_communication = $7, notes = $8, updated_at = now(), version = version + 1
						WHERE users_id=$9 AND id = $10 AND version = $11
						RETURNING version, updated_at`

	args := []any{
		jobApp.CompanyName,
		jobApp.RoleTitle,
		jobApp.Status,
		jobApp.InterviewAt,
		jobApp.StaleAfter,
		jobApp.AppliedAt,
		jobApp.LastCommunication,
		jobApp.Notes,
		jobApp.UserID,
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

// Delete removes the job application with the specified ID from the database. If no matching record is found, it returns a ErrRecordNotFound error.
func (m JobApplicationModel) Delete(id int64, userID int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM applications WHERE users_id = $1 AND id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		userID,
		id,
	}

	result, err := m.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrRecordNotFound
	}

	return err
}

func (m JobApplicationModel) MarkStaleApplications(ctx context.Context) (int64, error) {
	query := `UPDATE applications
						SET status = 'Ghosted'
						WHERE stale_after < NOW()
						AND status IN ('Applied', 'Interviewing')`

	result, err := m.pool.Exec(ctx, query)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

// ValidateJobApplication checks the fields of a JobApplication struct to ensure they meet the required criteria. It uses the provided validator to collect any validation errors.
func ValidateJobApplication(v *validator.Validator, jobApp *JobApplication) {
	v.CheckField(jobApp.CompanyName != "", "company_name", "must be provided")
	v.CheckField(len(jobApp.CompanyName) <= 200, "company_name", "must not be more than 200 bytes long")

	v.CheckField(jobApp.RoleTitle != "", "role_title", "must be provided")
	v.CheckField(len(jobApp.RoleTitle) <= 200, "role_title", "must not be more than 200 bytes long")

	v.CheckField(jobApp.Status != "", "status", "must be provided")
	v.CheckField(len(jobApp.Status) <= 200, "status", "must not be more than 200 bytes long")
	v.CheckField(jobApp.Status.IsValid(), "status", "must be a valid status")

	v.CheckField(!jobApp.AppliedAt.IsZero(), "applied_at", "must be provided")
	v.CheckField(len(jobApp.Notes) <= 8000, "notes", "must not be more than 1000 bytes long")

	// interview_at gotta be ahead the date of application or the applicant is cooked
	v.CheckField(jobApp.InterviewAt == nil || jobApp.InterviewAt.After(jobApp.AppliedAt), "interview_at", "must be after the applied date")

	// same shit for the stale_after field
	v.CheckField(jobApp.StaleAfter == nil || jobApp.StaleAfter.After(jobApp.AppliedAt), "stale_after", "must be after the applied date")

	// if the status is interviewing, then the interview_at field must be provided
	v.CheckField(
		jobApp.Status != StatusInterviewing || jobApp.InterviewAt != nil,
		"interview_at",
		"must be provided when status is Interviewing",
	)
}
