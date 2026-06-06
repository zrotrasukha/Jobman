package validator

import (
	"regexp"
	"slices"
)

var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$")
)

// Validator struct hold the Errors map, which is used to store validtion errors
type Validator struct {
	Errors map[string]string
}

// Validator returns a new instance of Validator struct with empty Errors map.
func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

// AddError adds an error message to the Errors map for a specific key. If an error message already exists for the given key, it does not overwrite it
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Checkfield checks for any condition beign false. If false, it adds the error message to the Errors map for the specified key.
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// Valid returns length of errors in Errors map.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// PermittedValues is a generic function that checks if a given value in given list of permitted values.
func PermittedValues[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// Matches checks if a string matches a given regular expression pattern. It returns true if the string matches the pattern, and false otherwise.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
