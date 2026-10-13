package _struct

import "reflect"

type DescriptorRegistry struct {
	descriptors map[string]*FieldDescriptor
}

func NewDescriptorRegistry() *DescriptorRegistry {
	return &DescriptorRegistry{
		descriptors: make(map[string]*FieldDescriptor),
	}
}

func (fi *DescriptorRegistry) register(reflectedField reflect.StructField) *FieldDescriptor {
	if fi.descriptors == nil {
		fi.descriptors = make(map[string]*FieldDescriptor)
	}
	name := reflectedField.Name

	existing, found := fi.descriptors[name]
	if found {
		return existing
	}
	descriptor := CaptureFieldDescriptor(reflectedField)
	fi.descriptors[name] = descriptor
	return descriptor
}
