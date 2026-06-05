package testhelpers

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	SharedPool *pgxpool.Pool
)

// These constants are being used in truncating tables
var (
	TableApplications = "applications"
	TableUsers        = "users"
)

func GetPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	return SharedPool
}

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
		postgres.WithSQLDriver("pgx5"),
	)
	if err != nil {
		log.Fatal("error occured running postgres container: ", err)
	}

	connStr, err := pgCtr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatal("Invalid or no connection string:", err)
	}

	superConn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer superConn.Close(ctx)

	// Note to myself: I have done this because the citext extension can only be created by a superuser,
	// the testcontainers-go postgres module gives us a superuser by default, but in the real app, we have
	// jobman role which has limitations, To bypass this issue, the query to create the citext extension is
	// being executed here before running the migrations. Makkhan tests, makkhan main app.
	_, err = superConn.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS citext;")
	if err != nil {
		log.Fatal("Failed to create citext extension:", err)
	}

	migrationURL := "pgx5://" + connStr[len("postgres://"):]
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

func TruncateTable(t *testing.T, pool *pgxpool.Pool, tableName string) {
	t.Helper()
	query := fmt.Sprintf(`TRUNCATE %s RESTART IDENTITY`, tableName)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if _, err := pool.Exec(ctx, query); err != nil {
		t.Fatalf("Failed to clear applications table: %v", err)
	}
}
