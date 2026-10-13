package tag

import (
	"fmt"
	"reflect"

	"github.com/qlustra/intel/engines/struct/tag/formats"
)

type FieldTagInspector struct {
	tags reflect.StructTag
}

func CaptureFieldTag(tags reflect.StructTag) *FieldTagInspector {
	interrogator := &FieldTagInspector{
		tags: tags,
	}
	return interrogator
}

func (i *FieldTagInspector) IsPresent(name string) bool {
	_, present := i.Lookup(name)
	return present
}

func (i *FieldTagInspector) Lookup(name string) (string, bool) {
	return i.tags.Lookup(name)
}

func (i *FieldTagInspector) Raw(name string) string {
	raw, present := i.Lookup(name)
	if !present {
		return ""
	}
	return raw
}

func (i *FieldTagInspector) Parse[
	AS any,
	FORMAT interface {
		*AS
		formats.TagFormatLike
	},
](name string) (*Tag[AS], bool, error) {
	raw, present := i.Lookup(name)
	if !present {
		return nil, false, nil
	}

	format := new(AS)
	parser := FORMAT(format)
	parser.Parse(raw)
	if err := parser.Validate(); err != nil {
		return nil, true, fmt.Errorf("tag %q: %w", name, err)
	}

	return &Tag[AS]{Name: name, Format: format}, true, nil
}
