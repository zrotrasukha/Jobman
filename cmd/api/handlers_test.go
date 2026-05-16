package main

import (
	"net/http"
	"testing"

	"github.com/zrotrasukha/jobman/internal/assert"
)

func TestHealth(T *testing.T) {
	ts := newTestServer()

	status, body := ts.Get(T, "/v1/healthcheck")

	assert.Equal(T, status, http.StatusOK)
	assert.Equal(T, body, "Status: OK,\nEnvironment: test")
}
