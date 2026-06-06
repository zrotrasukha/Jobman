package main

import (
	"net/http"
)

// healthcheckHandler is an HTTP handler function that responds to health check requests. It returns 200 OK status OK if the application is healthy, along with a JSON response containing the application's status, environment, and version. If there is an error while writing the JSON response, it sends a server error response to the client.
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	message := envelop{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	err := app.writeJSON(w, http.StatusOK, message, nil)
	if err != nil {
		app.serverErrResponse(w, r, err)
	}
}
