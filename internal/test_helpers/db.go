package testhelpers

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func NewTestPool() (*pgxpool.Pool, func()) {
	ctx := context.Background()

	dbName := "jobman"
	dbUser := "test"
	dbPassword := "$ecreT_T"

	pgCtr, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		log.Fatal("error occured running postgres container: ", err)
	}

	connStr, err := pgCtr.ConnectionString(ctx, "sslmode=disable")
	migrationURL := "pgx5://" + connStr[len("postgres://"):]

	if err != nil {
		log.Fatal("Invalid or no connection string:", err)
	}

	m, err := migrate.New("file://./../../migrations/", migrationURL)
	if err != nil {
		log.Fatal("Failed to create migrate instance:", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Failed to run migrations:", err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatal("Failed to create connection pool:", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	terminate := func() {
		pool.Close()
		pgCtr.Terminate(ctx)
	}

	return pool, terminate
}

func ClearApplications(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	query := `TRUNCATE applications RESTART IDENTITY`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if _, err := pool.Exec(ctx, query); err != nil {
		t.Fatalf("Failed to clear applications table: %v", err)
	}
}
