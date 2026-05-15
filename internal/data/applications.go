package data

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zrotrasukha/GO---Job-Application-Manager/internal/validator"
)

type Application struct {
	ID                int64  `json:"id"`
	CompanyName       string `json:"company_name"`
	RoleTitle         string `json:"role_title"`
	Status            string `json:"status"`
	AppliedAt         string `json:"applied_at"`
	UpdatedAt         string `json:"updated_at"`
	LastCommunication string `json:"last_communication"`
	Notes             string `json:"notes"`
	Version           int32  `json:"version"` // needed for optimistic locking
}

type ApplicationModel struct {
	pool *pgxpool.Pool
}

func (m ApplicationModel) Insert(jobApp *Application) error {
	query := `INSERT INTO applications (company_name, role_title, status, applied_at, updated_at, last_communication, notes)
						VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	args := []any{
		jobApp.CompanyName,
		jobApp.RoleTitle,
		jobApp.Status,
		jobApp.AppliedAt,
		jobApp.UpdatedAt,
		jobApp.LastCommunication,
		jobApp.Notes,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.pool.QueryRow(ctx, query, args...).Scan(&jobApp.ID)
}

func ValidateApplication(v *validator.Validator, jobApp *Application) {
	v.CheckField(jobApp.CompanyName != "", "company_name", "must be provided")
	v.CheckField(len(jobApp.CompanyName) <= 200, "company_name", "must not be more than 200 bytes long")

}
