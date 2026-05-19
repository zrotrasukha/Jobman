package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/zrotrasukha/jobman/internal/assert"
)

func TestHealthcheck(t *testing.T) {
	ts := newTestServer(t)
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

func TestCreateApplicationHandler(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		ts := newTestServer(t)

		input := `{
			"company_name": "Test Company",
			"role_title": "Test Role",
			"status": "Applied",
			"notes": "Test Notes"
		}`

		status, header, body := ts.Post(t, "/v1/applications", input)
		body = strings.TrimSpace(body)

		assert.Equal(t, status, http.StatusCreated)
		assert.Equal(t, header.Get("Location"), "/v1/applications/1")

		var got envelop
		err := json.Unmarshal([]byte(body), &got)
		if err != nil {
			t.Fatal(err)
		}
		assert.EqualJSON(t, []byte(body), &got)
	})
}
