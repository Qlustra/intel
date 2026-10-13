package _struct

import (
	"reflect"

	tagEngine "github.com/qlustra/intel/engines/struct/tag"
	typeEngine "github.com/qlustra/intel/engines/types"
)

type FieldDescriptor struct {
	typeInspector *typeEngine.TypeInspector
	tags          *tagEngine.FieldTagInspector
	structField   reflect.StructField
}

func CaptureFieldDescriptor(field reflect.StructField) *FieldDescriptor {
	typeInspector := typeEngine.CaptureType(field.Type)
	tagInspector := tagEngine.CaptureFieldTag(field.Tag)
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

func (i *FieldDescriptor) FieldTypeInspector() *typeEngine.TypeInspector {
	return i.typeInspector
}

func (i *FieldDescriptor) Tags() *tagEngine.FieldTagInspector {
	return i.tags
}

func (i *FieldDescriptor) FieldName() string {
	return i.structField.Name
}

func (i *FieldDescriptor) PkgPath() string {
	return i.structField.PkgPath
}

func (i *FieldDescriptor) Anonymous() bool {
	return IsAnonymousField(i.structField)
}

func (i *FieldDescriptor) Offset() uintptr {
	return i.structField.Offset
}

func (i *FieldDescriptor) ReflectedType() reflect.Type {
	return i.structField.Type
}

func (i *FieldDescriptor) Index() []int {
	return i.structField.Index
}

func (i *FieldDescriptor) IsBlank() bool {
	return IsBlankField(i.structField)
}

func (i *FieldDescriptor) IsProperty() bool {
	return IsProperty(i.structField)
}

func (i *FieldDescriptor) IsEmbedding() bool {
	return IsEmbedding(i.structField)
}

func (i *FieldDescriptor) IsEmbeddedStruct() bool {
	return IsEmbeddedStruct(i.structField)
}

func (i *FieldDescriptor) IsEmbeddedInterface() bool {
	return IsEmbeddedInterface(i.structField)
}

func (i *FieldDescriptor) IsDirect() bool {
	return IsDeclared(i.structField)
}

func (i *FieldDescriptor) IsPromoted() bool {
	return IsPromoted(i.structField)
}

func (i *FieldDescriptor) IsExported() bool {
	return i.structField.IsExported()
}

func IsEmbedding(structField reflect.StructField) bool {
	return structField.Anonymous
}

func IsDirectEmbedding(structField reflect.StructField) bool {
	return structField.Anonymous && IsDeclared(structField)
}

func IsEmbeddedStruct(structField reflect.StructField) bool {
	return IsAnonymousField(structField) && typeEngine.IsStructKind(structField.Type)
}

func IsDirectEmbeddedStruct(structField reflect.StructField) bool {
	return IsEmbeddedStruct(structField) && IsDeclared(structField)
}

func IsEmbeddedInterface(structField reflect.StructField) bool {
	return IsAnonymousField(structField) && typeEngine.IsInterfaceKind(structField.Type)
}

func IsDirectEmbeddedInterface(structField reflect.StructField) bool {
	return IsEmbeddedInterface(structField) && IsEmbeddedInterface(structField)
}

var IsAnonymousField = IsEmbedding

func IsBlankField(structField reflect.StructField) bool {
	return structField.Name == "_"
}

func IsProperty(structField reflect.StructField) bool {
	return !IsAnonymousField(structField) && !IsBlankField(structField)
}

func IsDirectProperty(structField reflect.StructField) bool {
	return IsProperty(structField) && IsDeclared(structField)
}

func IsDeclared(structField reflect.StructField) bool {
	return len(structField.Index) == 1
}

func IsPromoted(structField reflect.StructField) bool {
	return len(structField.Index) > 1
}
