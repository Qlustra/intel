package capture

import (
	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
)

// EmbeddingInterrogator describes an embedded field. The embedding may be a
// struct, an interface, or another embeddable Go type.
type EmbeddingInterrogator struct {
	fieldInterrogatorMember
	*inspect.FieldDescriptor
	structs   *FieldRegistry[ExtensionInterrogator]
	contracts *FieldRegistry[ContractInterrogator]
}

func captureEmbedding(descriptor *inspect.FieldDescriptor) *EmbeddingInterrogator {
	interrogator := &EmbeddingInterrogator{
		FieldDescriptor: descriptor,
	}
	return interrogator
}

// AsStruct narrows an embedding to an extension if its field type is a
// struct. It does not dereference embedded pointers; see [Extensions].
func (i *EmbeddingInterrogator) AsStruct() (*ExtensionInterrogator, logic.Check) {
	descriptor := i.FieldDescriptor
	check := descriptor.IsEmbeddedStruct()
	if !check.OK() {
		return nil, check
	}
	return i.structs.register(i.FieldDescriptor), check
}

// AsInterface narrows an embedding to a contract if its field type is an
// interface. It checks the declared type, not a value; see [Contracts].
func (i *EmbeddingInterrogator) AsInterface() (*ContractInterrogator, logic.Check) {
	descriptor := i.FieldDescriptor
	check := descriptor.IsEmbeddedInterface()
	if !check.OK() {
		return nil, check
	}
	return i.contracts.register(i.FieldDescriptor), check
}
