package types

import (
	"iter"

	structEngine "github.com/qlustra/intel/engines/struct"
)

type StructEmbeddings struct {
	fieldsInspector *structEngine.FieldsInspector
	fields          *FieldRegistry[EmbeddingInterrogator]
	structs         *FieldRegistry[StructInterrogator]
	interfaces      *FieldRegistry[InterfaceInterrogator]
}

func (p *StructEmbeddings) register(descriptor *structEngine.FieldDescriptor) *EmbeddingInterrogator {
	interrogator := p.fields.register(descriptor)
	interrogator.structs = p.structs
	interrogator.interfaces = p.interfaces
	return interrogator
}

func (p *StructEmbeddings) AllAddressable() iter.Seq2[string, *EmbeddingInterrogator] {
	return func(yield func(string, *EmbeddingInterrogator) bool) {
		for name, descriptor := range p.fieldsInspector.AllAddressable() {
			yield(name, p.register(descriptor))
		}
	}
}

func (p *StructEmbeddings) Named(name string) (*EmbeddingInterrogator, bool) {
	descriptor, found := p.fieldsInspector.LookupNamed(name)
	if !found {
		return nil, false
	}
	return p.register(descriptor), true
}

func (p *StructEmbeddings) OfType[T any]() (*EmbeddingInterrogator, bool, bool) {
	descriptor, found, ambiguous := p.fieldsInspector.LookupType[T]()
	if !found {
		return nil, found, ambiguous
	}
	return p.register(descriptor), true, ambiguous
}

func (p *StructEmbeddings) Satisfies[T any]() (*EmbeddingInterrogator, bool, bool) {
	descriptor, found, ambiguous := p.fieldsInspector.LookupContract[T]()
	if !found {
		return nil, false, ambiguous
	}
	return p.register(descriptor), true, ambiguous
}

func (p *StructEmbeddings) Declared(name string) (*EmbeddingInterrogator, bool) {
	descriptor, found := p.fieldsInspector.Declared(name)
	if !found {
		return nil, false
	}
	return p.register(descriptor), true
}

func (p *StructEmbeddings) DeclaredType[T any]() (*EmbeddingInterrogator, bool) {
	descriptor, found := p.fieldsInspector.DeclaredType[T]()
	if !found {
		return nil, false
	}
	return p.register(descriptor), true
}

func (p *StructEmbeddings) DeclaredContract[T any]() (*EmbeddingInterrogator, bool) {
	descriptor, found := p.fieldsInspector.DeclaredContract[T]()
	if !found {
		return nil, false
	}
	return p.register(descriptor), true
}
