package capture

import (
	"fmt"
	"reflect"

	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
)

// TypeInterrogator describes a captured Go type. Its narrowing methods return
// specialized views without changing or dereferencing that type.
type TypeInterrogator struct {
	*inspect.TypeInspector
	typ reflect.Type
}

// Type captures the exact Go type T without creating a value of that type.
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

// TypeReflection captures a reflect.Type for interrogation.
func TypeReflection(typ reflect.Type) *TypeInterrogator {
	inspector := inspect.CaptureType(typ)
	typeInterrogator := &TypeInterrogator{
		TypeInspector: inspector,
		typ:           typ,
	}
	return typeInterrogator
}

// TypeInspector wraps an existing type inspector.
func TypeInspector(inspector *inspect.TypeInspector) *TypeInterrogator {
	typeInterrogator := &TypeInterrogator{
		TypeInspector: inspector,
		typ:           inspector.CapturedType(),
	}
	return typeInterrogator
}

// AsStruct narrows the captured type to a struct. It does not dereference
// pointer types. A non-struct type produces a failed check and no handle.
func (i *TypeInterrogator) AsStruct() (*StructInterrogator, logic.Check) {
	check := i.IsStructKind()
	if !check.OK() {
		return nil, check
	}
	structInterrogator := captureStruct(i.TypeInspector)
	return structInterrogator, check
}

// AsInterface narrows an interface type without dereferencing pointers. A
// different kind returns a failed check and no handle.
func (i *TypeInterrogator) AsInterface() (*InterfaceInterrogator, logic.Check) {
	check := i.IsInterfaceKind()
	if !check.OK() {
		return nil, check
	}
	interfaceInterrogator := captureInterface(i.TypeInspector)
	return interfaceInterrogator, check
}

// AsSlice narrows a slice type without dereferencing pointers. A different
// kind returns a failed check and no handle; arrays are not slices.
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
