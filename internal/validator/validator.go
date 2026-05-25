package validator

import "slices"

// Validator struct is used to collect validation errors for various fields. It contains a map where the keys are the field names and the values are the corresponding error messages. The struct provides methods to add errors, check field validity, and determine if the overall validation is successful (i.e., if there are no errors).
type Validator struct {
	Errors map[string]string
}

// New creates and returns a new instance of the Validator struct with an initialized Errors map. This allows for collecting validation errors as they are added during the validation process.
func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

// AddError adds an error message to the Errors map for a specific field (key). If an error for that field already exists, it does not overwrite it, ensuring that only the first error message for each field is retained.
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// CheckField evaluates a boolean condition (ok) for a specific field (key). If the condition is false, it adds an error message to the Errors map for that field using the AddError method. This method is typically used to validate individual fields and collect any errors that occur during the validation process.
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// Valid checks if there are any validation errors collected in the Errors map. It returns true if there are no errors (i.e., the length of the Errors map is zero), indicating that the validation was successful. If there are any errors, it returns false, indicating that the validation failed.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// PermittedValues is a generic function that checks if a given value is included in a list of permitted values. It uses the slices.Contains function from the standard library to determine if the value exists in the provided slice of permitted values. The function returns true if the value is found in the slice, and false otherwise. This is commonly used for validating that a field's value is one of a predefined set of acceptable options.
func PermittedValues[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}
