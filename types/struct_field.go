package types

import (
	structEngine "github.com/qlustra/intel/engines/struct"
)

type FieldInterrogator struct {
	fieldInterrogatorMember
	*structEngine.FieldDescriptor
	properties *FieldRegistry[PropertyInterrogator]
	structs    *FieldRegistry[StructInterrogator]
	interfaces *FieldRegistry[InterfaceInterrogator]
	embeddings *FieldRegistry[EmbeddingInterrogator]
}

func captureField(descriptor *structEngine.FieldDescriptor) *FieldInterrogator {
	interrogator := &FieldInterrogator{
		FieldDescriptor: descriptor,
	}
	return interrogator
}

func (i *FieldInterrogator) AsProperty() (*PropertyInterrogator, bool) {
	if i.FieldDescriptor.IsProperty() {
		return i.properties.register(i.FieldDescriptor), true
	}
	return nil, false
}

func (i *FieldInterrogator) AsEmbedding() (*EmbeddingInterrogator, bool) {
	if i.FieldDescriptor.IsEmbedding() {
		return i.embeddings.register(i.FieldDescriptor), true
	}
	return nil, false
}

func (i *FieldInterrogator) AsEmbeddedStruct() (*StructInterrogator, bool) {
	if i.FieldDescriptor.IsEmbeddedStruct() {
		return i.structs.register(i.FieldDescriptor), true
	}
	return nil, false
}

func (i *FieldInterrogator) AsEmbeddedInterface() (*InterfaceInterrogator, bool) {
	if i.FieldDescriptor.IsEmbeddedInterface() {
		return i.interfaces.register(i.FieldDescriptor), true
	}
	return nil, false
}
