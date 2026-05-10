package data

import "github.com/jackc/pgx/v5/pgxpool"

type Models struct {
	Application ApplicationModel
}

func NewModels(pool *pgxpool.Pool) Models {
	return Models{
		Application: ApplicationModel{pool: pool},
	}
}
