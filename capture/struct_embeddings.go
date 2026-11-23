package capture

import (
	"iter"

	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
	"github.com/qlustra/intel/options"
)

type Embeddings struct {
	fieldsInspector *inspect.FieldsInspector
	fields          *FieldRegistry[EmbeddingInterrogator]
	extensions      *FieldRegistry[ExtensionInterrogator]
	contracts       *FieldRegistry[ContractInterrogator]
}

func (p *Embeddings) register(descriptor *inspect.FieldDescriptor) *EmbeddingInterrogator {
	interrogator := p.fields.register(descriptor)
	interrogator.structs = p.extensions
	interrogator.contracts = p.contracts
	return interrogator
}

func (p *Embeddings) AllAddressable(ordering ...options.DepthOrdering) iter.Seq2[string, *EmbeddingInterrogator] {
	return func(yield func(string, *EmbeddingInterrogator) bool) {
		for name, descriptor := range p.fieldsInspector.AllAddressable(ordering...) {
			if !yield(name, p.register(descriptor)) {
				return
			}
		}
	}
}

func (p *Embeddings) Named(name string) (*EmbeddingInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.LookupNamed(name)
	if !check.OK() {
		return nil, check
	}
	return p.register(descriptor), check
}

func (p *Embeddings) Embeds[T any]() (*EmbeddingInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.Embeds[T]()
	if !check.OK() {
		return nil, check
	}
	return p.register(descriptor), check
}

func (p *Embeddings) Satisfies[T any]() (*EmbeddingInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.Satisfies[T]()
	if !check.OK() {
		return nil, check
	}
	return p.register(descriptor), check
}

func (p *Embeddings) Declared(name string) (*EmbeddingInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.Declared(name)
	if !check.OK() {
		return nil, check
	}
	return p.register(descriptor), check
}

func (p *Embeddings) DeclaredType[T any]() (*EmbeddingInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.DeclaredType[T]()
	if !check.OK() {
		return nil, check
	}
	return p.register(descriptor), check
}

func (p *Embeddings) DeclaredContract[T any]() (*EmbeddingInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.DeclaredContract[T]()
	if !check.OK() {
		return nil, check
	}
	return p.register(descriptor), check
}
