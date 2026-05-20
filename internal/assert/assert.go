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

	var actualInterface any
	var expectedInterface any

	// find the right type shit
	var expectedBytes []byte
	switch v := expected.(type) {
	case string:
		expectedBytes = []byte(v)
	case []byte:
		expectedBytes = v
	default: // in case if the expected is map typ
		var err error
		expectedBytes, err = json.Marshal(v)
		if err != nil {
			t.Fatalf("faild to marshal expected type %T: %v", v, err)
		}
	}

	// Be the final judgement fair
	if err := json.Unmarshal(actual, &actualInterface); err != nil {
		t.Fatalf("failde to marshal actual json: %v", err)
	}
	if err := json.Unmarshal(expectedBytes, &expectedInterface); err != nil {
		t.Fatalf("failde to marshal actual json: %v", err)
	}

	if !reflect.DeepEqual(actualInterface, expectedInterface) {
		t.Errorf("\ngot: %s\nwant: %s", string(actual), string(expectedBytes))
	}

}
