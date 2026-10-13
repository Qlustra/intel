package types

import (
	structEngine "github.com/qlustra/intel/engines/struct"
)

type EmbeddingInterrogator struct {
	fieldInterrogatorMember
	*structEngine.FieldDescriptor
	structs    *FieldRegistry[StructInterrogator]
	interfaces *FieldRegistry[InterfaceInterrogator]
}

func captureEmbedding(descriptor *structEngine.FieldDescriptor) *EmbeddingInterrogator {
	interrogator := &EmbeddingInterrogator{
		FieldDescriptor: descriptor,
	}
	return interrogator
}

func (i *EmbeddingInterrogator) AsStruct() (*StructInterrogator, bool) {
	descriptor := i.FieldDescriptor
	if descriptor.IsEmbeddedStruct() {
		return i.structs.register(i.FieldDescriptor), true
	}
	return nil, false
}

func (i *EmbeddingInterrogator) AsInterface() (*InterfaceInterrogator, bool) {
	descriptor := i.FieldDescriptor
	if descriptor.IsEmbeddedStruct() {
		return i.interfaces.register(i.FieldDescriptor), true
	}
	return nil, false
}
