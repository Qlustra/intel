package capture

import (
	"iter"

	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
	"github.com/qlustra/intel/options"
)

// StructFields is the general member view. It can describe embeddings and
// promoted fields. Use [Properties] to select non-embedded, non-blank fields
// or [Embeddings] to select embedded fields.
type StructFields struct {
	fieldsInspector *inspect.FieldsInspector
	fields          *FieldRegistry[FieldInterrogator]
	properties      *FieldRegistry[PropertyInterrogator]
	extensions      *FieldRegistry[ExtensionInterrogator]
	contracts       *FieldRegistry[ContractInterrogator]
	embeddings      *Embeddings
}

func (p *StructFields) register(descriptor *inspect.FieldDescriptor) *FieldInterrogator {
	interrogator := p.fields.register(descriptor)
	interrogator.properties = p.properties
	interrogator.extensions = p.extensions
	interrogator.contracts = p.contracts
	interrogator.embeddings = p.embeddings
	return interrogator
}

// AllAddressable visits visible fields, including promoted fields, in the
// requested depth order (breadth-first by default). This describes type
// visibility, not access to storage on a particular value.
func (p *StructFields) AllAddressable(ordering ...options.DepthOrdering) iter.Seq2[string, *FieldInterrogator] {
	return func(yield func(string, *FieldInterrogator) bool) {
		for name, descriptor := range p.fieldsInspector.AllAddressable(ordering...) {
			if !yield(name, p.register(descriptor)) {
				return
			}
		}
	}
}

// AllDeclared visits fields declared directly on the struct in this view.
func (p *StructFields) AllDeclared() iter.Seq2[string, *FieldInterrogator] {
	return func(yield func(string, *FieldInterrogator) bool) {
		for name, descriptor := range p.fieldsInspector.AllDeclared() {
			if !yield(name, p.register(descriptor)) {
				return
			}
		}
	}
}

// Lookup finds a visible field by name, including promotion. Missing,
// ambiguous, or filtered fields fail the check; see [StructFields.Declared].
func (p *StructFields) Lookup(name string) (*FieldInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.LookupNamed(name)
	if !check.OK() {
		return nil, check
	}
	return p.register(descriptor), check
}

// Declared finds a field declared directly on the struct, not a promoted one.
func (p *StructFields) Declared(name string) (*FieldInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.Declared(name)
	if !check.OK() {
		return nil, check
	}
	return p.register(descriptor), check
}
