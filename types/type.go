package types

import (
	"reflect"

	typesEngine "github.com/qlustra/intel/engines/types"
)

type TypeInterrogator struct {
	*typesEngine.TypeInspector
	typ reflect.Type
}

func CaptureType[T any]() *TypeInterrogator {
	typ := reflect.TypeFor[T]()
	inspector := typesEngine.CaptureType(typ)
	typeInterrogator := &TypeInterrogator{
		TypeInspector: inspector,
		typ:           typ,
	}
	return typeInterrogator
}

func (i *TypeInterrogator) AsStruct() (*StructInterrogator, bool) {
	if !i.IsStructKind() {
		return nil, false
	}
	structInterrogator := captureStruct(i.TypeInspector)
	return structInterrogator, true
}

func (i *TypeInterrogator) AsInterface() (*InterfaceInterrogator, bool) {
	if !i.IsInterfaceKind() {
		return nil, false
	}
	interfaceInterrogator := captureInterface(i.TypeInspector)
	return interfaceInterrogator, true
}
