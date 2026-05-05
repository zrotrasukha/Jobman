package main

import (
	"fmt"
	"net/http"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
	}

	app.logger.Info("Starting server", "port", app.config.port, "env", app.config.env)

	err := srv.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
