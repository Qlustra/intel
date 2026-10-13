package types

import "reflect"

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

func (i *TypeInspector) TypeName() string {
	return i.typ.Name()
}

func (i *TypeInspector) PkgPath() string {
	return i.typ.PkgPath()
}

func (i *TypeInspector) Equals(otherType reflect.Type) bool {
	return i.typ == otherType
}

func (i *TypeInspector) IsPrimitive() bool {
	if !i.IsPrimitiveKind() {
		return false
	}
	return i.Is[string]() || i.Is[int]() || i.Is[bool]()
}

func (i *TypeInspector) Is[T any]() bool {
	return Is[T](i.typ)
}

func (i *TypeInspector) Satisfies(otherType reflect.Type) bool {
	return SatisfiesType(i.typ, otherType)
}

func (i *TypeInspector) SatisfiesType[O any]() bool {
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

func Is[T any](typ reflect.Type) bool {
	return typ == reflect.TypeFor[T]()
}

func SatisfiesType(typ, otherType reflect.Type) bool {
	return otherType != nil &&
		IsInterfaceKind(otherType) &&
		typ.Implements(otherType)
}

func Satisfies[T any](typ reflect.Type) bool {
	return SatisfiesType(typ, reflect.TypeFor[T]())
}
