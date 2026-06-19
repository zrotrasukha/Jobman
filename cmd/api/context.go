package main

import (
	"context"
	"net/http"

	"github.com/zrotrasukha/jobman/internal/data"
)

type contextKey string

const userContextKey = contextKey("user")

func (app *application) ContextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) ContextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}

func (app *application) ContextSetApplication(r *http.Request, application *data.JobApplication) *http.Request {
	ctx := context.WithValue(r.Context(), "application", application)
	return r.WithContext(ctx)
}

func (app *application) ContextGetApplication(r *http.Request) *data.JobApplication {
	application, ok := r.Context().Value("application").(*data.JobApplication)
	if !ok {
		panic("missing application value in request context")
	}
	return application
}
