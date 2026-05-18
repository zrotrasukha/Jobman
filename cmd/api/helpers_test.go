package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/zrotrasukha/jobman/internal/assert"
)

func TestWriteJSON(t *testing.T) {
	app := newTestApplication(t)
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
	app := newTestApplication(t)
	ts := newTestServer(t)

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

//
//TODO: Write test for readParamID
