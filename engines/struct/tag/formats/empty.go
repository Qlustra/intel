package formats

import "fmt"

type Empty struct {
	TagFormat
}

func (t Empty) Validate() error {
	if t.Raw != "" {
		return fmt.Errorf("expected an empty tag value, got %q", t.Raw)
	}
	return nil
}
