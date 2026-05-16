package main

import "net/http"

func (app *application) reqLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.logger.Info("request started", "method", r.Method, "url", r.URL.String())
		next.ServeHTTP(w, r)
	})
}
