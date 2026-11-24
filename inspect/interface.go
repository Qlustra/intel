package inspect

import (
	"go/token"
	"reflect"
)

// IsExportedIdentifier reports whether name is a valid, exported Go
// identifier.
func IsExportedIdentifier(name string) bool {
	return token.IsIdentifier(name) && token.IsExported(name)
}

// InterfaceInspector inspects the method set of an interface type.
type InterfaceInspector struct {
	typ reflect.Type
}

// InspectInterface creates an inspector for typ. typ must be a non-nil
// interface type for its method-set queries.
func InspectInterface(typ reflect.Type) *InterfaceInspector {
	return &InterfaceInspector{typ: typ}
}

// IsEmpty reports whether the interface method set has no methods.
func (i *InterfaceInspector) IsEmpty() bool {
	return i.typ.NumMethod() == 0
}

type InterfaceMethodInspector struct {
	method reflect.Method
}

// InspectInterfaceMethod looks up name on an interface type. It returns false
// when no such method exists; typ must be a non-nil interface type.
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

// ResultAssignableTo reports whether the method has exactly one result and
// that result is assignable to typ. A nil typ returns false.
func (i *InterfaceMethodInspector) ResultAssignableTo(typ reflect.Type) bool {
	return i.OutputArity() == 1 && typ != nil && i.method.Type.Out(0).AssignableTo(typ)
}

// ReceiverInspector describes the pointer receiver type and its methods without
// retaining a receiver value.
type ReceiverInspector struct {
	typ reflect.Type
}

// InspectReceiver resolves any pointer layers in typ and describes the method
// set of a pointer to the resulting concrete type. Nil and interface types
// cannot be inspected and return false.
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

// Method looks up an exported method in the described pointer receiver's
// method set.
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
