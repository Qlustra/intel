package formats

import (
	"fmt"

	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
)

type Empty struct {
	inspect.TagFormat
}

func (t Empty) Validate() logic.Check {
	if t.Raw != "" {
		return logic.Nope("empty_tag_expected", fmt.Sprintf("got %q", t.Raw))
	}
	return logic.Ok()
}
