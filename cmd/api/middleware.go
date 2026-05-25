package main

import "net/http"

// reqLogger middleware is use to log the incoming HTTP requests. It logs the HTTP method and the URL of each request when it starts processing.
func (app *application) reqLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.logger.Info("request started", "method", r.Method, "url", r.URL.String())
		next.ServeHTTP(w, r)
	})
}
