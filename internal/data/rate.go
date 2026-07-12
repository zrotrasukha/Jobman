package data

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidRateFormat = fmt.Errorf("invalid rate format")

type Rate float64

// MarshalJSON serializes the Rate to a percentage string (e.g., 5.5 to "5.5%").
func (r Rate) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%g%%", r)
	quotedValue := strconv.Quote(jsonValue)
	return []byte(quotedValue), nil
}

// UnmarshalJSON deserializes a percentage string (e.g., "5.5%") back to Rate.
// A bare `null` is rejected here deliberately — for *Rate fields, encoding/json
// short-circuits null and never calls this method at all, so reaching this
// function with the literal bytes `null` would mean something unexpected
// happened upstream (e.g. this method got called on a non-pointer Rate field
// that received null). Fail loudly rather than silently defaulting to 0.
func (r *Rate) UnmarshalJSON(jsonValue []byte) error {
	if bytes.Equal(jsonValue, []byte("null")) {
		return ErrInvalidRateFormat
	}
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRateFormat
	}
	if !strings.HasSuffix(unquotedJSONValue, "%") {
		return ErrInvalidRateFormat
	}
	valStr := strings.TrimSuffix(unquotedJSONValue, "%")
	f, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return ErrInvalidRateFormat
	}
	*r = Rate(f)
	return nil
}
