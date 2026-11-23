package formats

import (
	"slices"
	"strings"

	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
)

type CSV struct {
	inspect.TagFormat
	Parts     []string
	Validator func([]string) logic.Check
}

func (t *CSV) Parse(raw string) {
	t.Raw = raw
	t.Parts = strings.Split(raw, ",")
}

func (t CSV) Validate() logic.Check {
	if t.Validator != nil {
		return t.Validator(t.Parts)
	}
	return logic.Ok()
}

func (t *CSV) Contains(value string) bool {
	return slices.Contains(t.Parts, value)
}
