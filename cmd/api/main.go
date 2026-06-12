package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zrotrasukha/jobman/internal/data"
	"github.com/zrotrasukha/jobman/internal/mailer"
)

var version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		minIdleConns int
		maxIdleTime  time.Duration
	}
	mailer struct {
		sender   string
		host     string
		port     int
		username string
		password string
	}
}

type application struct {
	config config
	logger *slog.Logger
	models data.Models
	wg     sync.WaitGroup
	mailer *mailer.Mailer
}

// openDB establishes a connection pool to the PostgreSQL database using the provided configuration. It sets the maximum number of open connections, minimum number of idle connections, and maximum idle time for the connections in the pool. The function also performs a ping to the database to ensure that the connection is valid before returning the pool. If any errors occur during this process, they are returned to the caller.
func openDB(cfg config) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	poolCfg.MaxConns = int32(cfg.db.maxOpenConns)
	poolCfg.MinIdleConns = int32(cfg.db.minIdleConns)
	poolCfg.MaxConnIdleTime = cfg.db.maxIdleTime

	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil

}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "environment", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.minIdleConns, "db-min-idle-conns", 25, "PostgreSQL min idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	flag.StringVar(&cfg.mailer.sender, "mailer-sender", os.Getenv("MAILER_SENDER"), "Mailer sender")
	flag.StringVar(&cfg.mailer.host, "mailer-host", os.Getenv("MAILER_HOST"), "Mailer host")
	flag.IntVar(&cfg.mailer.port, "mailer-port", 587, "Mailer port")
	flag.StringVar(&cfg.mailer.username, "mailer-username", os.Getenv("MAILER_USERNAME"), "Mailer username")
	flag.StringVar(&cfg.mailer.password, "mailer-password", os.Getenv("MAILER_PASSWORD"), "Mailer password")

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		log.Printf("Version:\t%s\n", version)
		os.Exit(0)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	logger.Info("database connection pool established")

	mailer, err := mailer.New(cfg.mailer.host, cfg.mailer.port, cfg.mailer.username, cfg.mailer.password, cfg.mailer.sender)
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("mailer initialized")

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		wg:     sync.WaitGroup{},
		mailer: mailer,
	}

	err = app.serve()
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
}
