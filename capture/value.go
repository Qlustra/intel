package capture

import (
	"fmt"
	"reflect"

	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
)

const (
	ValueIsInvalidCode          logic.CheckCode = "value_is_invalid"
	ValueIsNilCode              logic.CheckCode = "value_is_nil"
	ValueTypeMismatchCode       logic.CheckCode = "value_type_mismatch"
	InvocationCallbackIsNilCode logic.CheckCode = "invocation_callback_is_nil"
	InvocationPanickedCode      logic.CheckCode = "invocation_panicked"
)

// ValueInterrogator holds addressable captured storage. See [Value] for the
// distinction between capturing a pointer and copying a non-pointer.
type ValueInterrogator struct {
	// val holds *T, while typ describes T.
	val reflect.Value
	typ *TypeInterrogator
}

// Value captures the storage addressed by a supplied non-nil pointer, or a
// shallow copy of a non-pointer value in new storage. The copy can still share
// referenced data such as a slice's backing array. Untyped nil and nil pointers
// return a failed check and no handle.
func Value(value any) (*ValueInterrogator, logic.Check) {
	val := reflect.ValueOf(value)
	if !val.IsValid() {
		return nil, logic.Nope(ValueIsInvalidCode, "no value was supplied")
	}
	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return nil, logic.Nope(ValueIsNilCode, fmt.Sprintf("%v is nil", val.Type()))
		}
	} else {
		instance := reflect.New(val.Type())
		instance.Elem().Set(val)
		val = instance
	}

	interrogator := &ValueInterrogator{
		val: val,
		typ: TypeReflection(val.Type().Elem()),
	}
	return interrogator, logic.Ok()
}

// AsStruct checks the captured type and shares the existing instance.
func (v *ValueInterrogator) AsStruct() (*StructValueInterrogator, logic.Check) {
	if v == nil || !v.val.IsValid() {
		return nil, logic.Nope(ValueIsInvalidCode, "no value was captured")
	}

	structInterrogator, check := v.typ.AsStruct()
	if !check.OK() {
		return nil, check
	}
	structValueInterrogator := &StructValueInterrogator{
		val: v.val,
		typ: structInterrogator,
	}
	return structValueInterrogator, logic.Ok()
}

// AsSlice checks the captured type and retains the captured slice's address.
func (v *ValueInterrogator) AsSlice() (*SliceValueInterrogator, logic.Check) {
	if v == nil || !v.val.IsValid() || v.val.Kind() != reflect.Pointer || v.val.IsNil() ||
		!v.val.CanInterface() || !v.val.Elem().CanSet() || v.typ == nil || v.typ.TypeInspector == nil {
		return nil, logic.Nope(ValueIsInvalidCode, "no writable slice value was captured")
	}
	if v.typ.CapturedType() != v.val.Type().Elem() || v.typ.typ != v.val.Type().Elem() {
		return nil, logic.Nope(ValueTypeMismatchCode, "captured slice type does not match value")
	}
	typ, check := v.typ.AsSlice()
	if !check.OK() {
		return nil, check
	}
	return &SliceValueInterrogator{val: v.val, typ: typ}, logic.Ok()
}

// StructValueInterrogator operates on captured struct storage. Obtain one via
// [ValueInterrogator.AsStruct] or [StructInterrogator.CreateValue].
type StructValueInterrogator struct {
	// val holds *S, while typ describes S.
	val reflect.Value
	typ *StructInterrogator
}

// Address returns the pointer to the captured storage, not a copy of its
// contents. A pointer supplied to [Value] is returned with its original Go type.
// The address of a captured field of type T is *T, so a pointer-valued field
// has a pointer-to-pointer address. If no accessible address is available,
// Address returns nil and a failed check.
func (v *ValueInterrogator) Address() (any, logic.Check) {
	if v == nil || !v.val.IsValid() || v.val.Kind() != reflect.Pointer || v.val.IsNil() || !v.val.CanInterface() {
		return nil, logic.Nope(ValueIsInvalidCode, "captured value has no accessible address")
	}
	return v.val.Interface(), logic.Ok()
}

// Field captures the writable, accessible field at the descriptor's index
// path in the captured struct value. The path must resolve and the terminal
// field's name and type must match. Field does not check the descriptor's
// origin or tag identity. Intermediate pointers must be non-nil; none are
// allocated. On failure, Field returns no handle.
func (i *StructValueInterrogator) Field(field *inspect.FieldDescriptor) (*ValueInterrogator, logic.Check) {
	if i == nil || !i.val.IsValid() || i.val.Kind() != reflect.Pointer || i.val.IsNil() || i.typ == nil {
		return nil, logic.Nope(ValueIsInvalidCode, "no struct value was captured")
	}
	if i.typ.TypeInspector == nil || i.val.Type().Elem() != i.typ.CapturedType() {
		return nil, logic.Nope(ValueTypeMismatchCode, "captured struct type does not match value")
	}
	path, _, check := i.typ.fieldPath(field)
	if !check.OK() {
		return nil, check
	}
	value := i.val.Elem()
	for n, index := range path {
		value = value.Field(index)
		if n < len(path)-1 && value.Kind() == reflect.Pointer {
			if value.IsNil() {
				return nil, logic.Nope(ValueIsNilCode, "field path crosses a nil pointer")
			}
			value = value.Elem()
		}
	}
	if !value.CanSet() || !value.CanAddr() || !value.Addr().CanInterface() {
		return nil, logic.Nope(ValueIsInvalidCode, "field is not writable and accessible")
	}
	return &ValueInterrogator{val: value.Addr(), typ: TypeReflection(value.Type())}, logic.Ok()
}

// Invoke checks the instance against API, then calls invoke with that typed
// instance. methodName labels diagnostics; the callback selects the method.
// Only callback panics are recovered. A callback failure check is returned as-is.
func (i *StructValueInterrogator) Invoke[API, RESULT any](
	methodName string,
	invoke func(API) (RESULT, logic.Check),
) (result RESULT, check logic.Check) {
	if i == nil || !i.val.IsValid() {
		return result, logic.Nope(ValueIsInvalidCode, "no struct value was captured")
	}
	receiverType := i.val.Type()
	if invoke == nil {
		return result, logic.Nope(InvocationCallbackIsNilCode,
			fmt.Sprintf("%v.%s(): callback is nil", receiverType, methodName))
	}

	instance, ok := i.val.Interface().(API)
	if !ok {
		return result, logic.Nope(ValueTypeMismatchCode,
			fmt.Sprintf("%v.%s(): instance does not satisfy %v",
				receiverType, methodName, reflect.TypeFor[API]()))
	}

	defer func() {
		if failure := recover(); failure != nil {
			check = logic.Nope(InvocationPanickedCode,
				fmt.Sprintf("%v.%s(): callback panicked: %v", receiverType, methodName, failure))
		}
	}()
	return invoke(instance)
}
