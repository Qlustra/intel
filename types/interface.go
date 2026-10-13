package types

import (
	"reflect"

	structEngine "github.com/qlustra/intel/engines/struct"
	typesEngine "github.com/qlustra/intel/engines/types"
)

type InterfaceInterrogator struct {
	fieldInterrogatorMember
	*typesEngine.TypeInspector
	typ reflect.Type
}

func captureInterface(inspector *typesEngine.TypeInspector) *InterfaceInterrogator {
	interrogator := &InterfaceInterrogator{
		TypeInspector: inspector,
		typ:           inspector.CapturedType(),
	}
	return interrogator
}

func captureMembership(descriptor *structEngine.FieldDescriptor) *InterfaceInterrogator {
	return captureInterface(descriptor.FieldTypeInspector())
}

//func (i *InterfaceInterrogator) Field(name string) (*FieldInterrogator, bool) {
//	field, found := i.typ.FieldByName(name)
//	if !found {
//		return nil, false
//	}
//	interrogator := CaptureField(field)
//	return interrogator, true
//}

//func (i *InterfaceInterrogator) Embedding(name string) (*EmbeddingInterrogator, bool) {
//	field, found := i.Field(name)
//	if !found {
//		return nil, false
//	}
//	return field.AsEmbedding()
//}
