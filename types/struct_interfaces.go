package types

import (
	structEngine "github.com/qlustra/intel/engines/struct"
)

type EmbeddedInterfaces struct {
	fieldsInspector *structEngine.FieldsInspector
	fields          *FieldRegistry[InterfaceInterrogator]
}

func (p *EmbeddedInterfaces) Lookup(name string) (*InterfaceInterrogator, bool) {
	descriptor, found := p.fieldsInspector.LookupNamed(name)
	if !found {
		return nil, false
	}
	return p.fields.register(descriptor), true
}

func (p *EmbeddedInterfaces) Declared(name string) (*InterfaceInterrogator, bool) {
	descriptor, found := p.fieldsInspector.Declared(name)
	if !found {
		return nil, false
	}
	return p.fields.register(descriptor), true
}

func (p *EmbeddedInterfaces) OfType[T any]() (*InterfaceInterrogator, bool, bool) {
	descriptor, found, ambiguous := p.fieldsInspector.LookupType[T]()
	if !found {
		return nil, false, false
	}
	return p.fields.register(descriptor), true, ambiguous
}

func (p *EmbeddedInterfaces) Satisfies[T any]() (*InterfaceInterrogator, bool, bool) {
	descriptor, found, ambiguous := p.fieldsInspector.LookupContract[T]()
	if !found {
		return nil, false, ambiguous
	}
	return p.fields.register(descriptor), true, ambiguous
}
