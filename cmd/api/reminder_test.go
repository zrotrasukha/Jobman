package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/zrotrasukha/jobman/internal/assert"
	"github.com/zrotrasukha/jobman/internal/data"
	"github.com/zrotrasukha/jobman/internal/data/mocks"
)

func TestListRemindersHandler(t *testing.T) {
	interviewDate := time.Date(2026, 8, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name                string
		wantStatus          int
		mockReminderModel   mocks.MockReminderModel
		wantBody            string
	}{
		{
			name:       "returns upcoming reminders",
			wantStatus: http.StatusOK,
			// default mock returns one reminder
			mockReminderModel: mocks.MockReminderModel{},
			wantBody: `{
				"reminders": [
					{
						"id": 1,
						"company_name": "Test Company",
						"role_title": "Software Engineer",
						"interview_at": "2026-08-15T10:00:00Z"
					}
				]
			}`,
		},
		{
			name:       "returns empty array when no upcoming interviews",
			wantStatus: http.StatusOK,
			mockReminderModel: mocks.MockReminderModel{
				GetUpcomingFunc: func(userID int64, limit int) ([]*data.Reminder, error) {
					return nil, nil
				},
			},
			wantBody: `{"reminders": []}`,
		},
		{
			name:       "respects limit of 10",
			wantStatus: http.StatusOK,
			mockReminderModel: mocks.MockReminderModel{
				GetUpcomingFunc: func(userID int64, limit int) ([]*data.Reminder, error) {
					if limit != 10 {
						t.Errorf("expected limit 10, got %d", limit)
					}
					return []*data.Reminder{
						{ID: 1, CompanyName: "Alpha", RoleTitle: "SWE", InterviewAt: interviewDate},
					}, nil
				},
			},
			wantBody: `{
				"reminders": [
					{
						"id": 1,
						"company_name": "Alpha",
						"role_title": "SWE",
						"interview_at": "2026-08-15T10:00:00Z"
					}
				]
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			models := mocks.NewMockModels()
			models.Reminder = tt.mockReminderModel

			ts := newTestServer(t, models)
			sc, body := ts.Get(t, "/v1/reminders")

			assert.Equal(t, sc, tt.wantStatus)
			assert.EqualJSON(t, []byte(body), []byte(tt.wantBody))
		})
	}
}
