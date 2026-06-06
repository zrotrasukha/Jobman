package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zrotrasukha/jobman/internal/data"
	"github.com/zrotrasukha/jobman/internal/data/mocks"
)

func newTestApplication() *application {
	cfg := config{
		env: "test",
	}
	return &application{
		config: cfg,
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		models: mocks.NewMockModels(),
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, model data.Models) *testServer {
	app := newTestApplication()
	app.models = model

	ts := httptest.NewServer(app.routes())
	t.Cleanup(ts.Close)

	return &testServer{ts}
}

// ts.Get() will send a GET request to the test server with the given url path, and return the status code and response body.
func (ts *testServer) Get(t *testing.T, urlPath string) (int, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	body = bytes.TrimSpace(body)
	return rs.StatusCode, string(body)
}

// ts.Post() will seend a POST request to the test server with the given url path and body, and return the status code, headers, and response body.
func (ts *testServer) Post(t *testing.T, urlPath string, body string) (int, http.Header, string) {
	url := ts.URL + urlPath
	buf := bytes.NewBufferString(body)
	rs, err := ts.Client().Post(url, "application/json", buf)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	respBody, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	respBody = bytes.TrimSpace(respBody)
	// return status, headers, body
	return rs.StatusCode, rs.Header, string(respBody)
}
