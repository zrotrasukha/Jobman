package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	mux := httprouter.New()

	mux.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	mux.HandlerFunc(http.MethodPost, "/v1/applications", app.CreateApplicationHandler)
	mux.HandlerFunc(http.MethodGet, "/v1/applications", app.GetApplicationsHandler)
	mux.HandlerFunc(http.MethodGet, "/v1/applications/:id", app.GetApplicationHandler)
	mux.HandlerFunc(http.MethodPut, "/v1/applications/:id", app.UpdateApplicationHandler)
	mux.HandlerFunc(http.MethodDelete, "/v1/applications/:id", app.DeleteApplicationHandler)

	return mux

}
