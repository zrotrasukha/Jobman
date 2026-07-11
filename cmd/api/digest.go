package main

import (
	"net/http"
	"time"

	"github.com/zrotrasukha/jobman/internal/data"
	"github.com/zrotrasukha/jobman/internal/validator"
)

// ListDigestHandler handles the GET /digest endpoint and returns a JSON response with the digest metrics for the authenticated user.
func (app *application) ListDigestHandler(w http.ResponseWriter, r *http.Request) {
	user := app.ContextGetUser(r)

	qs := r.URL.Query()
	window := app.readString(qs, "window", "7d")

	v := validator.New()
	data.ValidateDigest(v, window)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	now := time.Now()
	var from, to time.Time

	switch window {
	case "7d":
		from, to = app.currentWeekWindow(now)
	case "1m":
		from, to = app.currentMonthWindow(now)
	case "1y":
		from, to = app.currentYearWindow(now)
	}

	digest, err := app.models.Digest.GetDigest(user.Id, from, to)
	if err != nil {
		app.serverErrResponse(w, r, err)
		return
	}

	response := envelop{
		"cached": false,
		"digest": digest,
	}

	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrResponse(w, r, err)
	}
}
