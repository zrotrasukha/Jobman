package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// routes sets up the router and a http.Handler.
func (app *application) routes() http.Handler {
	mux := httprouter.New()
	mux.NotFound = http.HandlerFunc(app.notFoundResponse)

	activatedUserChain := alice.New(app.authenticate, func(next http.Handler) http.Handler {
		return app.requireActivatedUser(next.ServeHTTP)
	})

	ownerRequiredChain := activatedUserChain.Append(app.requireApplicationOwner)

	mux.Handler(http.MethodGet, "/v1/healthcheck", activatedUserChain.ThenFunc(app.healthcheckHandler))
	mux.Handler(http.MethodGet, "/v1/applications", activatedUserChain.ThenFunc(app.ListApplicationHandler))
	mux.Handler(http.MethodPost, "/v1/applications", activatedUserChain.ThenFunc(app.CreateApplicationHandler))

	mux.Handler(http.MethodGet, "/v1/applications/:id", ownerRequiredChain.ThenFunc(app.GetApplicationHandler))
	mux.Handler(http.MethodPatch, "/v1/applications/:id", ownerRequiredChain.ThenFunc(app.UpdateApplicationHandler))
	mux.Handler(http.MethodDelete, "/v1/applications/:id", ownerRequiredChain.ThenFunc(app.DeleteApplicationHandler))

	mux.Handler(http.MethodGet, "/v1/reminders", activatedUserChain.ThenFunc(app.ListRemindersHandler))

	mux.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	mux.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	mux.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.CreateAuthenticationTokenHandler)

	return app.recoverPanic(app.reqLogger(mux))
}
