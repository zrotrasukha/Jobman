package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zrotrasukha/jobman/internal/data"
	"github.com/zrotrasukha/jobman/internal/validator"
)

// reqLogger middleware is use to log the incoming HTTP requests. It logs the HTTP method and the URL of each request when it starts processing.
func (app *application) reqLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.logger.Info("request started", "method", r.Method, "url", r.URL.String())
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrResponse(w, r, fmt.Errorf("panic: %v", err))
				return
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = app.ContextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headParts := strings.Split(authorizationHeader, " ")
		if len(headParts) != 2 || headParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headParts[1]

		v := validator.New()
		if data.ValidateTokenPlainText(v, token); !v.Valid() {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		user, err := app.models.User.GetForToken(token, data.ScopeAuthentication)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrResponse(w, r, err)
			}
			return
		}

		r = app.ContextSetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthenticatedUser(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.ContextGetUser(r)

		if user.IsAnonymous() {
			app.authenticationRequired(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireActivatedUser(next http.HandlerFunc) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user := app.ContextGetUser(r)

		if !user.Activated {
			app.inactiveAccountResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	}

	return app.requireAuthenticatedUser(http.HandlerFunc(fn))
}

func (app *application) requireApplicationOwner(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.ContextGetUser(r)
		id, err := app.readParamID(r)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		jobApp, err := app.models.Application.Get(id, user.Id)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrResponse(w, r, err)
			}
			return
		}

		r = app.ContextSetApplication(r, jobApp)

		next.ServeHTTP(w, r)
	})

}
