package main

import (
	"context"
)

func (app *application) stalenessWorker(ctx context.Context) error {
	return app.runPeriodic(ctx, "staleness", app.config.stalenessInterval, func() error {
		rowsAffected, err := app.models.Application.MarkStaleApplications(ctx)
		if err != nil {
			return err
		}

		if rowsAffected > 0 {
			app.logger.Info("staleness worker: marked stale applications as ghosted", "rows_affected", rowsAffected)
		}
		return nil
	})
}
