package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

// serve starts the HTTP server and listens for incoming requests. It also sets up a goroutine to handle graceful shutdown when an interrupt signal is received. The verver will wait for any ongoing requests to complete before shutting down, and it will log the shutdown process. If any errors occur during startup or shutdown, they are returned to the caller.
func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		sig := <-sigChan
		app.logger.Info("Interruption occurred", "signal", sig.String())
		cancel()
	}()

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		app.logger.Info("Starting server", "port", app.config.port, "env", app.config.env)
		err := srv.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	})

	g.Go(func() error {
		return app.stalenessWorker(gCtx)
	})

	g.Go(func() error {
		<-gCtx.Done()
		app.logger.Info("shutting down server")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			return err
		}

		app.logger.Info("closing cache connection")
		return app.cache.Close()
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
