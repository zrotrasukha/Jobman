package assert

import (
	"encoding/json"
	"reflect"
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}

func EqualJSON(t *testing.T, actual []byte, expected any) {
	t.Helper()

	expectedB, err := json.MarshalIndent(expected, "", "\t")
	if err != nil {
		t.Fatalf("failed to marshal expected JSON: %v", err)
	}

	var expectedInterface any
	var actualInterface any

	if err := json.Unmarshal(expectedB, &expectedInterface); err != nil {
		t.Fatalf("failed to unmarshal expected JSON: %v", err)
	}

	if err := json.Unmarshal(actual, &actualInterface); err != nil {
		t.Fatalf("failed to unmarshal actual JSON: %v", err)
	}

	if !reflect.DeepEqual(expectedInterface, actualInterface) {
		t.Errorf("\ngot: %s;\nwant: %s", string(actual), string(expectedB))
	}
}
