package inspect

import (
	"fmt"
	"reflect"

	"github.com/qlustra/intel/logic"
)

type TypeInspector struct {
	typ reflect.Type
}

func CaptureType(typ reflect.Type) *TypeInspector {
	typeInspector := &TypeInspector{
		typ: typ,
	}
	return typeInspector
}

func (i *TypeInspector) CapturedType() reflect.Type {
	return i.typ
}

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

func (i *TypeInspector) Equals(otherType reflect.Type) logic.Check {
	return IsType(i.typ, otherType)
}

const TypeIsNotPrimitiveCode logic.CheckCode = "type_not_primitive"

func (i *TypeInspector) IsPrimitive() logic.Check {
	if i.Is[string]().OK() || i.Is[int]().OK() || i.Is[bool]().OK() {
		return logic.Ok()
	}
	return logic.Nope(TypeIsNotPrimitiveCode, fmt.Sprintf("%v is not string, int, or bool", i.typ))
}

func (i *TypeInspector) Is[T any]() logic.Check {
	return Is[T](i.typ)
}

func (i *TypeInspector) Satisfies(otherType reflect.Type) logic.Check {
	return SatisfiesType(i.typ, otherType)
}

func (i *TypeInspector) SatisfiesType[O any]() logic.Check {
	return Satisfies[O](i.typ)
}

func (i *TypeInspector) AssignableTo(otherType reflect.Type) bool {
	return i.typ.AssignableTo(otherType)
}

func (i *TypeInspector) AssignableToType[O any]() bool {
	return i.AssignableTo(reflect.TypeFor[O]())
}

func (i *TypeInspector) ConvertibleTo(otherType reflect.Type) bool {
	return i.typ.ConvertibleTo(otherType)
}

func (i *TypeInspector) ConvertibleToType[O any]() bool {
	return i.ConvertibleTo(reflect.TypeFor[O]())
}
