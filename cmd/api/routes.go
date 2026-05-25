package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// routes sets up the HTTP routes for the application. It returns a mux which is being consumed in serve.go. The mux is configured to handle various HTTP methods and paths, and it also includes a custom NotFound handler for any routes that are not defined.
func (app *application) routes() http.Handler {
	mux := httprouter.New()
	mux.NotFound = http.HandlerFunc(app.notFoundResponse)

	mux.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	mux.HandlerFunc(http.MethodPost, "/v1/applications", app.CreateApplicationHandler)
	mux.HandlerFunc(http.MethodGet, "/v1/applications", app.ListApplicationHandler)
	mux.HandlerFunc(http.MethodGet, "/v1/applications/:id", app.GetApplicationHandler)
	mux.HandlerFunc(http.MethodPut, "/v1/applications/:id", app.UpdateApplicationHandler)
	mux.HandlerFunc(http.MethodDelete, "/v1/applications/:id", app.DeleteApplicationHandler)

	return app.reqLogger(mux)
}
