package main

import (
	"net/http"
	"testing"

	"github.com/zrotrasukha/jobman/internal/assert"
	"github.com/zrotrasukha/jobman/internal/data"
	"github.com/zrotrasukha/jobman/internal/data/mocks"
)

func TestRegisterUserHandler(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		url           string
		wantStatus    int
		mockUserModel mocks.MockUserModel
		wantBody      string
	}{
		{
			name:       "valid request",
			wantStatus: http.StatusCreated,
			input: `{
				"name": "Test User",
				"email": "testuser@mail.com",
				"password": "password123"
			}	`,
			url:           "/v1/users",
			mockUserModel: mocks.MockUserModel{},
			wantBody: `{
				"user": {
					"id": 1,
					"created_at": "2026-08-12T11:45:00Z",
					"name": "Test User",
					"email": "test@user.com",
					"activated": false
				}
			}`,
		},
		{
			name: "invalid email",
			input: `
				"name":     "test user",
				"email":    "bad@mail",
				"password": "just!secret?",
			}`,
			url:        "/v1/users",
			wantStatus: http.StatusBadRequest,
			wantBody:   `{ "error": "invalid JSON (at character 11)" }`,
		},
		{
			name: "password",
			input: `{ 
				"name": "test user",
				"email": "test@user.com"
				}`,
			url:           "/v1/users",
			wantStatus:    http.StatusUnprocessableEntity,
			mockUserModel: mocks.MockUserModel{},
			wantBody: `{
				"error": {
					"password": "must be provided"
				}
			}`,
		},
		{
			name: "no email",
			input: `{
				"name": "test user",
				"password": "blahblahpassword"
			}`,
			url:           "/v1/users",
			wantStatus:    http.StatusUnprocessableEntity,
			mockUserModel: mocks.MockUserModel{},
			wantBody: `{
        	"error": {
        		"email": "must be provided"
        	}
        }`,
		},
		{
			name: "short password length",
			input: `{
				"name": "test user",
				"email": "test@gmail.com",
				"password": "short"
			}`,
			url:           "/v1/users",
			wantStatus:    http.StatusUnprocessableEntity,
			mockUserModel: mocks.MockUserModel{},
			wantBody: `{
				"error": {
					"password": "must be at least 8 bytes long"
				}
			}`,
		},
		{
			name: "long password length",
			input: `{
				"name": "long user",
				"email": "long@gmail.com",
				"password": "blahhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhnhhhnn"
			}`,
			url:           "/v1/users",
			wantStatus:    http.StatusUnprocessableEntity,
			mockUserModel: mocks.MockUserModel{},
			wantBody: `{
				"error": {
					"password": "must not be more than 72 bytes long"
				}
			}`,
		},
		{
			name: "duplicate email",
			input: `{
				"name": "duplicate user",
				"email": "duplicate@gmail.com",
				"password": "pa55word"
			}`,
			url:        "/v1/users",
			wantStatus: http.StatusUnprocessableEntity,
			mockUserModel: mocks.MockUserModel{
				InsertFunc: func(user *data.User) error {
					return data.ErrDuplicateEmail
				},
			},
			wantBody: `{
				"error": {
					"email": "a user with this email address already exists"
				}
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			models := mocks.NewMockModels()
			models.User = tt.mockUserModel

			ts := newTestServer(t, models)
			sc, _, body := ts.Post(t, tt.url, tt.input)

			assert.Equal(t, sc, tt.wantStatus)
			assert.EqualJSON(t, []byte(body), []byte(tt.wantBody))
		})
	}
}
