package formats

import (
	"fmt"

	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
)

// Empty accepts only an empty tag value. Presence is checked separately by
// [inspect.FieldTagInspector.Parse].
type Empty struct {
	inspect.TagFormat
}

// Validate fails if the parsed tag value is non-empty.
func (t Empty) Validate() logic.Check {
	if t.Raw != "" {
		return logic.Nope("empty_tag_expected", fmt.Sprintf("got %q", t.Raw))
	}
	return logic.Ok()
}
