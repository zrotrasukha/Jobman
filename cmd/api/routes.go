package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// routes sets up the router and a http.Handler.
func (app *application) routes() http.Handler {
	mux := httprouter.New()
	mux.NotFound = http.HandlerFunc(app.notFoundResponse)

	mux.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	mux.HandlerFunc(http.MethodPost, "/v1/applications", app.CreateApplicationHandler)
	mux.HandlerFunc(http.MethodGet, "/v1/applications", app.ListApplicationHandler)
	mux.HandlerFunc(http.MethodGet, "/v1/applications/:id", app.GetApplicationHandler)
	mux.HandlerFunc(http.MethodPatch, "/v1/applications/:id", app.UpdateApplicationHandler)
	mux.HandlerFunc(http.MethodDelete, "/v1/applications/:id", app.DeleteApplicationHandler)

	mux.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	mux.HandlerFunc(http.MethodPut, "/v1/tokens/authentication", app.activateUserHandler)

	return app.recoverPanic(app.reqLogger(mux))
}
