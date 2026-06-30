package main

import (
	"net/http"

	"github.com/zrotrasukha/jobman/internal/data"
)

// ListRemindersHandler returns upcoming interview reminders for the authenticated user
func (app *application) ListRemindersHandler(w http.ResponseWriter, r *http.Request) {
	user := app.ContextGetUser(r)

	reminders, err := app.models.Reminder.GetUpcoming(user.Id, 10)
	if err != nil {
		app.serverErrResponse(w, r, err)
		return
	}

	// guardrail for nil response, as cli always expects non-nil slice to range over it
	if reminders == nil {
		reminders = []*data.Reminder{}
	}

	err = app.writeJSON(w, http.StatusOK, envelop{"reminders": reminders}, nil)
	if err != nil {
		app.serverErrResponse(w, r, err)
	}
}
