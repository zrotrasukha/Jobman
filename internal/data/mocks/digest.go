package mocks

import (
	"time"

	"github.com/zrotrasukha/jobman/internal/data"
)

type MockDigestModel struct {
	GetDigestFunc func(userID int64, from, to time.Time) (*data.Digest, error)
}

func (m MockDigestModel) GetDigest(userID int64, from, to time.Time) (*data.Digest, error) {
	if m.GetDigestFunc != nil {
		return m.GetDigestFunc(userID, from, to)
	}

	rate := data.Rate(25.0)

	return &data.Digest{
		Window: "7d",
		From:   from,
		To:     to,
		Funnel: &data.Funnel{
			Applied:      20,
			Replied:      6,
			Interviewing: 3,
			Offered:      1,
			RoundCleared: 2,
			Selected:     1,
			Declined:     1,
			Rejected:     4,
			Ghosted:      12,
		},
		GhostCohort: &data.GhostCohort{
			Matured:        20,
			MaturedGhosted: 5,
			Rate:           &rate,
		},
	}, nil
}
