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
	ts := newTestServer(t, data.Models{})
	app := newTestApplication()

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
		name                  string
		input                 string
		mockApplicationsModel mocks.MockJobApplicationModel
		wantHeader            bool
		wantStatus            int
		wantBody              string
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
			mockApplicationsModel: mocks.MockJobApplicationModel{
				InsertFunc: func(jobApp *data.JobApplication) error {
					jobApp.ID = 1
					jobApp.Version = 1
					return nil
				},
			},
			wantHeader: true,
			wantStatus: http.StatusCreated,
			wantBody: `{
        	"application": {
        		"id": 1,
        		"company_name": "Test Company",
        		"role_title": "Test Role",
        		"status": "Applied",
        		"applied_at": "2026-08-12T11:45:00Z",
        		"updated_at": "0001-01-01T00:00:00Z",
        		"interview_at": null,
        		"stale_after": "2026-09-11T11:45:00Z",
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
			mockApplicationsModel: mocks.MockJobApplicationModel{},
			wantHeader:            false,
			wantStatus:            http.StatusBadRequest,
			wantBody:              `{"error": "unknown field \"invalid_field\""}`,
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
			mockApplicationsModel: mocks.MockJobApplicationModel{},
			wantHeader:            false,
			wantStatus:            http.StatusBadRequest,
			wantBody:              `{"error": "invalid status: Invalid Status"}`,
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
			mockApplicationsModel: mocks.MockJobApplicationModel{},
			wantHeader:            false,
			wantStatus:            http.StatusUnprocessableEntity,
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
			mockApplicationsModel: mocks.MockJobApplicationModel{},
			wantHeader:            false,
			wantStatus:            http.StatusUnprocessableEntity,
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
			mockApplicationsModel: mocks.MockJobApplicationModel{
				InsertFunc: func(jobApp *data.JobApplication) error {
					return errors.New("db fails to load")
				},
			},
			wantHeader: false,
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"error": "the server encountered a problem and could not process your request"}`,
		},
		{
			name: "valid input with interview_at calculates 5-day staleness",
			input: `{
                "company_name": "Google",
                "role_title": "Backend Engineer",
                "status": "Interviewing",
                "applied_at": "2026-08-12T11:45:00Z",
                "interview_at": "2026-08-15T10:00:00Z"
            }`,
			mockApplicationsModel: mocks.MockJobApplicationModel{
				InsertFunc: func(jobApp *data.JobApplication) error {
					jobApp.ID = 1
					jobApp.Version = 1
					return nil
				},
			},
			wantHeader: true,
			wantStatus: http.StatusCreated, // Verifies that stale_after is exactly Aug 15 + 5 days = Aug 20
			wantBody: `{
                "application": {
                    "id": 1,
                    "company_name": "Google",
                    "role_title": "Backend Engineer",
                    "status": "Interviewing",
                    "applied_at": "2026-08-12T11:45:00Z",
                    "interview_at": "2026-08-15T10:00:00Z",
                    "stale_after": "2026-08-20T10:00:00Z",
                    "updated_at": "0001-01-01T00:00:00Z",
                    "last_communication": null,
                    "notes": "",
                    "version": 1
                }
            }`,
		},
		{
			name: "invalid interview_at format",
			input: `{
                "company_name": "Google",
                "role_title": "Backend Engineer",
                "status": "Interviewing",
                "applied_at": "2026-08-12T11:45:00Z",
                "interview_at": "bad-date-format"
            }`,
			mockApplicationsModel: mocks.MockJobApplicationModel{},
			wantHeader:            false,
			wantStatus:            http.StatusUnprocessableEntity,
			wantBody: `{
                "error": {
                    "interview_at": "must be a valid RFC3339 date"
                }
            }`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			models := mocks.NewMockModels()
			models.Application = tt.mockApplicationsModel

			ts := newTestServer(t, models)
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
		name                 string
		id                   string
		wantCode             int
		mockApplicationModel mocks.MockJobApplicationModel
		wantBody             string
	}{
		{
			name:                 "invalid id",
			id:                   "abc",
			wantCode:             http.StatusBadRequest,
			mockApplicationModel: mocks.MockJobApplicationModel{},
			wantBody:             `{"error": "invalid id parameter"}`,
		},
		{
			name:     "non-existent id",
			id:       "999",
			wantCode: http.StatusNotFound,
			mockApplicationModel: mocks.MockJobApplicationModel{
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
			mockApplicationModel: mocks.MockJobApplicationModel{
				GETFunc: func(id int64) (*data.JobApplication, error) {
					return nil, errors.New("db fails to load")
				},
			},
			wantBody: `{"error": "the server encountered a problem and could not process your request"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			models := mocks.NewMockModels()
			models.Application = tt.mockApplicationModel

			ts := newTestServer(t, models)
			urlPath := "/v1/applications/" + tt.id
			sc, body := ts.Get(t, urlPath)

			assert.Equal(t, sc, tt.wantCode)
			assert.EqualJSON(t, []byte(body), []byte(tt.wantBody))
		})
	}
}

func TestListApplicationHandler(t *testing.T) {
	tests := []struct {
		name                 string
		url                  string
		wantCode             int
		mockApplicationModel mocks.MockJobApplicationModel
		wantBody             string
	}{
		{
			name:     "invalid parameter - page too low",
			url:      "/v1/applications?page=0",
			wantCode: http.StatusUnprocessableEntity,
			mockApplicationModel: mocks.MockJobApplicationModel{
				GetAllFunc: func(searchString string, filters data.Filters) ([]*data.JobApplication, *data.Metadata, error) {
					t.Fatal("GetAll should not be called if validation fails!")
					return nil, nil, nil
				},
			},
			wantBody: `{"error": {"page": "must be greater than zero"}}`,
		},
		{
			name:     "invalid parameter - page size too high",
			url:      "/v1/applications?page_size=999",
			wantCode: http.StatusUnprocessableEntity,
			mockApplicationModel: mocks.MockJobApplicationModel{
				GetAllFunc: func(searchString string, filters data.Filters) ([]*data.JobApplication, *data.Metadata, error) {
					t.Fatal("GetAll should not be called if validation fails!")
					return nil, nil, nil
				},
			},
			wantBody: `{"error": {"page_size": "must be a maximum of 100"}}`,
		},
		{
			name:     "invalid parameter - dangerous sort clause",
			url:      "/v1/applications?sort=drop_table_users",
			wantCode: http.StatusUnprocessableEntity,
			mockApplicationModel: mocks.MockJobApplicationModel{
				GetAllFunc: func(searchString string, filters data.Filters) ([]*data.JobApplication, *data.Metadata, error) {
					t.Fatal("GetAll should not be called if validation fails!")
					return nil, nil, nil
				},
			},
			wantBody: `{"error": {"sort": "invalid sort value"}}`,
		},
		{
			name:     "db fails",
			url:      "/v1/applications",
			wantCode: http.StatusInternalServerError,
			mockApplicationModel: mocks.MockJobApplicationModel{
				GetAllFunc: func(searchString string, filters data.Filters) ([]*data.JobApplication, *data.Metadata, error) {
					return nil, nil, errors.New("database connection completely fried")
				},
			},
			wantBody: `{"error": "the server encountered a problem and could not process your request"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			models := mocks.NewMockModels()
			models.Application = tt.mockApplicationModel

			ts := newTestServer(t, models)
			sc, body := ts.Get(t, tt.url)

			assert.Equal(t, sc, tt.wantCode)
			assert.EqualJSON(t, []byte(body), []byte(tt.wantBody))
		})
	}
}
