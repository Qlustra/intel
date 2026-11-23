package formats

import (
	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
)

type Single struct {
	inspect.TagFormat
	// Validator is optional; unrestricted values such as help text are valid
	// unless a caller supplies a more specific rule.
	Validator func(string) logic.Check
}

func (t Single) Validate() logic.Check {
	if t.Validator != nil {
		return t.Validator(t.Raw)
	}
	return logic.Ok()
}
