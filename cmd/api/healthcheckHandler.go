package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	message := envelop{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	err := WriteJSON(w, http.StatusOK, message, nil)
	if err != nil {
		app.serverErrResponse(w, r)
	}
}
