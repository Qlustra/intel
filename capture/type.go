package capture

import (
	"fmt"
	"reflect"

	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
)

type TypeInterrogator struct {
	*inspect.TypeInspector
	typ reflect.Type
}

func Type[T any]() *TypeInterrogator {
	return TypeReflection(reflect.TypeFor[T]())
}

// TypeOf captures the dynamic type of value. Untyped nil has no dynamic type.
func TypeOf(value any) (*TypeInterrogator, logic.Check) {
	typ := reflect.TypeOf(value)
	if typ == nil {
		return nil, logic.Nope(ValueIsInvalidCode, "no type was supplied")
	}
	return TypeReflection(typ), logic.Ok()
}

func TypeReflection(typ reflect.Type) *TypeInterrogator {
	inspector := inspect.CaptureType(typ)
	typeInterrogator := &TypeInterrogator{
		TypeInspector: inspector,
		typ:           typ,
	}
	return typeInterrogator
}

func TypeInspector(inspector *inspect.TypeInspector) *TypeInterrogator {
	typeInterrogator := &TypeInterrogator{
		TypeInspector: inspector,
		typ:           inspector.CapturedType(),
	}
	return typeInterrogator
}

func (i *TypeInterrogator) AsStruct() (*StructInterrogator, logic.Check) {
	check := i.IsStructKind()
	if !check.OK() {
		return nil, check
	}
	structInterrogator := captureStruct(i.TypeInspector)
	return structInterrogator, check
}

func (i *TypeInterrogator) AsInterface() (*InterfaceInterrogator, logic.Check) {
	check := i.IsInterfaceKind()
	if !check.OK() {
		return nil, check
	}
	interfaceInterrogator := captureInterface(i.TypeInspector)
	return interfaceInterrogator, check
}

func (i *TypeInterrogator) AsSlice() (*SliceInterrogator, logic.Check) {
	element, check := i.SliceElement()
	if !check.OK() {
		return nil, check
	}
	sliceInterrogator := captureSlice(i.TypeInspector, element)
	return sliceInterrogator, check
}

func (i *TypeInterrogator) String() string {
	return fmt.Sprintf("type %s %s", i.TypeName(), i.Kind().String())
}
