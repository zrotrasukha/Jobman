package data

import "github.com/jackc/pgx/v5/pgxpool"

type Application struct {
	ID                int64  `json:"id"`
	Company_name      string `json:"company_name"`
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
