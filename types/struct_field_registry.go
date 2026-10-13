package types

import structEngine "github.com/qlustra/intel/engines/struct"

type fieldInterrogatorMember interface {
	interrogate()
}

type fieldCapturer[T fieldInterrogatorMember] func(descriptor *structEngine.FieldDescriptor) *T

type FieldRegistry[T fieldInterrogatorMember] struct {
	fields  map[string]*T
	capture fieldCapturer[T]
}

func createFieldRegistry[T fieldInterrogatorMember](capture fieldCapturer[T]) *FieldRegistry[T] {
	return &FieldRegistry[T]{
		capture: capture,
	}
}

func (r *FieldRegistry[T]) register(descriptor *structEngine.FieldDescriptor) *T {
	if r.fields == nil {
		r.fields = make(map[string]*T)
	}
	name := descriptor.FieldName()

	existing, found := r.fields[name]
	if found {
		return existing
	}
	interrogator := r.capture(descriptor)
	r.fields[name] = interrogator
	return interrogator
}
