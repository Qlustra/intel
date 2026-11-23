package capture

import (
	"iter"

	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
	"github.com/qlustra/intel/options"
)

type Properties struct {
	fieldsInspector *inspect.FieldsInspector
	fields          *FieldRegistry[PropertyInterrogator]
}

func (p *Properties) AllAddressable(ordering ...options.DepthOrdering) iter.Seq2[string, *PropertyInterrogator] {
	return func(yield func(string, *PropertyInterrogator) bool) {
		for name, descriptor := range p.fieldsInspector.AllAddressable(ordering...) {
			if !yield(name, p.fields.register(descriptor)) {
				return
			}
		}
	}
}

func (p *Properties) AllDeclared() iter.Seq2[string, *PropertyInterrogator] {
	return func(yield func(string, *PropertyInterrogator) bool) {
		for name, descriptor := range p.fieldsInspector.AllDeclared() {
			if !yield(name, p.fields.register(descriptor)) {
				return
			}
		}
	}
}

func (p *Properties) Lookup(name string) (*PropertyInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.LookupNamed(name)
	if !check.OK() {
		return nil, check
	}
	return p.fields.register(descriptor), check
}

func (p *Properties) Declared(name string) (*PropertyInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.Declared(name)
	if !check.OK() {
		return nil, check
	}
	return p.fields.register(descriptor), check
}
