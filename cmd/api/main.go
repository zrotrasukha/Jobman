package main

import (
	"flag"
	"log/slog"
	"os"
)

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *slog.Logger
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.Parse()

	app := &application{
		config: cfg,
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	err := app.serve()
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
}
