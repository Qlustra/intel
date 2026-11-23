package capture

import (
	"fmt"
	"reflect"

	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
)

const (
	ZeroReceiverTypeInvalidCode  logic.CheckCode = "zero_receiver_type_invalid"
	MethodNotFoundCode           logic.CheckCode = "method_not_found"
	GetterMethodMismatchCode     logic.CheckCode = "getter_method_mismatch"
	GetterSignatureInvalidCode   logic.CheckCode = "getter_signature_invalid"
	MethodInputIndexInvalidCode  logic.CheckCode = "method_input_index_invalid"
	MethodResultIndexInvalidCode logic.CheckCode = "method_result_index_invalid"
)

// ZeroReceiver resolves pointer layers and captures receiver type metadata.
func (i *TypeInterrogator) ZeroReceiver() (*ReceiverInterrogator, logic.Check) {
	if i == nil || i.TypeInspector == nil || i.CapturedType() == nil {
		return nil, logic.Nope(ZeroReceiverTypeInvalidCode, "no type was captured")
	}
	inspector, ok := inspect.InspectReceiver(i.CapturedType())
	if !ok {
		return nil, logic.Nope(ZeroReceiverTypeInvalidCode, fmt.Sprintf("%v cannot resolve to a receiver", i.CapturedType()))
	}
	return &ReceiverInterrogator{typ: inspector}, logic.Ok()
}

// ReceiverInterrogator captures receiver method metadata, not an instance.
type ReceiverInterrogator struct{ typ *inspect.ReceiverInspector }

func (r *ReceiverInterrogator) Type() *TypeInterrogator {
	if r == nil || r.typ == nil {
		return nil
	}
	return TypeReflection(r.typ.Type())
}

func (r *ReceiverInterrogator) Method(name string) (*MethodInterrogator, logic.Check) {
	if r == nil || r.typ == nil {
		return nil, logic.Nope(ZeroReceiverTypeInvalidCode, "no receiver type was captured")
	}
	method, ok := r.typ.Method(name)
	if !ok {
		return nil, logic.Nope(MethodNotFoundCode, name)
	}
	return &MethodInterrogator{method: method}, logic.Ok()
}

type MethodInterrogator struct {
	method *inspect.ReceiverMethodInspector
}

func (m *MethodInterrogator) Name() string { return m.method.Name() }
func (m *MethodInterrogator) ReceiverType() *TypeInterrogator {
	return TypeReflection(m.method.ReceiverType())
}
func (m *MethodInterrogator) InputArity() int  { return m.method.InputArity() }
func (m *MethodInterrogator) OutputArity() int { return m.method.OutputArity() }
func (m *MethodInterrogator) Input(index int) (*TypeInterrogator, logic.Check) {
	if m == nil || index < 0 || index >= m.InputArity() {
		return nil, logic.Nope(MethodInputIndexInvalidCode, fmt.Sprint(index))
	}
	return TypeReflection(m.method.InputType(index)), logic.Ok()
}
func (m *MethodInterrogator) Result(index int) (*TypeInterrogator, logic.Check) {
	if m == nil || index < 0 || index >= m.OutputArity() {
		return nil, logic.Nope(MethodResultIndexInvalidCode, fmt.Sprint(index))
	}
	return TypeReflection(m.method.ResultType(index)), logic.Ok()
}
func (m *MethodInterrogator) ResultAssignableTo(typ reflect.Type) bool {
	return m != nil && m.OutputArity() == 1 && typ != nil && m.method.ResultType(0).AssignableTo(typ)
}

func (r *ReceiverInterrogator) Getter(method *MethodInterrogator) (*GetterInterrogator, logic.Check) {
	if r == nil || r.typ == nil {
		return nil, logic.Nope(ZeroReceiverTypeInvalidCode, "no receiver type was captured")
	}
	if method == nil || method.method == nil || method.method.ReceiverType() != r.typ.Type() {
		return nil, logic.Nope(GetterMethodMismatchCode, "method does not belong to receiver")
	}
	if method.InputArity() != 0 || method.OutputArity() != 1 {
		return nil, logic.Nope(GetterSignatureInvalidCode, "getter must take no inputs and return one result")
	}
	return &GetterInterrogator{receiverType: r.typ.Type(), method: method}, logic.Ok()
}

type GetterInterrogator struct {
	receiverType reflect.Type
	method       *MethodInterrogator
}

func (g *GetterInterrogator) Read() (result *ResultInterrogator, check logic.Check) {
	if g == nil || g.receiverType == nil || g.method == nil || g.method.method == nil {
		return nil, logic.Nope(ZeroReceiverTypeInvalidCode, "no receiver metadata was captured")
	}
	defer func() {
		if failure := recover(); failure != nil {
			result = nil
			check = logic.Nope(InvocationPanickedCode, fmt.Sprintf("%v.%s() panicked: %v", g.receiverType, g.method.Name(), failure))
		}
	}()
	value := reflect.New(g.receiverType.Elem())
	results := g.method.method.Function().Call([]reflect.Value{value})
	return &ResultInterrogator{value: results[0]}, logic.Ok()
}

type ResultInterrogator struct{ value reflect.Value }

func (r *ResultInterrogator) Type() *TypeInterrogator { return TypeReflection(r.value.Type()) }
func (r *ResultInterrogator) Get() any                { return r.value.Interface() }
