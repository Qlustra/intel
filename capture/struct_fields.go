package capture

import (
	"iter"

	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
	"github.com/qlustra/intel/options"
)

type StructFields struct {
	fieldsInspector *inspect.FieldsInspector
	fields          *FieldRegistry[FieldInterrogator]
	properties      *FieldRegistry[PropertyInterrogator]
	extensions      *FieldRegistry[ExtensionInterrogator]
	contracts       *FieldRegistry[ContractInterrogator]
	embeddings      *FieldRegistry[EmbeddingInterrogator]
}

func (p *StructFields) register(descriptor *inspect.FieldDescriptor) *FieldInterrogator {
	interrogator := p.fields.register(descriptor)
	interrogator.properties = p.properties
	interrogator.extensions = p.extensions
	interrogator.contracts = p.contracts
	interrogator.embeddings = p.embeddings
	return interrogator
}

func (p *StructFields) AllAddressable(ordering ...options.DepthOrdering) iter.Seq2[string, *FieldInterrogator] {
	return func(yield func(string, *FieldInterrogator) bool) {
		for name, descriptor := range p.fieldsInspector.AllAddressable(ordering...) {
			if !yield(name, p.register(descriptor)) {
				return
			}
		}
	}
}

func (p *StructFields) AllDeclared() iter.Seq2[string, *FieldInterrogator] {
	return func(yield func(string, *FieldInterrogator) bool) {
		for name, descriptor := range p.fieldsInspector.AllDeclared() {
			if !yield(name, p.register(descriptor)) {
				return
			}
		}
	}
}

func (p *StructFields) Lookup(name string) (*FieldInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.LookupNamed(name)
	if !check.OK() {
		return nil, check
	}
	return p.register(descriptor), check
}

func (p *StructFields) Declared(name string) (*FieldInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.Declared(name)
	if !check.OK() {
		return nil, check
	}
	return p.register(descriptor), check
}
