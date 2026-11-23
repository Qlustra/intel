package inspect

import (
	"go/token"
	"reflect"
)

// IsExportedIdentifier reports whether name is a valid Go identifier that is exported.
func IsExportedIdentifier(name string) bool {
	return token.IsIdentifier(name) && token.IsExported(name)
}

// InterfaceInspector inspects the method set of an interface type.
type InterfaceInspector struct {
	typ reflect.Type
}

func InspectInterface(typ reflect.Type) *InterfaceInspector {
	return &InterfaceInspector{typ: typ}
}

func (i *InterfaceInspector) IsEmpty() bool {
	return i.typ.NumMethod() == 0
}

type InterfaceMethodInspector struct {
	method reflect.Method
}

func InspectInterfaceMethod(typ reflect.Type, name string) (*InterfaceMethodInspector, bool) {
	method, ok := typ.MethodByName(name)
	if !ok {
		return nil, false
	}
	return &InterfaceMethodInspector{method: method}, true
}

func (i *InterfaceMethodInspector) Name() string {
	return i.method.Name
}

func (i *InterfaceMethodInspector) IsExported() bool {
	return i.method.PkgPath == ""
}

func (i *InterfaceMethodInspector) InputArity() int {
	return i.method.Type.NumIn()
}

func (i *InterfaceMethodInspector) OutputArity() int {
	return i.method.Type.NumOut()
}

func (i *InterfaceMethodInspector) ResultAssignableTo(typ reflect.Type) bool {
	return i.OutputArity() == 1 && typ != nil && i.method.Type.Out(0).AssignableTo(typ)
}

// ReceiverInspector describes the pointer receiver type and its methods without
// retaining a receiver value.
type ReceiverInspector struct {
	typ reflect.Type
}

// InspectReceiver resolves pointer layers to a concrete receiver type and
// describes its pointer method set.
func InspectReceiver(typ reflect.Type) (*ReceiverInspector, bool) {
	if typ == nil {
		return nil, false
	}
	seen := make(map[reflect.Type]struct{})
	// TODO: extract into a pointer chain resolver helper later
	for typ.Kind() == reflect.Pointer {
		if _, exists := seen[typ]; exists {
			return nil, false
		}
		seen[typ] = struct{}{}
		typ = typ.Elem()
	}
	if typ.Kind() == reflect.Interface {
		return nil, false
	}
	return &ReceiverInspector{typ: reflect.PointerTo(typ)}, true
}

func (i *ReceiverInspector) Type() reflect.Type { return i.typ }

func (i *ReceiverInspector) Method(name string) (*ReceiverMethodInspector, bool) {
	method, ok := i.typ.MethodByName(name)
	if !ok || method.PkgPath != "" {
		return nil, false
	}
	return &ReceiverMethodInspector{receiverType: i.typ, method: method}, true
}

// ReceiverMethodInspector exposes pure metadata for a pointer receiver method.
type ReceiverMethodInspector struct {
	receiverType reflect.Type
	method       reflect.Method
}

func (i *ReceiverMethodInspector) Name() string               { return i.method.Name }
func (i *ReceiverMethodInspector) ReceiverType() reflect.Type { return i.receiverType }
func (i *ReceiverMethodInspector) InputArity() int            { return i.method.Type.NumIn() - 1 }
func (i *ReceiverMethodInspector) OutputArity() int           { return i.method.Type.NumOut() }
func (i *ReceiverMethodInspector) InputType(index int) reflect.Type {
	return i.method.Type.In(index + 1)
}
func (i *ReceiverMethodInspector) ResultType(index int) reflect.Type {
	return i.method.Type.Out(index)
}
func (i *ReceiverMethodInspector) Function() reflect.Value { return i.method.Func }
