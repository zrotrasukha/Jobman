package data

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestStatus_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		status Status
		want   bool
	}{
		{"valid Applied", StatusApplied, true},
		{"valid Interviewing", StatusInterviewing, true},
		{"valid Offered", StatusOffered, true},
		{"valid Rejected", StatusRejected, true},
		{"valid Ghosted", StatusGhosted, true},
		{"invalid lowercase applied", Status("applied"), false},
		{"invalid arbitrary string", Status("custom-status"), false},
		{"invalid empty string", Status(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsValid(); got != tt.want {
				t.Errorf("Status.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatus_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		status  Status
		want    string
		wantErr bool
	}{
		{"marshal Applied", StatusApplied, `"Applied"`, false},
		{"marshal Ghosted", StatusGhosted, `"Ghosted"`, false},
		{"marshal invalid status", Status("invalid"), `"invalid"`, false}, // MarshalJSON just serializes the underlying string
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.status)
			if (err != nil) != tt.wantErr {
				t.Fatalf("json.Marshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if string(got) != tt.want {
				t.Errorf("json.Marshal() = %s, want %s", string(got), tt.want)
			}
		})
	}
}

func TestStatus_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		jsonInput string
		want      Status
		wantErr   bool
		errSubstr string
	}{
		{"unmarshal Applied", `"Applied"`, StatusApplied, false, ""},
		{"unmarshal Interviewing", `"Interviewing"`, StatusInterviewing, false, ""},
		{"unmarshal Offered", `"Offered"`, StatusOffered, false, ""},
		{"unmarshal Rejected", `"Rejected"`, StatusRejected, false, ""},
		{"unmarshal Ghosted", `"Ghosted"`, StatusGhosted, false, ""},
		{"unmarshal invalid lowercase", `"applied"`, "", true, "invalid status: applied"},
		{"unmarshal invalid custom", `"Pending"`, "", true, "invalid status: Pending"},
		{"unmarshal malformed JSON type", `123`, "", true, ""}, // json.Unmarshal should fail on non-string inputs
		{"unmarshal empty quotes", `""`, "", true, "invalid status:"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s Status
			err := json.Unmarshal([]byte(tt.jsonInput), &s)
			if (err != nil) != tt.wantErr {
				t.Fatalf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if s != tt.want {
					t.Errorf("json.Unmarshal() = %s, want %s", s, tt.want)
				}
			} else if tt.errSubstr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errSubstr) {
					t.Errorf("json.Unmarshal() error = %v, expected to contain %q", err, tt.errSubstr)
				}
			}
		})
	}
}
