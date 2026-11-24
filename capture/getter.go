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

// ZeroReceiver resolves pointer layers and captures receiver type metadata;
// it does not capture a live instance. See [ReceiverInterrogator.Getter] and
// [GetterInterrogator.Read] to invoke a getter on a fresh zero receiver.
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

// Type returns the resolved receiver type, or nil if no receiver was captured.
func (r *ReceiverInterrogator) Type() *TypeInterrogator {
	if r == nil || r.typ == nil {
		return nil
	}
	return TypeReflection(r.typ.Type())
}

// Method finds a receiver method by name, returning a failed check if absent.
// The returned metadata can be checked with [ReceiverInterrogator.Getter].
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

// MethodInterrogator describes a receiver method, excluding the receiver from
// its input arity and indexed inputs. It does not invoke the method.
type MethodInterrogator struct {
	method *inspect.ReceiverMethodInspector
}

func (m *MethodInterrogator) Name() string { return m.method.Name() }
func (m *MethodInterrogator) ReceiverType() *TypeInterrogator {
	return TypeReflection(m.method.ReceiverType())
}
func (m *MethodInterrogator) InputArity() int  { return m.method.InputArity() }
func (m *MethodInterrogator) OutputArity() int { return m.method.OutputArity() }

// Input returns the type of the zero-based input or a failed check if out of range.
func (m *MethodInterrogator) Input(index int) (*TypeInterrogator, logic.Check) {
	if m == nil || index < 0 || index >= m.InputArity() {
		return nil, logic.Nope(MethodInputIndexInvalidCode, fmt.Sprint(index))
	}
	return TypeReflection(m.method.InputType(index)), logic.Ok()
}

// Result returns the type of the zero-based result or a failed check if out of range.
func (m *MethodInterrogator) Result(index int) (*TypeInterrogator, logic.Check) {
	if m == nil || index < 0 || index >= m.OutputArity() {
		return nil, logic.Nope(MethodResultIndexInvalidCode, fmt.Sprint(index))
	}
	return TypeReflection(m.method.ResultType(index)), logic.Ok()
}

// ResultAssignableTo reports whether the sole result is assignable to typ.
// It returns false for a nil type or a method with other than one result.
func (m *MethodInterrogator) ResultAssignableTo(typ reflect.Type) bool {
	return m != nil && m.OutputArity() == 1 && typ != nil && m.method.ResultType(0).AssignableTo(typ)
}

// Getter checks that method belongs to this receiver and has no inputs and
// exactly one result. It does not check the result type or invoke the method;
// use [MethodInterrogator.ResultAssignableTo] when a result type is required.
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

// GetterInterrogator holds checked getter metadata for [GetterInterrogator.Read].
type GetterInterrogator struct {
	receiverType reflect.Type
	method       *MethodInterrogator
}

// Read invokes the getter on a newly allocated zero receiver, not a live
// captured instance. A panic becomes a failed check. For live invocation see
// [StructValueInterrogator.Invoke].
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

// ResultInterrogator holds the result of a successful [GetterInterrogator.Read].
type ResultInterrogator struct{ value reflect.Value }

// Type describes the result's type.
func (r *ResultInterrogator) Type() *TypeInterrogator { return TypeReflection(r.value.Type()) }

// Get returns the result as an interface value. Use [Assign] for a checked typed copy.
func (r *ResultInterrogator) Get() any { return r.value.Interface() }
