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

	js, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, status, http.StatusOK)
	assert.EqualJSON(t, []byte(body), js)
}

func TestCreateApplicationHandler(t *testing.T) {

	tests := []struct {
		name       string
		input      string
		wantHeader bool
		wantStatus int
		wantBody   string
	}{
		{
			name: "valid input",
			input: `{
				"company_name": "Test Company",
				"role_title": "Test Role",
				"status": "Applied",
				"applied_at" : "2024-06-01T12:00:00Z",
				"notes": "Test Notes"
			}`,
			wantHeader: true,
			wantStatus: http.StatusCreated,
			wantBody: `{
        	"application": {
        		"id": 1,
        		"company_name": "Test Company",
        		"role_title": "Test Role",
        		"status": "Applied",
        		"applied_at": "2024-06-01T12:00:00Z",
        		"last_communication": null,
        		"notes": "Test Notes",
        		"version": 0
        	}
        }`,
		},
		{
			name: "invalid field",
			input: `{
				"company_name": "Test Company",
				"role_title": "Test Role",
				"applied_at" : "2024-06-01T12:00:00Z",
				"status": "Applied",
				"notes": "Test Notes",
				"invalid_field": "Invalid Value"
			}`,
			wantHeader: false,
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error": "unknown field \"invalid_field\""}`,
		},
		{
			name: "invalid status",
			input: `{
				"company_name": "Test Company",
				"role_title": "Test Role",
				"applied_at" : "2024-06-01T12:00:00Z",
				"status": "Invalid Status",
				"notes": "Test Notes"
			}`,
			wantHeader: false,
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error": "invalid status: Invalid Status"}`,
		},
		{
			name: "missing required field",
			input: `{
				"company_name": "",
				"role_title": "",
				"applied_at": "",
				"status": "Applied",
				"notes": "Applied through referral"
			}`,
			wantHeader: false,
			wantStatus: http.StatusUnprocessableEntity,
			wantBody: `{
		      "error": {
		              "applied_at": "must be provided",
		              "company_name": "must be provided",
		              "role_title": "must be provided"
		      }
			}`,
		},
		{
			name: "invalid applied_at format",
			input: `{
				"company_name": "Test Company",
				"role_title": "Test Role",
				"applied_at" : "invalid-date-format",
				"status": "Applied",
				"notes": "Test Notes"
			}`,
			wantHeader: false,
			wantStatus: http.StatusUnprocessableEntity,
			wantBody: `{
				"error": {
								"applied_at": "must be a valid RFC3339 date"
				}
			}`,
		},
	}

	ts := newTestServer(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, header, body := ts.Post(t, "/v1/applications", tt.input)
			body = strings.TrimSpace(body)

			assert.Equal(t, status, tt.wantStatus)
			if tt.wantHeader {
				assert.Equal(t, header.Get("Location"), "/v1/applications/1")
			}

			assert.EqualJSON(t, []byte(body), []byte(tt.wantBody))
		})
	}
}
