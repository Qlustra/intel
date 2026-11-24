// Package formats provides tag formats for [inspect.FieldTagInspector.Parse].
// The formats interpret tag text, not the conventions of a particular tag key.
package formats

import (
	"slices"
	"strings"

	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
)

// CSV splits a tag on commas without trimming whitespace, removing empty
// parts, or interpreting quotes. Its optional Validator checks the parts.
type CSV struct {
	inspect.TagFormat
	Parts     []string
	Validator func([]string) logic.Check
}

// Parse replaces Raw and Parts with the supplied text and its comma split.
func (t *CSV) Parse(raw string) {
	t.Raw = raw
	t.Parts = strings.Split(raw, ",")
}

// Validate calls the configured Validator, or succeeds if none is set.
func (t CSV) Validate() logic.Check {
	if t.Validator != nil {
		return t.Validator(t.Parts)
	}
	return logic.Ok()
}

// Contains reports whether a part is exactly equal to value.
func (t *CSV) Contains(value string) bool {
	return slices.Contains(t.Parts, value)
}
