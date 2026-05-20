package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/zrotrasukha/jobman/internal/assert"
	"github.com/zrotrasukha/jobman/internal/data"
	"github.com/zrotrasukha/jobman/internal/data/mocks"
)

func TestHealthcheck(t *testing.T) {
	ts := newTestServer(t, mocks.MockJobApplicationModel{})
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
		mockModels mocks.MockJobApplicationModel
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
				"applied_at" : "2026-08-12T11:45:00Z",
				"notes": "Test Notes"
			}`,
			mockModels: mocks.MockJobApplicationModel{},
			wantHeader: true,
			wantStatus: http.StatusCreated,
			wantBody: `{
        	"application": {
        		"id": 1,
        		"company_name": "Test Company",
        		"role_title": "Test Role",
        		"status": "Applied",
						"applied_at" : "2026-08-12T11:45:00Z",
						"updated_at" : "2026-08-12T11:45:00Z",
        		"last_communication": null,
        		"notes": "Test Notes",
        		"version": 1
        	}
        }`,
		},
		{
			name: "invalid field",
			input: `{
				"company_name": "Test Company",
				"role_title": "Test Role",
				"applied_at" : "2026-08-12T11:45:00Z",
				"status": "Applied",
				"notes": "Test Notes",
				"invalid_field": "Invalid Value"
			}`,
			mockModels: mocks.MockJobApplicationModel{},
			wantHeader: false,
			wantStatus: http.StatusBadRequest,

			wantBody: `{"error": "unknown field \"invalid_field\""}`,
		},
		{
			name: "invalid status",
			input: `{
				"company_name": "Test Company",
				"role_title": "Test Role",
				"applied_at" : "2026-08-12T11:45:00Z",
				"status": "Invalid Status",
				"notes": "Test Notes"
			}`,
			mockModels: mocks.MockJobApplicationModel{},
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
			mockModels: mocks.MockJobApplicationModel{},
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
			mockModels: mocks.MockJobApplicationModel{},
			wantHeader: false,
			wantStatus: http.StatusUnprocessableEntity,
			wantBody: `{
				"error": {
								"applied_at": "must be a valid RFC3339 date"
				}
			}`,
		},
		{
			name: "db fails",
			input: `{
				"company_name": "Test Company",
				"role_title": "Test Role",
				"applied_at" : "2026-08-12T11:45:00Z",
				"status": "Applied",
				"notes": "Test Notes"	
			}`,
			mockModels: mocks.MockJobApplicationModel{
				InsertFunc: func(jobApp *data.JobApplication) error {
					return errors.New("db fails to load")
				},
			},
			wantHeader: false,
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"error": "the server encountered a problem and could not process your request"}`,
		},
	}

	for _, tt := range tests {
		ts := newTestServer(t, tt.mockModels)

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

func TestGetApplicationHandler(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		wantCode   int
		mockModels mocks.MockJobApplicationModel
		wantBody   string
	}{
		{
			name:       "valid id",
			id:         "1",
			wantCode:   http.StatusOK,
			mockModels: mocks.MockJobApplicationModel{},
			wantBody: `{
                "application": {
                        "id": 1,
                        "company_name": "Test Company",
                        "role_title": "Test Role",
                        "status": "Applied",
                        "applied_at": "2026-08-12T11:45:00Z",
                        "updated_at": "2026-08-12T11:45:00Z",
                        "last_communication": null,
                        "notes": "Test Notes",
                        "version": 1
                }
        }`,
		},
		{
			name:       "invalid id",
			id:         "abc",
			wantCode:   http.StatusBadRequest,
			mockModels: mocks.MockJobApplicationModel{},
			wantBody:   `{"error": "invalid id parameter"}`,
		},
		{
			name:     "non-existent id",
			id:       "999",
			wantCode: http.StatusNotFound,
			mockModels: mocks.MockJobApplicationModel{
				GETFunc: func(id int64) (*data.JobApplication, error) {
					return nil, data.ErrRecordNotFound
				},
			},
			wantBody: `{"error": "the requested resource could not be found"}`,
		},
		{
			name:     "db fails",
			id:       "1",
			wantCode: http.StatusInternalServerError,
			mockModels: mocks.MockJobApplicationModel{
				GETFunc: func(id int64) (*data.JobApplication, error) {
					return nil, errors.New("db fails to load")
				},
			},
			wantBody: `{"error": "the server encountered a problem and could not process your request"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := newTestServer(t, tt.mockModels)
			urlPath := "/v1/applications/" + tt.id
			sc, body := ts.Get(t, urlPath)

			assert.Equal(t, sc, tt.wantCode)
			assert.EqualJSON(t, []byte(body), tt.wantBody)
		})
	}
}
