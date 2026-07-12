package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/zrotrasukha/jobman/internal/assert"
	"github.com/zrotrasukha/jobman/internal/data"
	"github.com/zrotrasukha/jobman/internal/data/mocks"
)

func TestListDigestHandler(t *testing.T) {
	ts := newTestServer(t, mocks.NewMockModels())

	t.Run("valid default request", func(t *testing.T) {
		status, body := ts.Get(t, "/v1/digest")
		assert.Equal(t, status, http.StatusOK)

		var resp struct {
			Cached bool `json:"cached"`
			Digest struct {
				Window      string    `json:"window"`
				From        string    `json:"from"`
				To          string    `json:"to"`
				Funnel      struct {
					Applied      int `json:"applied"`
					Replied      int `json:"replied"`
					Interviewing int `json:"interviewing"`
					Offered      int `json:"offered"`
					Rejected     int `json:"rejected"`
					Ghosted      int `json:"ghosted"`
				} `json:"funnel"`
				GhostCohort struct {
					Matured        int       `json:"matured"`
					MaturedGhosted int       `json:"matured_ghosted"`
					Pending        int       `json:"pending"`
					Rate           data.Rate `json:"rate"`
				} `json:"ghost_cohort"`
			} `json:"digest"`
		}

		err := json.Unmarshal([]byte(body), &resp)
		if err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		assert.Equal(t, resp.Cached, false)
		assert.Equal(t, resp.Digest.Window, "7d")
		assert.Equal(t, resp.Digest.Funnel.Applied, 20)
		assert.Equal(t, resp.Digest.Funnel.Replied, 6)
		assert.Equal(t, resp.Digest.Funnel.Interviewing, 3)
		assert.Equal(t, resp.Digest.Funnel.Offered, 1)
		assert.Equal(t, resp.Digest.Funnel.Rejected, 4)
		assert.Equal(t, resp.Digest.Funnel.Ghosted, 12)
		assert.Equal(t, resp.Digest.GhostCohort.Matured, 20)
		assert.Equal(t, resp.Digest.GhostCohort.MaturedGhosted, 5)
		assert.Equal(t, resp.Digest.GhostCohort.Pending, 10)
		assert.Equal(t, resp.Digest.GhostCohort.Rate, data.Rate(25.0))
	})

	t.Run("invalid window query parameter", func(t *testing.T) {
		status, body := ts.Get(t, "/v1/digest?window=30d")
		assert.Equal(t, status, http.StatusUnprocessableEntity)

		if !strings.Contains(body, "must be one of 7d, 1m, 1y or all") {
			t.Errorf("expected error message for unsupported window, got %s", body)
		}
	})
}
