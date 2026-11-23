package capture

import (
	"iter"

	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
	"github.com/qlustra/intel/options"
)

type Extensions struct {
	fieldsInspector *inspect.FieldsInspector
	fields          *FieldRegistry[ExtensionInterrogator]
}

func (p *Extensions) AllAddressable(ordering ...options.DepthOrdering) iter.Seq2[string, *ExtensionInterrogator] {
	return func(yield func(string, *ExtensionInterrogator) bool) {
		for name, descriptor := range p.fieldsInspector.AllAddressable(ordering...) {
			if !yield(name, p.fields.register(descriptor)) {
				return
			}
		}
	}
}

func (p *Extensions) Lookup(name string) (*ExtensionInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.LookupNamed(name)
	if !check.OK() {
		return nil, check
	}
	return p.fields.register(descriptor), check
}

func (p *Extensions) Declared(name string) (*ExtensionInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.Declared(name)
	if !check.OK() {
		return nil, check
	}
	return p.fields.register(descriptor), check
}

func (p *Extensions) Embeds[T any]() (*ExtensionInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.Embeds[T]()
	if !check.OK() {
		return nil, check
	}
	return p.fields.register(descriptor), check
}

func (p *Extensions) Satisfies[T any]() (*ExtensionInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.Satisfies[T]()
	if !check.OK() {
		return nil, check
	}
	return p.fields.register(descriptor), check
}

func (p *Extensions) DeclaredType[T any]() (*ExtensionInterrogator, logic.Check) {
	descriptor, check := p.fieldsInspector.DeclaredType[T]()
	if !check.OK() {
		return nil, check
	}
	return p.fields.register(descriptor), check
}
