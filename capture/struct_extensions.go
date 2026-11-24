package capture

import (
	"iter"

	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
	"github.com/qlustra/intel/options"
)

// Extensions selects embedded fields of struct kind. An embedded pointer to
// a struct is in [Embeddings] but not in Extensions; this view does not
// dereference the field type.
type Extensions struct {
	fieldsInspector *inspect.FieldsInspector
	fields          *FieldRegistry[ExtensionInterrogator]
}

// AllAddressable visits visible extensions, including promoted extensions,
// breadth-first by default. Visibility does not guarantee live value access.
func (p *Extensions) AllAddressable(ordering ...options.DepthOrdering) iter.Seq2[string, *ExtensionInterrogator] {
	return func(yield func(string, *ExtensionInterrogator) bool) {
		for name, descriptor := range p.fieldsInspector.AllAddressable(ordering...) {
			if !yield(name, p.fields.register(descriptor)) {
				return
			}
		}
	}
}

// Lookup finds a visible extension by name, including promotion.
func (p *Extensions) Lookup(name string) (*ExtensionInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.LookupNamed(name)
	if !check.OK() {
		return nil, check
	}
	return p.fields.register(descriptor), check
}

// Declared finds an extension declared directly on the struct by name.
func (p *Extensions) Declared(name string) (*ExtensionInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.Declared(name)
	if !check.OK() {
		return nil, check
	}
	return p.fields.register(descriptor), check
}

// Embeds finds a visible extension of exactly type T.
func (p *Extensions) Embeds[T any]() (*ExtensionInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.Embeds[T]()
	if !check.OK() {
		return nil, check
	}
	return p.fields.register(descriptor), check
}

// Satisfies finds a visible extension whose field type satisfies interface T.
func (p *Extensions) Satisfies[T any]() (*ExtensionInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.Satisfies[T]()
	if !check.OK() {
		return nil, check
	}
	return p.fields.register(descriptor), check
}

// DeclaredType finds a directly declared extension of exactly type T.
func (p *Extensions) DeclaredType[T any]() (*ExtensionInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.DeclaredType[T]()
	if !check.OK() {
		return nil, check
	}
	return p.fields.register(descriptor), check
}
