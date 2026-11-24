package capture

import (
	"iter"

	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
	"github.com/qlustra/intel/options"
)

// Contracts selects embedded fields of interface kind. A named,
// non-embedded interface field belongs to [Properties] instead. This view
// does not report every interface the enclosing struct satisfies, or whether
// an embedded interface has a non-nil value on a particular instance.
type Contracts struct {
	fieldsInspector *inspect.FieldsInspector
	fields          *FieldRegistry[ContractInterrogator]
}

// AllAddressable visits visible contracts, including promoted contracts,
// breadth-first by default. Visibility does not guarantee live value access.
func (p *Contracts) AllAddressable(ordering ...options.DepthOrdering) iter.Seq2[string, *ContractInterrogator] {
	return func(yield func(string, *ContractInterrogator) bool) {
		for name, descriptor := range p.fieldsInspector.AllAddressable(ordering...) {
			if !yield(name, p.fields.register(descriptor)) {
				return
			}
		}
	}
}

// Lookup finds a visible contract by name, including promotion.
func (p *Contracts) Lookup(name string) (*ContractInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.LookupNamed(name)
	if !check.OK() {
		return nil, check
	}
	return p.fields.register(descriptor), check
}

// Declared finds a contract declared directly on the struct by name.
func (p *Contracts) Declared(name string) (*ContractInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.Declared(name)
	if !check.OK() {
		return nil, check
	}
	return p.fields.register(descriptor), check
}

// Embeds looks for a visible contract whose field type is exactly T. Use
// [Contracts.Satisfies] to match by interface satisfaction instead.
func (p *Contracts) Embeds[T any]() (*ContractInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.Embeds[T]()
	if !check.OK() {
		return nil, check
	}
	return p.fields.register(descriptor), check
}

// Satisfies looks for a visible contract whose field type satisfies the
// interface T. It checks the embedded field's type, not the enclosing struct's
// method set or the value stored in the field.
func (p *Contracts) Satisfies[T any]() (*ContractInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.Satisfies[T]()
	if !check.OK() {
		return nil, check
	}
	return p.fields.register(descriptor), check
}
