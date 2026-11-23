package inspect

import (
	"fmt"
	"reflect"

	"github.com/qlustra/intel/logic"
)

// ======================================== Kind

const (
	TypeIsNotStructKindCode    logic.CheckCode = "type_is_not_a_struct"
	TypeIsNotInterfaceKindCode logic.CheckCode = "type_is_not_an_interface"
	TypeIsNotSliceKindCode     logic.CheckCode = "type_is_not_a_slice"
)

func IsStructKind(typ reflect.Type) logic.Check {
	if typ == nil {
		return logic.Nope(TypeIsNotStructKindCode, "typ == nil")
	}
	kind := typ.Kind()
	if kind != reflect.Struct {
		return logic.Nope(TypeIsNotStructKindCode, fmt.Sprintf("%s != Struct", kind.String()))
	}
	return logic.Ok()
}

func IsInterfaceKind(typ reflect.Type) logic.Check {
	if typ == nil {
		return logic.Nope(TypeIsNotInterfaceKindCode, "typ == nil")
	}
	kind := typ.Kind()
	if kind != reflect.Interface {
		return logic.Nope(TypeIsNotInterfaceKindCode, fmt.Sprintf("%s != Interface", kind.String()))
	}
	return logic.Ok()
}

func IsSliceKind(typ reflect.Type) logic.Check {
	if typ == nil {
		return logic.Nope(TypeIsNotSliceKindCode, "typ == nil")
	}
	kind := typ.Kind()
	if kind != reflect.Slice {
		return logic.Nope(TypeIsNotSliceKindCode, fmt.Sprintf("%s != Slice", kind.String()))
	}
	return logic.Ok()
}

// ======================================== Type

const (
	TypeIsNotEqualCode     logic.CheckCode = "types_are_not_equal"
	TypeDoesNotSatisfyCode logic.CheckCode = "type_does_not_satisfy_other_type"
)

func Is[T any](typ reflect.Type) logic.Check {
	return IsType(typ, reflect.TypeFor[T]())
}

func IsType(typ, otherType reflect.Type) logic.Check {
	if typ == nil {
		return logic.Nope(TypeIsNotEqualCode, "typ == nil")
	}
	if otherType == nil {
		return logic.Nope(TypeIsNotEqualCode, "other == nil")
	}
	if typ != otherType {
		return logic.Nope(TypeIsNotEqualCode, fmt.Sprintf("%v != %v", typ, otherType))
	}
	return logic.Ok()
}

func SatisfiesType(typ, otherType reflect.Type) logic.Check {
	if typ == nil {
		return logic.Nope(TypeDoesNotSatisfyCode, "typ == nil")
	}
	if otherType == nil {
		return logic.Nope(TypeDoesNotSatisfyCode, "other == nil")
	}
	isInterfaceCheck := IsInterfaceKind(otherType)
	if !isInterfaceCheck.OK() {
		return logic.NopeBecause(TypeDoesNotSatisfyCode, "", isInterfaceCheck)
	}
	if !typ.Implements(otherType) {
		return logic.Nope(TypeDoesNotSatisfyCode, fmt.Sprintf("%v does not implement %v", typ, otherType))
	}

	return logic.Ok()
}

func Satisfies[T any](typ reflect.Type) logic.Check {
	return SatisfiesType(typ, reflect.TypeFor[T]())
}

// ======================================== StructField

const (
	FieldIsNotEmbeddingCode                 logic.CheckCode = "struct_field_is_not_an_embedding"
	FieldIsNotDirectlyEmbeddingCode         logic.CheckCode = "struct_field_is_not_a_directly_declared_embedding"
	FieldIsNotEmbeddedStructCode            logic.CheckCode = "struct_field_is_not_an_embedded_struct"
	FieldIsNotDirectlyEmbeddedStructCode    logic.CheckCode = "struct_field_is_not_a_directly_embedded_struct"
	FieldIsNotEmbeddedInterfaceCode         logic.CheckCode = "struct_field_is_not_an_embedded_interface"
	FieldIsNotDirectlyEmbeddedInterfaceCode logic.CheckCode = "struct_field_is_not_a_directly_embedded_interface"
	FieldIsNotDeclaredCode                  logic.CheckCode = "struct_field_is_not_directly_declared_on_the_struct"
	FieldIsNotBlankFieldCode                logic.CheckCode = "struct_field_is_not_blank"
	FieldIsNotPropertyCode                  logic.CheckCode = "struct_field_is_not_a_property"
	FieldIsNotDirectPropertyCode            logic.CheckCode = "struct_field_is_not_a_directly_declared_property"
	FieldIsNotPromotedCode                  logic.CheckCode = "struct_field_is_not_promoted"
)

func IsEmbedding(structField reflect.StructField) logic.Check {
	if structField.Anonymous {
		return logic.Ok()
	}
	return logic.Nope(FieldIsNotEmbeddingCode, fmt.Sprintf("%s not Anonymous", structField.Name))
}

func IsDirectEmbedding(structField reflect.StructField) logic.Check {
	isEmbeddingCheck := IsEmbedding(structField)
	if !isEmbeddingCheck.OK() {
		return logic.NopeBecause(FieldIsNotDirectlyEmbeddingCode, "", isEmbeddingCheck)
	}
	isDeclaredCheck := IsDeclared(structField)
	if !isDeclaredCheck.OK() {
		return logic.NopeBecause(FieldIsNotDirectlyEmbeddingCode, "", isDeclaredCheck)
	}
	return logic.Ok()
}

func IsEmbeddedStruct(structField reflect.StructField) logic.Check {
	isAnonymousCheck := IsAnonymousField(structField)
	if !isAnonymousCheck.OK() {
		return logic.NopeBecause(FieldIsNotEmbeddedStructCode, "", isAnonymousCheck)
	}
	isStructKindCheck := IsStructKind(structField.Type)
	if !isStructKindCheck.OK() {
		return logic.NopeBecause(FieldIsNotEmbeddedStructCode, "", isStructKindCheck)
	}
	return logic.Ok()
}

func IsDirectEmbeddedStruct(structField reflect.StructField) logic.Check {
	isEmbeddedStructCheck := IsEmbeddedStruct(structField)
	if !isEmbeddedStructCheck.OK() {
		return logic.NopeBecause(FieldIsNotDirectlyEmbeddedStructCode, "", isEmbeddedStructCheck)
	}
	isDeclaredCheck := IsDeclared(structField)
	if !isDeclaredCheck.OK() {
		return logic.NopeBecause(FieldIsNotDirectlyEmbeddedStructCode, "", isDeclaredCheck)
	}
	return logic.Ok()
}

func IsEmbeddedInterface(structField reflect.StructField) logic.Check {
	isAnonymousFieldCheck := IsAnonymousField(structField)
	if !isAnonymousFieldCheck.OK() {
		return logic.NopeBecause(FieldIsNotEmbeddedInterfaceCode, "", isAnonymousFieldCheck)
	}
	isInterfaceKindCheck := IsInterfaceKind(structField.Type)
	if !isInterfaceKindCheck.OK() {
		return logic.NopeBecause(FieldIsNotEmbeddedInterfaceCode, "", isInterfaceKindCheck)
	}
	return logic.Ok()
}

func IsDirectEmbeddedInterface(structField reflect.StructField) logic.Check {
	isEmbeddedInterfaceCheck := IsEmbeddedInterface(structField)
	if !isEmbeddedInterfaceCheck.OK() {
		return logic.NopeBecause(FieldIsNotDirectlyEmbeddedInterfaceCode, "", isEmbeddedInterfaceCheck)
	}
	isDeclaredCheck := IsDeclared(structField)
	if !isDeclaredCheck.OK() {
		return logic.NopeBecause(FieldIsNotDirectlyEmbeddedInterfaceCode, "", isDeclaredCheck)
	}
	return logic.Ok()
}

var IsAnonymousField = IsEmbedding

func IsBlankField(structField reflect.StructField) logic.Check {
	if structField.Name == "_" {
		return logic.Ok()
	}
	return logic.Nope(FieldIsNotBlankFieldCode, fmt.Sprintf("%s not _", structField.Name))
}

func IsProperty(structField reflect.StructField) logic.Check {
	if IsAnonymousField(structField).OK() {
		return logic.Nope(FieldIsNotPropertyCode, fmt.Sprintf("%s is Anonymous", structField.Name))
	}
	if IsBlankField(structField).OK() {
		return logic.Nope(FieldIsNotPropertyCode, fmt.Sprintf("%s is blank", structField.Name))
	}
	return logic.Ok()
}

func IsDirectProperty(structField reflect.StructField) logic.Check {
	isPropertyCheck := IsProperty(structField)
	if !isPropertyCheck.OK() {
		return logic.NopeBecause(FieldIsNotDirectPropertyCode, "", isPropertyCheck)
	}
	isDeclaredCheck := IsDeclared(structField)
	if !isDeclaredCheck.OK() {
		return logic.NopeBecause(FieldIsNotDirectPropertyCode, "", isDeclaredCheck)
	}
	return logic.Ok()
}

func IsDeclared(structField reflect.StructField) logic.Check {
	if len(structField.Index) == 1 {
		return logic.Ok()
	}
	return logic.Nope(FieldIsNotDeclaredCode, fmt.Sprintf("%s.Index is not 1", structField.Name))
}

func IsPromoted(structField reflect.StructField) logic.Check {
	if len(structField.Index) > 1 {
		return logic.Ok()
	}
	return logic.Nope(FieldIsNotPromotedCode, fmt.Sprintf("%s.Index is not > 1", structField.Name))
}
