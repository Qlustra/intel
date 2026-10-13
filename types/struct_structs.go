package types

import (
	structEngine "github.com/qlustra/intel/engines/struct"
)

type EmbeddedStructs struct {
	fieldsInspector *structEngine.FieldsInspector
	fields          *FieldRegistry[StructInterrogator]
}

func (p *EmbeddedStructs) Lookup(name string) (*StructInterrogator, bool) {
	descriptor, found := p.fieldsInspector.LookupNamed(name)
	if !found {
		return nil, false
	}
	return p.fields.register(descriptor), true
}

func (p *EmbeddedStructs) Declared(name string) (*StructInterrogator, bool) {
	descriptor, found := p.fieldsInspector.Declared(name)
	if !found {
		return nil, false
	}
	return p.fields.register(descriptor), true
}

func (p *EmbeddedStructs) OfType[T any]() (*StructInterrogator, bool, bool) {
	descriptor, found, ambiguous := p.fieldsInspector.LookupType[T]()
	if !found {
		return nil, false, false
	}
	return p.fields.register(descriptor), true, ambiguous
}

func (p *EmbeddedStructs) Satisfies[T any]() (*StructInterrogator, bool, bool) {
	descriptor, found, ambiguous := p.fieldsInspector.LookupContract[T]()
	if !found {
		return nil, false, ambiguous
	}
	return p.fields.register(descriptor), true, ambiguous
}
