package types

import (
	"iter"

	structEngine "github.com/qlustra/intel/engines/struct"
)

type StructProperties struct {
	fieldsInspector *structEngine.FieldsInspector
	fields          *FieldRegistry[PropertyInterrogator]
}

func (p *StructProperties) AllAddressable() iter.Seq2[string, *PropertyInterrogator] {
	return func(yield func(string, *PropertyInterrogator) bool) {
		for name, descriptor := range p.fieldsInspector.AllAddressable() {
			yield(name, p.fields.register(descriptor))
		}
	}
}

func (p *StructProperties) AllDeclared() iter.Seq2[string, *PropertyInterrogator] {
	return func(yield func(string, *PropertyInterrogator) bool) {
		for name, descriptor := range p.fieldsInspector.AllDeclared() {
			yield(name, p.fields.register(descriptor))
		}
	}
}

func (p *StructProperties) Lookup(name string) (*PropertyInterrogator, bool) {
	descriptor, found := p.fieldsInspector.LookupNamed(name)
	if !found {
		return nil, false
	}
	return p.fields.register(descriptor), true
}

func (p *StructProperties) Declared(name string) (*PropertyInterrogator, bool) {
	descriptor, found := p.fieldsInspector.Declared(name)
	if !found {
		return nil, false
	}
	return p.fields.register(descriptor), true
}
