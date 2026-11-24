package capture

import (
	"iter"

	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
	"github.com/qlustra/intel/options"
)

// Properties selects non-embedded, non-blank fields. A property may be
// promoted through an embedding; use [Properties.Declared] to require a
// direct declaration. Field type alone does not make a field an embedding.
type Properties struct {
	fieldsInspector *inspect.FieldsInspector
	fields          *FieldRegistry[PropertyInterrogator]
}

// AllAddressable visits the visible properties in the current view, including
// promoted properties. It visits breadth-first by default; changing traversal
// order does not change which properties are visible. This is type inspection,
// not a check that a particular instance has accessible storage for a property.
func (p *Properties) AllAddressable(ordering ...options.DepthOrdering) iter.Seq2[string, *PropertyInterrogator] {
	return func(yield func(string, *PropertyInterrogator) bool) {
		for name, descriptor := range p.fieldsInspector.AllAddressable(ordering...) {
			if !yield(name, p.fields.register(descriptor)) {
				return
			}
		}
	}
}

// AllDeclared visits properties declared directly on the struct in this view.
func (p *Properties) AllDeclared() iter.Seq2[string, *PropertyInterrogator] {
	return func(yield func(string, *PropertyInterrogator) bool) {
		for name, descriptor := range p.fieldsInspector.AllDeclared() {
			if !yield(name, p.fields.register(descriptor)) {
				return
			}
		}
	}
}

// Lookup finds a visible property by name, including promotion. Missing,
// ambiguous, or filtered properties fail the check.
func (p *Properties) Lookup(name string) (*PropertyInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.LookupNamed(name)
	if !check.OK() {
		return nil, check
	}
	return p.fields.register(descriptor), check
}

// Declared looks up a property declared directly on the struct. A promoted
// property does not match.
func (p *Properties) Declared(name string) (*PropertyInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.Declared(name)
	if !check.OK() {
		return nil, check
	}
	return p.fields.register(descriptor), check
}
