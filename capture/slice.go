package capture

import (
	"reflect"

	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
)

// SliceValueInterrogator operates on the captured slice storage. Obtain one
// with [ValueInterrogator.AsSlice]; element captures may become stale after
// slice growth, as described by [SliceValueInterrogator.AppendZero].
type SliceValueInterrogator struct {
	// val holds *T, while typ describes the slice T.
	val reflect.Value
	typ *SliceInterrogator
}

// AppendZero appends a fresh zero element and captures its actual address.
// As with Go append, later appends may reallocate the backing array; existing
// element captures are not stable handles to the current slice in that case.
// Fill an element before appending the next when its current identity matters.
func (i *SliceValueInterrogator) AppendZero() (*ValueInterrogator, logic.Check) {
	if i == nil || !i.val.IsValid() || i.val.Kind() != reflect.Pointer || i.val.IsNil() ||
		!i.val.CanInterface() || i.typ == nil || i.typ.TypeInspector == nil {
		return nil, logic.Nope(ValueIsInvalidCode, "no writable slice value was captured")
	}
	slice := i.val.Elem()
	if !slice.CanSet() || slice.Kind() != reflect.Slice {
		return nil, logic.Nope(ValueIsInvalidCode, "captured value is not a writable slice")
	}
	if i.typ.CapturedType() != slice.Type() {
		return nil, logic.Nope(ValueTypeMismatchCode, "captured slice type does not match value")
	}
	slice.Set(reflect.Append(slice, reflect.Zero(slice.Type().Elem())))
	element := slice.Index(slice.Len() - 1)
	return &ValueInterrogator{val: element.Addr(), typ: TypeReflection(element.Type())}, logic.Ok()
}

// Elements captures each current slice element at its actual address, in index
// order. The slice and its elements are not modified. Later growth can
// reallocate the backing array; reacquire handles after appending.
func (i *SliceValueInterrogator) Elements() ([]*ValueInterrogator, logic.Check) {
	if i == nil || !i.val.IsValid() || i.val.Kind() != reflect.Pointer || i.val.IsNil() ||
		!i.val.CanInterface() || i.typ == nil || i.typ.TypeInspector == nil {
		return nil, logic.Nope(ValueIsInvalidCode, "no writable slice value was captured")
	}
	slice := i.val.Elem()
	if !slice.CanSet() || slice.Kind() != reflect.Slice {
		return nil, logic.Nope(ValueIsInvalidCode, "captured value is not a writable slice")
	}
	if i.typ.CapturedType() != slice.Type() {
		return nil, logic.Nope(ValueTypeMismatchCode, "captured slice type does not match value")
	}
	elements := make([]*ValueInterrogator, slice.Len())
	for index := range elements {
		element := slice.Index(index)
		if !element.CanAddr() || !element.Addr().CanInterface() {
			return nil, logic.Nope(ValueIsInvalidCode, "slice element is not accessible")
		}
		elements[index] = &ValueInterrogator{val: element.Addr(), typ: TypeReflection(element.Type())}
	}
	return elements, logic.Ok()
}

// SliceInterrogator describes a slice type rather than a slice value.
type SliceInterrogator struct {
	*inspect.TypeInspector
	element *TypeInterrogator
}

func captureSlice(inspector, element *inspect.TypeInspector) *SliceInterrogator {
	interrogator := &SliceInterrogator{
		TypeInspector: inspector,
		element:       TypeInspector(element),
	}
	return interrogator
}

// Element describes the exact element type, which may itself be a pointer.
func (i *SliceInterrogator) Element() *TypeInterrogator {
	return i.element
}

// AsStructSlice narrows a slice whose element has struct kind. It does not
// dereference pointer elements; a different element kind fails the check.
func (i *SliceInterrogator) AsStructSlice() (*StructSliceInterrogator, logic.Check) {
	element, check := i.Element().AsStruct()
	if !check.OK() {
		return nil, check
	}
	interrogator := captureStructSlice(i, element)
	return interrogator, check
}
