package data

import (
	"encoding/json"
	"fmt"
)

// Status is a custom type that represents the status of a job application. It is defined as a string type, and can take on one of several predefined values such as "Applied", "Interviewing", "Offered", "Rejected", or "Ghosted". The IsValid method checks if a given status value is valid, while the MarshalJSON and UnmarshalJSON methods allow for custom JSON encoding and decoding of the Status type.
type Status string

const (
	StatusApplied      = "Applied"
	StatusInterviewing = "Interviewing"
	StatusOffered      = "Offered"
	StatusRejected     = "Rejected"
	StatusGhosted      = "Ghosted"
)

// IsValid checks if the Status value is one of the predefined valid statuses. It returns true if the status is valid, and false otherwise.
func (s Status) IsValid() bool {
	switch s {
	case StatusApplied, StatusInterviewing, StatusOffered, StatusRejected, StatusGhosted:
		return true
	default:
		return false
	}
}

// MarshalJSON is a custom JSON marshaling method for the Status type. It converts the Status value to a JSON string format, which allows it to be easily included in JSON responses when encoding job application data.
func (s Status) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf(`"%s"`, s)
	return []byte(jsonValue), nil
}

// UnmarshalJSON is a custom JSON unmarshaling method for the Status type. It takes a JSON byte slice as input, attempts to unmarshal it into a string, and then checks if the resulting string is a valid Status value. If the value is valid, it assigns it to the Status variable; otherwise, it returns an error indicating that the status is invalid.
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
