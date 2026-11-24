package inspect

import (
	"fmt"
	"reflect"

	"github.com/qlustra/intel/logic"
)

const (
	TagNotPresentCode    logic.CheckCode = "tag_not_present"
	InvalidTagFormatCode logic.CheckCode = "tag_invalid"
)

// TagValidatorLike validates one parsed value without inspecting a declaration or
// consulting its symbol ledger.
type TagValidatorLike interface {
	Validate() logic.Check
}

type TagFormatLike interface {
	TagValidatorLike
	Parse(raw string)
}

type Tag[F any] struct {
	Name   string
	Format *F
}

type TagFormat struct {
	Raw string
}

func (t *TagFormat) Parse(raw string) {
	t.Raw = raw
}

// Validate accepts opaque tag text. Specific formats supply their own rules.
func (t TagFormat) Validate() logic.Check {
	return logic.Ok()
}

type FieldTagInspector struct {
	tags reflect.StructTag
}

// CaptureFieldTag creates an inspector over tags. Lookup follows
// reflect.StructTag semantics; it does not assign meaning to tag keys.
func CaptureFieldTag(tags reflect.StructTag) *FieldTagInspector {
	interrogator := &FieldTagInspector{
		tags: tags,
	}
	return interrogator
}

// IsPresent succeeds when name is present, including when its value is empty.
func (i *FieldTagInspector) IsPresent(name string) logic.Check {
	_, present := i.Lookup(name)
	if present {
		return logic.Ok()
	}
	return logic.Nope(TagNotPresentCode, name)
}

// Lookup returns the value and presence of name, distinguishing an absent key
// from a key explicitly set to the empty string. See [FieldTagInspector.Raw]
// when that distinction is not needed and [FieldTagInspector.Parse] for a
// validated format.
func (i *FieldTagInspector) Lookup(name string) (string, bool) {
	return i.tags.Lookup(name)
}

// Raw returns the tag value, or an empty string when name is absent. Use
// Lookup when absence must be distinguished from a present empty value.
func (i *FieldTagInspector) Raw(name string) string {
	raw, present := i.Lookup(name)
	if !present {
		return ""
	}
	return raw
}

// Parse looks up name, parses its raw value with FORMAT, and validates the
// resulting format. Missing keys and invalid formats return a failed check;
// use the result only when the check is OK. FORMAT is a pointer to AS that
// implements TagFormatLike.
func (i *FieldTagInspector) Parse[
	AS any,
	FORMAT interface {
		*AS
		TagFormatLike
	},
](name string) (*Tag[AS], logic.Check) {
	raw, present := i.Lookup(name)
	if !present {
		return nil, logic.Nope(TagNotPresentCode, name)
	}

	format := new(AS)
	parser := FORMAT(format)
	parser.Parse(raw)
	if check := parser.Validate(); !check.OK() {
		typ := CaptureType(reflect.TypeFor[AS]())
		return nil, logic.NopeBecause(InvalidTagFormatCode, fmt.Sprintf("%s is not a valid tag format: %s", name, typ.TypeName()), check)
	}

	return &Tag[AS]{Name: name, Format: format}, logic.Ok()
}

func (i *FieldTagInspector) String() string {
	return string(i.tags)
}
