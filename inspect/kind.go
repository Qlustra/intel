package inspect

import (
	"reflect"

	"github.com/qlustra/intel/logic"
)

func (i *TypeInspector) Kind() reflect.Kind {
	return i.typ.Kind()
}

func (i *TypeInspector) IsPrimitiveKind() bool {
	return i.IsIntKind() || i.IsStringKind() || i.IsBoolKind()
}

func (i *TypeInspector) IsObjectKind() bool {
	return i.IsStructKind().OK() || i.IsChanKind() || i.IsSlice() || i.IsMapKind() || i.IsArrayKind()
}

func (i *TypeInspector) IsReferenceKind() bool {
	return i.IsPointerKind() || i.IsUnsafePointerKind() || i.IsUintptrKind()
}

func (i *TypeInspector) IsPointerKind() bool {
	return i.Kind() == reflect.Pointer
}

func (i *TypeInspector) IsStructKind() logic.Check {
	return IsStructKind(i.typ)
}

func (i *TypeInspector) IsInterfaceKind() logic.Check {
	return IsInterfaceKind(i.typ)
}

func (i *TypeInspector) IsChanKind() bool {
	return i.Kind() == reflect.Chan
}

func (i *TypeInspector) IsSlice() bool {
	return i.Kind() == reflect.Slice
}

func (i *TypeInspector) IsSliceKind() logic.Check {
	return IsSliceKind(i.typ)
}

func (i *TypeInspector) IsMapKind() bool {
	return i.Kind() == reflect.Map
}

func (i *TypeInspector) IsArrayKind() bool {
	return i.Kind() == reflect.Array
}

func (i *TypeInspector) IsStringKind() bool {
	return i.Kind() == reflect.String
}

func (i *TypeInspector) IsBoolKind() bool {
	return i.Kind() == reflect.Bool
}

func (i *TypeInspector) IsIntKind() bool {
	return i.Kind() == reflect.Int
}

func (i *TypeInspector) IsComplex64Kind() bool {
	return i.Kind() == reflect.Complex64
}

func (i *TypeInspector) IsComplex128() bool {
	return i.Kind() == reflect.Complex128
}

func (i *TypeInspector) IsFloat32Kind() bool {
	return i.Kind() == reflect.Float32
}

func (i *TypeInspector) IsFloat64() bool {
	return i.Kind() == reflect.Float64
}

func (i *TypeInspector) IsInt8Kind() bool {
	return i.Kind() == reflect.Int8
}

func (i *TypeInspector) IsInt16Kind() bool {
	return i.Kind() == reflect.Int16
}

func (i *TypeInspector) IsInt32Kind() bool {
	return i.Kind() == reflect.Int32
}

func (i *TypeInspector) IsInt64Kind() bool {
	return i.Kind() == reflect.Int64
}

func (i *TypeInspector) IsUintKind() bool {
	return i.Kind() == reflect.Uint
}

func (i *TypeInspector) IsUint8Kind() bool {
	return i.Kind() == reflect.Uint8
}

func (i *TypeInspector) IsUint16Kind() bool {
	return i.Kind() == reflect.Uint16
}

func (i *TypeInspector) IsUint32Kind() bool {
	return i.Kind() == reflect.Uint32
}

func (i *TypeInspector) IsUint64Kind() bool {
	return i.Kind() == reflect.Uint64
}

func (i *TypeInspector) IsUintptrKind() bool {
	return i.Kind() == reflect.Uintptr
}

func (i *TypeInspector) IsUnsafePointerKind() bool {
	return i.Kind() == reflect.UnsafePointer
}

func (i *TypeInspector) IsInvalidKind() bool {
	return i.Kind() == reflect.Invalid
}
