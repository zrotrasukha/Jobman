package main

import (
	"context"
	"time"
)

func (app *application) stalenessWorker(ctx context.Context) error {
	ticker := time.NewTicker(6 * time.Hour)
	defer ticker.Stop()

	for {
		rowsAffected, err := app.models.Application.MarkStaleApplications(ctx)
		if err != nil {
			app.logger.Error("staleness worker: failed to mark stale applications", "error", err)
		} else if rowsAffected > 0 {
			app.logger.Info("staleness worker: marked stale applications as ghosted", "rows_affected", rowsAffected)
		}

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			app.logger.Info("staleness worker: shutting down")
			return nil
		}
	}
}
