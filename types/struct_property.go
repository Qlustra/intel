package types

import (
	structEngine "github.com/qlustra/intel/engines/struct"
)

type PropertyInterrogator struct {
	fieldInterrogatorMember
	*structEngine.FieldDescriptor
}

func captureProperty(descriptor *structEngine.FieldDescriptor) *PropertyInterrogator {
	interrogator := &PropertyInterrogator{
		FieldDescriptor: descriptor,
	}
	return interrogator
}
