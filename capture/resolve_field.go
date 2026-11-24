package capture

import (
	"reflect"

	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
)

// ResolveField captures a writable, accessible field in the supplied instance,
// not a copy. The instance must be a non-nil pointer whose element type is
// exactly the struct type captured by expectedRoot.
//
// index is the full Go field index path from that struct, including embedding
// steps for a promoted field. Non-nil intermediate pointers are followed;
// nil intermediate pointers cause failure and are not allocated. A pointer
// field at the end of the path is captured as a field slot, even if it is nil.
//
// A failed check returns no handle. ResolveField checks that the path can be
// used on this instance, not where the path came from. It does not modify the
// instance or the supplied index.
func ResolveField(instance any, expectedRoot *TypeInterrogator, index []int) (*ValueInterrogator, logic.Check) {
	if expectedRoot == nil || expectedRoot.TypeInspector == nil || expectedRoot.typ == nil ||
		expectedRoot.CapturedType() != expectedRoot.typ || expectedRoot.typ.Kind() != reflect.Struct {
		return nil, logic.Nope(ValueTypeMismatchCode, "expected root type is not valid struct metadata")
	}
	value := reflect.ValueOf(instance)
	if !value.IsValid() || value.Kind() != reflect.Pointer || value.IsNil() || value.Type().Elem() != expectedRoot.typ {
		return nil, logic.Nope(ValueTypeMismatchCode, "instance root type must be a non-nil pointer to the exact expected struct type")
	}
	if len(index) == 0 {
		return nil, logic.Nope(FieldPathInvalidCode, "empty field path")
	}
	typ := expectedRoot.typ
	var field reflect.StructField
	for n, step := range index {
		if typ.Kind() != reflect.Struct || step < 0 || step >= typ.NumField() {
			return nil, logic.Nope(FieldPathInvalidCode, "invalid field index path")
		}
		field = typ.Field(step)
		if n < len(index)-1 {
			typ = field.Type
			if typ.Kind() == reflect.Pointer {
				typ = typ.Elem()
			}
		}
	}
	field.Index = append([]int(nil), index...)
	root, check := Value(instance)
	if !check.OK() {
		return nil, check
	}
	structure, check := root.AsStruct()
	if !check.OK() {
		return nil, check
	}
	return structure.Field(inspect.CaptureFieldDescriptor(field))
}
