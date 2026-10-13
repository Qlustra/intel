package formats

type Single struct {
	TagFormat
	// Validator is optional; unrestricted values such as help text are valid
	// unless a caller supplies a more specific rule.
	Validator func(string) error
}

func (t Single) Validate() error {
	if t.Validator != nil {
		return t.Validator(t.Raw)
	}
	return nil
}
