package formats

type TagFormat struct {
	Raw string
}

func (t *TagFormat) Parse(raw string) {
	t.Raw = raw
}

// Validate accepts opaque tag text. Specific formats supply their own rules.
func (t TagFormat) Validate() error {
	return nil
}
