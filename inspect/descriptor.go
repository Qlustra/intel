package inspect

import (
	"reflect"

	"github.com/qlustra/intel/logic"
)

type FieldDescriptor struct {
	typeInspector *TypeInspector
	tags          *FieldTagInspector
	structField   reflect.StructField
}

// CaptureFieldDescriptor records metadata from field. It does not verify that
// the field belongs to a particular root type or that its index path resolves
// on a value; use the checked capture field operations for those guarantees.
func CaptureFieldDescriptor(field reflect.StructField) *FieldDescriptor {
	typeInspector := CaptureType(field.Type)
	tagInspector := CaptureFieldTag(field.Tag)
	descriptor := &FieldDescriptor{
		typeInspector: typeInspector,
		tags:          tagInspector,
		structField:   field,
	}
	return descriptor
}

func (i *FieldDescriptor) StructField() reflect.StructField {
	return i.structField
}

func (i *FieldDescriptor) FieldTypeInspector() *TypeInspector {
	return i.typeInspector
}

// Tags returns an inspector over the field's struct tag text.
func (i *FieldDescriptor) Tags() *FieldTagInspector {
	return i.tags
}

func (i *FieldDescriptor) FieldName() string {
	return i.structField.Name
}

func (i *FieldDescriptor) PkgPath() string {
	return i.structField.PkgPath
}

func (i *FieldDescriptor) Anonymous() logic.Check {
	return IsAnonymousField(i.structField)
}

func (i *FieldDescriptor) Offset() uintptr {
	return i.structField.Offset
}

func (i *FieldDescriptor) ReflectedType() reflect.Type {
	return i.structField.Type
}

func (i *FieldDescriptor) Is[T any]() logic.Check {
	return Is[T](i.structField.Type)
}

// Index returns the field's Go index path as stored in the descriptor. A path
// from a promoted lookup includes the steps through its embeddings. The slice
// is not copied; treat it as read-only and clone it before modifying it.
func (i *FieldDescriptor) Index() []int {
	return i.structField.Index
}

func (i *FieldDescriptor) IsBlank() logic.Check {
	return IsBlankField(i.structField)
}

func (i *FieldDescriptor) IsProperty() logic.Check {
	return IsProperty(i.structField)
}

func (i *FieldDescriptor) IsEmbedding() logic.Check {
	return IsEmbedding(i.structField)
}

func (i *FieldDescriptor) IsEmbeddedStruct() logic.Check {
	return IsEmbeddedStruct(i.structField)
}

func (i *FieldDescriptor) IsEmbeddedInterface() logic.Check {
	return IsEmbeddedInterface(i.structField)
}

func (i *FieldDescriptor) IsDirect() logic.Check {
	return IsDeclared(i.structField)
}

func (i *FieldDescriptor) DeclaredBy(name string) logic.Check {
	check := IsDeclared(i.structField)
	if !check.OK() {
		return check
	}
	if i.FieldName() != name {
		return logic.Nope("field_name_not_matched", name)
	}
	return logic.Ok()
}

func (i *FieldDescriptor) DeclaredAs[T any]() logic.Check {
	check := IsDeclared(i.structField)
	if !check.OK() {
		return check
	}
	return Is[T](i.structField.Type)
}

func (i *FieldDescriptor) IsPromoted() logic.Check {
	return IsPromoted(i.structField)
}

func (i *FieldDescriptor) IsExported() bool {
	return i.structField.IsExported()
}
