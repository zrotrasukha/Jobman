package data

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// Models wraps all database models for easy dependency injection.
type Models struct {
	Application JobApplicationModelInterface
	User        UserModelInterface
}

// NewModels initializes and returns a Models struct with all DB repositories.
func NewModels(pool *pgxpool.Pool) Models {
	return Models{
		Application: JobApplicationModel{pool: pool},
		User:        UserModel{pool: pool},
	}
}
