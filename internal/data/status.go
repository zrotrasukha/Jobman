package data

import (
	"encoding/json"
	"fmt"
)

type Status string

const (
	StatusApplied      = "Applied"
	StatusInterviewing = "Interviewing"
	StatusOffered      = "Offered"
	StatusRejected     = "Rejected"
	StatusGhosted      = "Ghosted"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusApplied, StatusInterviewing, StatusOffered, StatusRejected, StatusGhosted:
		return true
	default:
		return false
	}
}

func (s Status) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf(`"%s"`, s)
	return []byte(jsonValue), nil
}

func (s *Status) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	status := Status(str)
	if !status.IsValid() {
		return fmt.Errorf("invalid status: %s", str)
	}

	*s = status
	return nil
}
