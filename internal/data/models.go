package data

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrRecordNotFound = errors.New("record not found")

type Models struct {
	Application ApplicationModel
}

func NewModels(pool *pgxpool.Pool) Models {
	return Models{
		Application: ApplicationModel{pool: pool},
	}
}
