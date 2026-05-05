package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http/httptest"
	"testing"
)

func newTestApplication() *application {

	cfg := config{
		env: "test",
	}
	return &application{
		config: cfg,
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer() *testServer {
	app := newTestApplication()
	ts := httptest.NewServer(app.routes())
	return &testServer{ts}
}

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
