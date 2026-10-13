package _struct

import (
	"iter"
	"reflect"

	typesEngine "github.com/qlustra/intel/engines/types"
)

type FieldMatcher func(structField reflect.StructField) bool
type DescriptorMatcher func(descriptor *FieldDescriptor) bool

type FieldsInspector struct {
	*DescriptorRegistry
	//*typesEngine.TypeInspector
	strct             reflect.Type
	fieldMatcher      FieldMatcher
	descriptorMatcher DescriptorMatcher
}

func InspectFields(structType *typesEngine.TypeInspector, descriptorRegistry *DescriptorRegistry, fieldMatcher FieldMatcher, descriptorMatcher DescriptorMatcher) *FieldsInspector {
	return &FieldsInspector{
		DescriptorRegistry: descriptorRegistry,
		//TypeInspector:      structType,
		strct:             structType.CapturedType(),
		fieldMatcher:      fieldMatcher,
		descriptorMatcher: descriptorMatcher,
	}
}

func (fi *FieldsInspector) matchField(structField reflect.StructField) bool {
	if fi.fieldMatcher == nil {
		return true
	}
	return fi.fieldMatcher(structField)
}

func (fi *FieldsInspector) matchDescriptor(descriptor *FieldDescriptor) bool {
	if descriptor == nil {
		return false
	}
	if fi.descriptorMatcher == nil {
		return true
	}
	return fi.descriptorMatcher(descriptor)
}

func (fi *FieldsInspector) tryRegister(structField reflect.StructField) *FieldDescriptor {
	if !fi.matchField(structField) {
		return nil
	}
	descriptor := fi.register(structField)
	if !fi.matchDescriptor(descriptor) {
		return nil
	}
	return descriptor
}

func (fi *FieldsInspector) AllAddressable() iter.Seq2[string, *FieldDescriptor] {
	return func(yield func(string, *FieldDescriptor) bool) {
		fi.strct.FieldByNameFunc(func(name string) bool {
			structField, found := fi.strct.FieldByName(name)
			if !found {
				return false
			}
			descriptor := fi.tryRegister(structField)
			if descriptor == nil {
				return false
			}
			return !yield(name, descriptor)
		})
	}
}

func (fi *FieldsInspector) AllDeclared() iter.Seq2[string, *FieldDescriptor] {
	return func(yield func(string, *FieldDescriptor) bool) {
		for structField := range fi.strct.Fields() {
			descriptor := fi.tryRegister(structField)
			if descriptor != nil {
				if !yield(descriptor.FieldName(), descriptor) {
					return
				}
			}
		}
	}
}

func (fi *FieldsInspector) LookupNamed(name string) (*FieldDescriptor, bool) {
	var descriptor *FieldDescriptor
	_, found := fi.strct.FieldByNameFunc(func(fieldName string) bool {
		if name != fieldName {
			return false
		}
		structField, found := fi.strct.FieldByName(fieldName)
		if !found {
			return false
		}
		descriptor = fi.tryRegister(structField)
		return descriptor != nil
	})
	return descriptor, found
}

func (fi *FieldsInspector) LookupType[T any]() (*FieldDescriptor, bool, bool) {
	var firstDescriptor *FieldDescriptor
	foundOne := false
	ambiguous := false
	fi.strct.FieldByNameFunc(func(fieldName string) bool {
		structField, found := fi.strct.FieldByName(fieldName)
		if !found || !typesEngine.Is[T](structField.Type) {
			return false
		}
		descriptor := fi.tryRegister(structField)
		foundOne = descriptor != nil
		if foundOne {
			if firstDescriptor == nil {
				firstDescriptor = descriptor
			} else {
				ambiguous = true
			}
		}
		return foundOne
	})
	return firstDescriptor, foundOne, ambiguous
}

func (fi *FieldsInspector) LookupContract[T any]() (*FieldDescriptor, bool, bool) {
	var firstDescriptor *FieldDescriptor
	foundOne := false
	ambiguous := false
	fi.strct.FieldByNameFunc(func(fieldName string) bool {
		structField, found := fi.strct.FieldByName(fieldName)
		if !found || !typesEngine.Satisfies[T](structField.Type) {
			return false
		}
		descriptor := fi.tryRegister(structField)
		foundOne = descriptor != nil
		if foundOne {
			if firstDescriptor == nil {
				firstDescriptor = descriptor
			} else {
				ambiguous = true
			}
		}
		return foundOne
	})
	return firstDescriptor, foundOne, ambiguous
}

func (fi *FieldsInspector) Declared(name string) (*FieldDescriptor, bool) {
	for structField := range fi.strct.Fields() {
		if name != structField.Name {
			continue
		}
		descriptor := fi.tryRegister(structField)
		if descriptor != nil {
			return descriptor, true
		}
	}
	return nil, false
}

func (fi *FieldsInspector) DeclaredType[T any]() (*FieldDescriptor, bool) {
	for structField := range fi.strct.Fields() {
		if !typesEngine.Is[T](structField.Type) {
			continue
		}
		descriptor := fi.tryRegister(structField)
		if descriptor != nil {
			return descriptor, true
		}
	}
	return nil, false
}

func (fi *FieldsInspector) DeclaredContract[T any]() (*FieldDescriptor, bool) {
	for structField := range fi.strct.Fields() {
		if !typesEngine.Satisfies[T](structField.Type) {
			continue
		}
		descriptor := fi.tryRegister(structField)
		if descriptor != nil {
			return descriptor, true
		}
	}
	return nil, false
}
