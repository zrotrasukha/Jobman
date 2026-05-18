package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/zrotrasukha/jobman/internal/assert"
)

func TestHealthcheck(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	app := newTestApplication(t)

	expected := envelop{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	status, body := ts.Get(t, "/v1/healthcheck")
	body = strings.TrimSpace(body)

	assert.Equal(t, status, http.StatusOK)
	assert.EqualJSON(t, []byte(body), expected)
}
