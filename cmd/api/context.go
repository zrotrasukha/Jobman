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
