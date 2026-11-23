package capture

import (
	"reflect"

	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
)

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
// order. The slice and its elements are not modified.
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

func (i *SliceInterrogator) Element() *TypeInterrogator {
	return i.element
}

func (i *SliceInterrogator) AsStructSlice() (*StructSliceInterrogator, logic.Check) {
	element, check := i.Element().AsStruct()
	if !check.OK() {
		return nil, check
	}
	interrogator := captureStructSlice(i, element)
	return interrogator, check
}
