package main

import (
	"fmt"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("Status: OK,\nEnvironment: %s", app.config.env)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}
