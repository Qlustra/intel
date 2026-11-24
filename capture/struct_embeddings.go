package capture

import (
	"iter"

	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
	"github.com/qlustra/intel/options"
)

// Embeddings selects fields declared using Go embedding syntax, including
// embedded pointers. A non-embedded field does not join this view just
// because it has the same type. [Extensions] and [Contracts] are narrower
// views of embeddings.
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

// AllAddressable visits visible embeddings, including promoted embeddings,
// breadth-first by default. Visibility does not guarantee live value access.
func (p *Embeddings) AllAddressable(ordering ...options.DepthOrdering) iter.Seq2[string, *EmbeddingInterrogator] {
	return func(yield func(string, *EmbeddingInterrogator) bool) {
		for name, descriptor := range p.fieldsInspector.AllAddressable(ordering...) {
			if !yield(name, p.register(descriptor)) {
				return
			}
		}
	}
}

// Named finds a visible embedding by name, including promotion; missing,
// ambiguous, or filtered fields fail the check.
func (p *Embeddings) Named(name string) (*EmbeddingInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.LookupNamed(name)
	if !check.OK() {
		return nil, check
	}
	return p.register(descriptor), check
}

// Embeds finds a visible embedding of exactly type T. See
// [Embeddings.DeclaredType] for a direct declaration.
func (p *Embeddings) Embeds[T any]() (*EmbeddingInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.Embeds[T]()
	if !check.OK() {
		return nil, check
	}
	return p.register(descriptor), check
}

// Satisfies finds a visible embedding whose field type satisfies interface T.
func (p *Embeddings) Satisfies[T any]() (*EmbeddingInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.Satisfies[T]()
	if !check.OK() {
		return nil, check
	}
	return p.register(descriptor), check
}

// Declared finds an embedding declared directly on the struct by name.
func (p *Embeddings) Declared(name string) (*EmbeddingInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.Declared(name)
	if !check.OK() {
		return nil, check
	}
	return p.register(descriptor), check
}

// DeclaredType finds a directly declared embedding of exactly type T.
func (p *Embeddings) DeclaredType[T any]() (*EmbeddingInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.DeclaredType[T]()
	if !check.OK() {
		return nil, check
	}
	return p.register(descriptor), check
}

// DeclaredContract finds a directly declared embedding whose field type
// satisfies interface T.
func (p *Embeddings) DeclaredContract[T any]() (*EmbeddingInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.DeclaredContract[T]()
	if !check.OK() {
		return nil, check
	}
	return p.register(descriptor), check
}
