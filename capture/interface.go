package capture

import (
	"reflect"

	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
)

type InterfaceInterrogator struct {
	*inspect.TypeInspector
}

func captureInterface(inspector *inspect.TypeInspector) *InterfaceInterrogator {
	return &InterfaceInterrogator{
		TypeInspector: inspector,
	}
}

// IsEmpty reports whether the captured interface declares no methods.
func (i *InterfaceInterrogator) IsEmpty() bool {
	return inspect.InspectInterface(i.CapturedType()).IsEmpty()
}

const InterfaceMethodNotFoundCode logic.CheckCode = "interface_method_not_found"

// Method returns an exported method declared by the interface.
func (i *InterfaceInterrogator) Method(name string) (*InterfaceMethodInterrogator, logic.Check) {
	method, ok := inspect.InspectInterfaceMethod(i.CapturedType(), name)
	if !ok || !method.IsExported() {
		return nil, logic.Nope(InterfaceMethodNotFoundCode, name)
	}
	return &InterfaceMethodInterrogator{method: method}, logic.Ok()
}

type InterfaceMethodInterrogator struct {
	method *inspect.InterfaceMethodInspector
}

func (i *InterfaceMethodInterrogator) Name() string {
	return i.method.Name()
}

func (i *InterfaceMethodInterrogator) InputArity() int {
	return i.method.InputArity()
}

func (i *InterfaceMethodInterrogator) OutputArity() int {
	return i.method.OutputArity()
}

// ResultAssignableTo reports whether the method's sole result is assignable to typ.
func (i *InterfaceMethodInterrogator) ResultAssignableTo(typ reflect.Type) bool {
	return i.method.ResultAssignableTo(typ)
}
