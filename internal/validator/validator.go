package validator

// Validator contains a map of strings to strings representing errors
type Validator struct {
	Errors map[string]string
}

// Returns an initialized Validator
func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

// Returns a boolean value. False is invalid, meaning there are errors, and true is valid
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// Adds an error to Validator with the provided key and value, if it doesn't already exists
func (v *Validator) AddError(key string, value string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = value
	}
}

// Evaluates the expression, and if !ok, adds an erroor to Validator
func (v *Validator) Check(ok bool, key string, value string) {
	if !ok {
		v.AddError(key, value)
	}
}
