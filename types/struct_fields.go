package types

import (
	"iter"

	structEngine "github.com/qlustra/intel/engines/struct"
)

type StructFields struct {
	fieldsInspector *structEngine.FieldsInspector
	fields          *FieldRegistry[FieldInterrogator]
	properties      *FieldRegistry[PropertyInterrogator]
	structs         *FieldRegistry[StructInterrogator]
	interfaces      *FieldRegistry[InterfaceInterrogator]
	embeddings      *FieldRegistry[EmbeddingInterrogator]
}

func (p *StructFields) register(descriptor *structEngine.FieldDescriptor) *FieldInterrogator {
	interrogator := p.fields.register(descriptor)
	interrogator.properties = p.properties
	interrogator.structs = p.structs
	interrogator.embeddings = p.embeddings
	return interrogator
}

func (p *StructFields) AllAddressable() iter.Seq2[string, *FieldInterrogator] {
	return func(yield func(string, *FieldInterrogator) bool) {
		for name, descriptor := range p.fieldsInspector.AllAddressable() {
			yield(name, p.register(descriptor))
		}
	}
}

func (p *StructFields) AllDeclared() iter.Seq2[string, *FieldInterrogator] {
	return func(yield func(string, *FieldInterrogator) bool) {
		for name, descriptor := range p.fieldsInspector.AllDeclared() {
			yield(name, p.register(descriptor))
		}
	}
}

func (p *StructFields) Lookup(name string) (*FieldInterrogator, bool) {
	descriptor, found := p.fieldsInspector.LookupNamed(name)
	if !found {
		return nil, false
	}
	return p.register(descriptor), true
}

func (p *StructFields) Declared(name string) (*FieldInterrogator, bool) {
	descriptor, found := p.fieldsInspector.Declared(name)
	if !found {
		return nil, false
	}
	return p.register(descriptor), true
}
