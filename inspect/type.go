package inspect

import (
	"fmt"
	"reflect"

	"github.com/qlustra/intel/logic"
)

type TypeInspector struct {
	typ reflect.Type
}

// CaptureType creates an inspector for typ. A nil type is retained and checked
// operations report failure where applicable; methods that directly query the
// reflected type require a non-nil type.
func CaptureType(typ reflect.Type) *TypeInspector {
	typeInspector := &TypeInspector{
		typ: typ,
	}
	return typeInspector
}

// CapturedType returns the reflected type held by the inspector.
func (i *TypeInspector) CapturedType() reflect.Type {
	return i.typ
}

// SliceElement returns an inspector for the immediate element type when the
// captured type is a slice. The result must only be used when the check is OK.
func (i *TypeInspector) SliceElement() (*TypeInspector, logic.Check) {
	check := i.IsSliceKind()
	if !check.OK() {
		return nil, check
	}
	return CaptureType(i.typ.Elem()), check
}

func (i *TypeInspector) TypeName() string {
	return i.typ.Name()
}

func (i *TypeInspector) PkgPath() string {
	return i.typ.PkgPath()
}

// Equals reports whether the captured type is exactly otherType.
func (i *TypeInspector) Equals(otherType reflect.Type) logic.Check {
	return IsType(i.typ, otherType)
}

const TypeIsNotPrimitiveCode logic.CheckCode = "type_not_primitive"

// IsPrimitive accepts exactly string, int, or bool; other types fail the check.
func (i *TypeInspector) IsPrimitive() logic.Check {
	if i.Is[string]().OK() || i.Is[int]().OK() || i.Is[bool]().OK() {
		return logic.Ok()
	}
	return logic.Nope(TypeIsNotPrimitiveCode, fmt.Sprintf("%v is not string, int, or bool", i.typ))
}

// Is reports whether the captured type is exactly T, not merely convertible
// or assignable to T.
func (i *TypeInspector) Is[T any]() logic.Check {
	return Is[T](i.typ)
}

// Satisfies reports whether the captured type's method set implements the
// interface type otherType. It does not substitute a pointer type or call a
// method; otherType must be an interface.
func (i *TypeInspector) Satisfies(otherType reflect.Type) logic.Check {
	return SatisfiesType(i.typ, otherType)
}

// SatisfiesType reports whether the captured type implements interface O.
func (i *TypeInspector) SatisfiesType[O any]() logic.Check {
	return Satisfies[O](i.typ)
}

// AssignableTo reports Go assignability to otherType.
func (i *TypeInspector) AssignableTo(otherType reflect.Type) bool {
	return i.typ.AssignableTo(otherType)
}

// AssignableToType reports Go assignability to O.
func (i *TypeInspector) AssignableToType[O any]() bool {
	return i.AssignableTo(reflect.TypeFor[O]())
}

// ConvertibleTo reports Go convertibility to otherType; it does not perform a
// conversion.
func (i *TypeInspector) ConvertibleTo(otherType reflect.Type) bool {
	return i.typ.ConvertibleTo(otherType)
}

// ConvertibleToType reports Go convertibility to O; it does not perform a
// conversion.
func (i *TypeInspector) ConvertibleToType[O any]() bool {
	return i.ConvertibleTo(reflect.TypeFor[O]())
}
