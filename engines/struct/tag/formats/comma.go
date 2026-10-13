package formats

import "strings"

type CSV struct {
	TagFormat
	Parts     []string
	Validator func([]string) error
}

func (t *CSV) Parse(raw string) {
	t.Raw = raw
	t.Parts = strings.Split(raw, ",")
}

func (t CSV) Validate() error {
	if t.Validator != nil {
		return t.Validator(t.Parts)
	}
	return nil
}
