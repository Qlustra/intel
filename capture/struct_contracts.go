package capture

import (
	"iter"

	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
	"github.com/qlustra/intel/options"
)

type Contracts struct {
	fieldsInspector *inspect.FieldsInspector
	fields          *FieldRegistry[ContractInterrogator]
}

func (p *Contracts) AllAddressable(ordering ...options.DepthOrdering) iter.Seq2[string, *ContractInterrogator] {
	return func(yield func(string, *ContractInterrogator) bool) {
		for name, descriptor := range p.fieldsInspector.AllAddressable(ordering...) {
			if !yield(name, p.fields.register(descriptor)) {
				return
			}
		}
	}
}

func (p *Contracts) Lookup(name string) (*ContractInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.LookupNamed(name)
	if !check.OK() {
		return nil, check
	}
	return p.fields.register(descriptor), check
}

func (p *Contracts) Declared(name string) (*ContractInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.Declared(name)
	if !check.OK() {
		return nil, check
	}
	return p.fields.register(descriptor), check
}

func (p *Contracts) Embeds[T any]() (*ContractInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.Embeds[T]()
	if !check.OK() {
		return nil, check
	}
	return p.fields.register(descriptor), check
}

func (p *Contracts) Satisfies[T any]() (*ContractInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.Satisfies[T]()
	if !check.OK() {
		return nil, check
	}
	return p.fields.register(descriptor), check
}
