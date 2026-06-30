package mocks

import (
	"time"

	"github.com/zrotrasukha/jobman/internal/data"
)

var FixedInterviewDate = time.Date(2026, 8, 15, 10, 0, 0, 0, time.UTC)

type MockReminderModel struct {
	GetUpcomingFunc func(userID int64, limit int) ([]*data.Reminder, error)
}

func (m MockReminderModel) GetUpcoming(userID int64, limit int) ([]*data.Reminder, error) {
	if m.GetUpcomingFunc != nil {
		return m.GetUpcomingFunc(userID, limit)
	}

	return []*data.Reminder{
		{
			ID:          1,
			CompanyName: "Test Company",
			RoleTitle:   "Software Engineer",
			InterviewAt: FixedInterviewDate,
		},
	}, nil
}
