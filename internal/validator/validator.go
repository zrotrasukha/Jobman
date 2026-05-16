package validator

type Validator struct {
	errors map[string]string
}

func New() *Validator {
	return &Validator{
		errors: make(map[string]string),
	}
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.errors[key]; !exists {
		v.errors[key] = message
	}
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func (v *Validator) Valid() bool {
	return len(v.errors) == 0
}
