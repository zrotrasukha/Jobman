package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/zrotrasukha/jobman/internal/assert"
	"github.com/zrotrasukha/jobman/internal/data/mocks"
)

func TestWriteJSON(t *testing.T) {
	app := newTestApplication()
	env := envelop{
		"message": "testing TestReadJSON helper",
	}

	w := httptest.NewRecorder()
	header := make(http.Header)
	header.Set("X-Test-Header", "test value")

	err := app.writeJSON(w, http.StatusOK, env, header)
	if err != nil {
		t.Fatal(err)
	}

	rs := w.Result()
	defer rs.Body.Close()

	//header check

	assert.Equal(t, rs.StatusCode, http.StatusOK)
	assert.Equal(t, rs.Header.Get("Content-Type"), "application/json")
	assert.Equal(t, rs.Header.Get("X-Test-Header"), "test value")

	b, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualJSON(t, b, env)
}

func TestReadJSON(t *testing.T) {
	type Input struct {
		Name string `json:"name"`
	}
	tests := []struct {
		name string
		body string
		err  string
	}{
		{
			name: "valid body",
			body: `{"name": "test"}`,
			err:  "",
		},
		{
			name: "empty body",
			body: ``,
			err:  "body must not be empty",
		},
		{
			name: "malformed body",
			body: `{"name": "test"`,
			err:  "malformed JSON",
		},
		{
			name: "invalid value for field",
			body: `{"name": 123}`,
			err:  "invalid value for field \"name\" (at character 12)",
		},
		{
			name: "unknown field",
			body: `{"name": "test", "age": 30}`,
			err:  `unknown field "age"`,
		},
		{
			name: "multiple JSON values in body",
			body: `{"name": "test"}{"name": "test2"}`,
			err:  "body must only contain a single JSON value",
		},
		{
			name: "large body",
			body: `{"name": "` + strings.Repeat("a", 1_048_577) + `"}`,
			err:  "body must not be larger than 1048576 bytes",
		},
	}
	app := newTestApplication()
	ts := newTestServer(t, mocks.NewMockModels())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			var dst Input

			req := httptest.NewRequest(http.MethodPost, ts.URL, strings.NewReader(tt.body))
			err := app.readJSON(w, req, &dst)

			if tt.err == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			} else {
				if err == nil {
					t.Fatalf("got nil error; want: %v", tt.err)
				} else {
					if !strings.Contains(err.Error(), tt.err) {
						t.Fatalf("got error: %v; want error: %v", err, tt.err)
					}
				}
			} //t.run end
		}) //iteration end
	}
}

func TestReadParamID(t *testing.T) {
	test := []struct {
		name    string
		paramID string
		wantID  int64
		wantErr bool
	}{
		{
			name:    "valid ID",
			paramID: "67",
			wantID:  67,
			wantErr: false,
		},
		{
			name:    "invalid ID",
			paramID: "playouteerwildsmanshitissickasfuck",
			wantID:  0,
			wantErr: true,
		},
		{
			name:    "negative ID",
			paramID: "-5",
			wantID:  0,
			wantErr: true,
		},
		{
			name:    "zero ID",
			paramID: "0",
			wantID:  0,
			wantErr: true,
		},
	}

	app := newTestApplication()
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			params := httprouter.Params{
				httprouter.Param{Key: "id", Value: tt.paramID},
			}

			ctx := context.WithValue(req.Context(), httprouter.ParamsKey, params)
			req = req.WithContext(ctx)

			gotID, err := app.readParamID(req)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("got nil error; want error: %t", tt.wantErr)
				} else {
					if !strings.Contains(err.Error(), "invalid id parameter") {
						t.Fatalf("got error: %v; want error: %v", err, "invalid id parameter")
					}
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				} else {
					assert.Equal(t, gotID, tt.wantID)
				}
			}

		})
	}

}

func TestCurrentWeekWindow(t *testing.T) {
	app := newTestApplication()

	tests := []struct {
		name     string
		input    time.Time
		wantFrom time.Time
		wantTo   time.Time
	}{
		{
			name:     "Tuesday in middle of week",
			input:    time.Date(2026, 7, 7, 10, 30, 0, 0, time.UTC), // Tuesday
			wantFrom: time.Date(2026, 7, 6, 0, 0, 0, 0, time.UTC),  // Monday
			wantTo:   time.Date(2026, 7, 12, 23, 59, 59, 999999999, time.UTC), // Sunday
		},
		{
			name:     "Monday (start of week)",
			input:    time.Date(2026, 7, 6, 8, 0, 0, 0, time.UTC),
			wantFrom: time.Date(2026, 7, 6, 0, 0, 0, 0, time.UTC),
			wantTo:   time.Date(2026, 7, 12, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name:     "Sunday (end of week)",
			input:    time.Date(2026, 7, 12, 23, 0, 0, 0, time.UTC),
			wantFrom: time.Date(2026, 7, 6, 0, 0, 0, 0, time.UTC),
			wantTo:   time.Date(2026, 7, 12, 23, 59, 59, 999999999, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFrom, gotTo := app.currentWeekWindow(tt.input)
			if !gotFrom.Equal(tt.wantFrom) {
				t.Errorf("gotFrom = %v; want %v", gotFrom, tt.wantFrom)
			}
			if !gotTo.Equal(tt.wantTo) {
				t.Errorf("gotTo = %v; want %v", gotTo, tt.wantTo)
			}
		})
	}
}
