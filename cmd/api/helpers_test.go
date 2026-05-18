package main

import (
	"io"
	"net/http"
	"net/http/httptest"
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

//TODO: Write test for readJSON
//TODO: Write test for readParamID
