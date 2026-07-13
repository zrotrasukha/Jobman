package data

import (
	"encoding/json"
	"fmt"
)

// Status represents the lifecycle stage of a job application.
type Status string

const (
	StatusApplied      = "Applied"
	StatusInterviewing = "Interviewing"
	StatusOffered      = "Offered"
	StatusRoundCleared = "RoundCleared" // repeatable, mid-process
	StatusSelected     = "Selected"     // earned via full pipeline
	StatusDeclined     = "Declined"     // user said no to offer 
	StatusRejected     = "Rejected"
	StatusGhosted      = "Ghosted"
)

// IsValid returns true if the status matches a predefined valid value.
func (s Status) IsValid() bool {
	switch s {
	case StatusApplied, StatusInterviewing, StatusOffered, StatusRejected, StatusGhosted, StatusRoundCleared, StatusSelected, StatusDeclined:
		return true
	default:
		return false
	}
}

// MarshalJSON encodes the Status as a JSON string.
func (s Status) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf(`"%s"`, s)
	return []byte(jsonValue), nil
}

// UnmarshalJSON decodes a JSON string into a Status and validates that it is a known value.
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
