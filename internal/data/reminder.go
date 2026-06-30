package data

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Reminder struct {
	ID          int64     `json:"id"`
	CompanyName string    `json:"company_name"`
	RoleTitle   string    `json:"role_title"`
	InterviewAt time.Time `json:"interview_at"`
}

type ReminderModelInterface interface {
	GetUpcoming(userID int64, limit int) ([]*Reminder, error)
}

type ReminderModel struct {
	pool *pgxpool.Pool
}

// GetUpcoming returns at most `limit` upcoming interviews for authenticated user for 7 days window
func (m ReminderModel) GetUpcoming(userID int64, limit int) ([]*Reminder, error) {
	query := `
		SELECT id, company_name, role_title, interview_at
		FROM applications
		WHERE users_id      = $1
		  AND status        = 'Interviewing'
		  AND interview_at IS NOT NULL
		  AND interview_at  > NOW()
		  AND interview_at <= NOW() + INTERVAL '7 days' 
		ORDER BY interview_at ASC
		LIMIT $2` // has 7 days reminder context window

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reminders []*Reminder

	for rows.Next() {
		var r Reminder
		if err := rows.Scan(&r.ID, &r.CompanyName, &r.RoleTitle, &r.InterviewAt); err != nil {
			return nil, err
		}
		reminders = append(reminders, &r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reminders, nil
}
