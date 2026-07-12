package data

import (
	"encoding/json"
	"testing"
)

func TestRate_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		rate     Rate
		expected string
	}{
		{
			name:     "integer rate",
			rate:     5,
			expected: `"5%"`,
		},
		{
			name:     "decimal rate",
			rate:     12.34,
			expected: `"12.34%"`,
		},
		{
			name:     "zero rate",
			rate:     0,
			expected: `"0%"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.rate)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(got) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(got))
			}
		})
	}
}

func TestRate_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		jsonInput string
		expected  Rate
		expectErr bool
	}{
		{
			name:      "valid integer percentage",
			jsonInput: `"5%"`,
			expected:  5,
			expectErr: false,
		},
		{
			name:      "valid decimal percentage",
			jsonInput: `"12.34%"`,
			expected:  12.34,
			expectErr: false,
		},
		{
			name:      "invalid format - no percent sign",
			jsonInput: `"5"`,
			expectErr: true,
		},
		{
			name:      "invalid format - space before percent sign",
			jsonInput: `"5 %"`,
			expectErr: true,
		},
		{
			name:      "invalid JSON input - not a quoted string",
			jsonInput: `5%`,
			expectErr: true,
		},
		{
			name:      "invalid numeric format",
			jsonInput: `"abc%"`,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r Rate
			err := json.Unmarshal([]byte(tt.jsonInput), &r)
			if tt.expectErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if r != tt.expected {
					t.Errorf("expected rate %g, got %g", tt.expected, r)
				}
			}
		})
	}
}
